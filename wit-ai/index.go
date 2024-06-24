package wit_ai

import (
	"ai_bot/contract"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/umbracle/ethgo"
	witai "github.com/wit-ai/wit-go/v2"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	client     *witai.Client
	nftNameMap = map[string]string{
		"qwe": "0x85A679b0b57C2486688a564221001F0b88323f9B", //bsc
		"cxp": "0xB6cBBbbF49664c749Fc519d7d03194C22645CC31", //remix
	}
	desSelectMap = map[string]uint64{
		"bsc_test":   13264668187771770619,
		"wemix_test": 9284632837123596123,
	}
	receiverMap = map[string]ethgo.Address{
		"wemix_test": ethgo.HexToAddress("0xD6C1e806B29D22B862e5c8AA2a35CE2e98B82002"),
		"bsc_test":   ethgo.HexToAddress("0xBcDD9f7835d0994Bde8bA1D24ffC2AF0cbBdF64e"),
	}
)

func Init(token string) {
	proxyStr := os.Getenv("HTTP_PROXY")
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		panic(err)
	}
	client = witai.NewClient(token)
	client.SetHTTPClient(&http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	})
}

func ParseMessage(query, privateKey string) (string, error) {
	msg, err := client.Parse(&witai.MessageRequest{
		Query: query,
	})
	if err != nil {
		logrus.Printf("Error parsing message: %v", err)
		return "", err
	}

	indent, _ := json.MarshalIndent(msg, "", "  ")
	fmt.Println(string(indent))
	jsonData := string(indent)

	intent := gjson.Get(jsonData, "intents.0.name").String()
	fmt.Println("Intent:", intent)

	placeContract := contract.NewNFTMarketPlaceContract(privateKey)

	switch intent {
	case "listNFT":
		return handleListNFT(jsonData, placeContract)
	case "buyNFT":
		return handleBuyNFT(jsonData, placeContract)
	default:
		return "", fmt.Errorf("unknown intent: %s", intent)
	}
}

func handleListNFT(jsonData string, placeContract *contract.NFTMarketPlaceContract) (string, error) {
	var (
		nftName      string
		tokenID      uint64
		price        int64
		priceInWei   *big.Int
		chainId      uint
		isCrossChain bool
	)

	if gjson.Get(jsonData, "entities.SystemNFT:SystemNFT").Exists() {
		nftName = gjson.Get(jsonData, "entities.SystemNFT:SystemNFT.0.value").String()
	} else {
		return "", fmt.Errorf("NFT name not found")
	}

	if gjson.Get(jsonData, "entities.TokenID:TokenID").Exists() {
		tokenID = gjson.Get(jsonData, "entities.TokenID:TokenID.0.value").Uint()
	} else {
		return "", fmt.Errorf("Token ID not found")
	}

	if gjson.Get(jsonData, "entities.Price:Price").Exists() {
		price = gjson.Get(jsonData, "entities.Price:Price.0.value").Int()
		priceInWei = new(big.Int).Mul(big.NewInt(price), big.NewInt(1e18))
	} else {
		return "", fmt.Errorf("Price not found")
	}

	// 打印提取的信息
	fmt.Println("NFT Name:", nftName)
	fmt.Println("Token ID:", tokenID)
	fmt.Println("Price:", price)

	nftAddress := nftNameMap[nftName]

	if nftName == "qwe" {
		chainId = 97
		isCrossChain = false
	} else if nftName == "cxp" {
		chainId = 1112
		isCrossChain = true
	}
	//打印所有函数参数
	fmt.Printf("nftAddress: %s, tokenID: %d, priceInWei: %s, chainId: %d, isCrossChain: %t\n", nftAddress, tokenID, priceInWei.String(), chainId, isCrossChain)

	err := placeContract.ListItem(ethgo.HexToAddress(nftAddress), uint(tokenID), priceInWei, chainId, isCrossChain)

	if err != nil {
		return "", err
	}

	return "NFT listed successfully", nil
}

func handleBuyNFT(jsonData string, placeContract *contract.NFTMarketPlaceContract) (string, error) {
	var (
		nftName                  string
		tokenID                  uint64
		success                  bool
		destinationChainSelector uint64
		isCrossChain             bool
		receiver                 ethgo.Address
	)

	if gjson.Get(jsonData, "entities.SystemNFT:SystemNFT").Exists() {
		nftName = gjson.Get(jsonData, "entities.SystemNFT:SystemNFT.0.value").String()
	} else {
		return "", fmt.Errorf("NFT name not found")
	}

	if gjson.Get(jsonData, "entities.TokenID:TokenID").Exists() {
		tokenID = gjson.Get(jsonData, "entities.TokenID:TokenID.0.value").Uint()
	} else {
		return "", fmt.Errorf("Token ID not found")
	}

	fmt.Println("NFT Name:", nftName)
	fmt.Println("Token ID:", tokenID)

	nftAddress := nftNameMap[nftName]
	if nftName == "qwe" {
		destinationChainSelector = desSelectMap["bsc_test"]
		receiver = receiverMap["bsc_test"]
		isCrossChain = false
	} else if nftName == "cxp" {
		destinationChainSelector = desSelectMap["wemix_test"]
		receiver = receiverMap["wemix_test"]
		isCrossChain = true
	}
	//print all param
	logrus.Printf("nftAddress: %s, tokenID: %d, destinationChainSelector: %d, receiver: %s, isCrossChain: %t\n", nftAddress, tokenID, destinationChainSelector, receiver, isCrossChain)
	success = placeContract.BuyItem(nftAddress, uint(tokenID), destinationChainSelector, receiver, isCrossChain)

	if !success {
		return "", fmt.Errorf("failed to buy NFT")
	}

	return "NFT bought successfully", nil
}

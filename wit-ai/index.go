package wit_ai

import (
	"ai_bot/contract"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
	"os"
)

type ClientWrapper struct {
	*witai.Client
}

func NewClientWrapper(token string) *ClientWrapper {
	client := witai.NewClient(token)
	return &ClientWrapper{client}
}

func (handler *ClientWrapper) ParseMessage(query string) (string, error) {
	msg, err := handler.Parse(&witai.MessageRequest{
		Query: query,
	})
	if err != nil {
		return "", err
	}

	indent, _ := json.MarshalIndent(msg, "", "  ")
	fmt.Println(string(indent))
	jsonData := string(indent)
	// 使用 gjson 提取信息
	nftName := ""
	tokenID := ""
	price := ""
	intent := gjson.Get(jsonData, "intents.0.name").String()
	fmt.Println("Intent:", intent)
	if intent == "listNFT" {
		if gjson.Get(jsonData, "entities.SystemNFT:SystemNFT").Exists() {
			nftName = gjson.Get(jsonData, "entities.SystemNFT:SystemNFT.0.value").String()
		}

		if gjson.Get(jsonData, "entities.TokenID:TokenID").Exists() {
			tokenID = gjson.Get(jsonData, "entities.TokenID:TokenID.0.value").String()
		}

		if gjson.Get(jsonData, "entities.Price:Price").Exists() {
			price = gjson.Get(jsonData, "entities.Price:Price.0.value").String()
		}

		// 打印提取的信息
		fmt.Println("NFT Name:", nftName)
		fmt.Println("Token ID:", tokenID)
		fmt.Println("Price:", price)
	}
	placeContract := contract.NewNFTMarketPlaceContract(os.Getenv("PRIVATE_KEY"))
	nftAddress := placeContract.GetNFTAddressByName(nftName)
	fmt.Println("NFT Address:", nftAddress)

	return string(indent), nil
}

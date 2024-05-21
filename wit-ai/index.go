package wit_ai

import (
	"ai_bot/contract"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
	"math/big"
)

var (
	client *witai.Client
)

func Init(token string) {
	client = witai.NewClient(token)
}

func ParseMessage(query, privateKey string) error {
	msg, err := client.Parse(&witai.MessageRequest{
		Query: query,
	})
	if err != nil {
		return err
	}

	indent, _ := json.MarshalIndent(msg, "", "  ")
	fmt.Println(string(indent))
	jsonData := string(indent)

	intent := gjson.Get(jsonData, "intents.0.name").String()
	fmt.Println("Intent:", intent)

	placeContract := contract.NewNFTMarketPlaceContract(privateKey)

	switch intent {
	case "listNFT":
		{
			var (
				nftName    string
				tokenID    uint64
				price      int64
				priceInWei *big.Int
			)

			if gjson.Get(jsonData, "entities.SystemNFT:SystemNFT").Exists() {
				nftName = gjson.Get(jsonData, "entities.SystemNFT:SystemNFT.0.value").String()
			} else {
				return fmt.Errorf("NFT name not found")
			}

			if gjson.Get(jsonData, "entities.TokenID:TokenID").Exists() {
				tokenID = gjson.Get(jsonData, "entities.TokenID:TokenID.0.value").Uint()
			} else {
				return fmt.Errorf("Token ID not found")
			}

			if gjson.Get(jsonData, "entities.Price:Price").Exists() {
				price = gjson.Get(jsonData, "entities.Price:Price.0.value").Int()
				priceInWei = new(big.Int).Mul(big.NewInt(price), big.NewInt(1e18))
			} else {
				return fmt.Errorf("Price not found")
			}

			// 打印提取的信息
			fmt.Println("NFT Name:", nftName)
			fmt.Println("Token ID:", tokenID)
			fmt.Println("Price:", price)

			nftAddress := placeContract.GetNFTAddressByName(nftName)

			err := placeContract.ListItem(nftAddress, uint(tokenID), priceInWei, 97, false)
			if err != nil {
				return err
			}
			fmt.Println("NFT Address:", nftAddress)
			//log.Debug("NFT Address:", nftAddress)
			break
		}
	}
	return nil
}

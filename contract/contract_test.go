package contract

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestNFTMarketPlaceContract_GetNFTAddressByName(t *testing.T) {
	godotenv.Load("../.env")
	placeContract := NewNFTMarketPlaceContract(os.Getenv("PRIVATE_KEY"))
	placeContract.GetNFTAddressByName("qwe")
}

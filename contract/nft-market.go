package contract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
	"math/big"
	"os"
)

type NFTMarketPlaceContract struct {
	Instance *contract.Contract
}

func NewNFTMarketPlaceContract(privateKey string) *NFTMarketPlaceContract {

	nftMarketPlaceContract := new(NFTMarketPlaceContract)

	nftMarketPlaceAbi, _ := abi.NewABIFromList(NFTMarketPlaceABI)

	contractAddr := ethgo.HexToAddress(os.Getenv("NFT_MARKET_ADDRESS"))

	client, _ := jsonrpc.NewClient(os.Getenv("BSC_TEST_RPC"))

	priKey, _ := wallet.NewWalletFromPrivKey(hexutil.MustDecode(privateKey))

	nftMarketPlaceContract.Instance = contract.NewContract(contractAddr, nftMarketPlaceAbi, contract.WithJsonRPC(client.Eth()), contract.WithSender(priKey))

	return nftMarketPlaceContract
}

func (nftMarketPlaceContract *NFTMarketPlaceContract) GetNFTAddressByName(name string) ethgo.Address {
	result, err := nftMarketPlaceContract.Instance.Call("nftContractsByNames", ethgo.Latest, name)
	if err != nil {
		log.Error("Error getting NFT address by name", "err", err)
		return ethgo.ZeroAddress
	}
	fmt.Print("NFT address by name", "address", result)
	return result["0"].(ethgo.Address)
}

func (nftMarketPlaceContract *NFTMarketPlaceContract) ListItem(nftAddr ethgo.Address, tokenId uint, price *big.Int, chainId uint, isCrossChain bool) error {
	tx, txErr := nftMarketPlaceContract.Instance.Txn("listItem", nftAddr, tokenId, price, chainId, isCrossChain)
	if txErr != nil {
		return txErr
	}
	err := tx.Do()
	if err != nil {
		return err
	}
	txInfo, waitErr := tx.Wait()
	if waitErr != nil {
		return waitErr
	}
	log.Info("Listing item transaction info", "txInfo", txInfo)
	return nil
}

func (nftMarketPlaceContract *NFTMarketPlaceContract) BuyItem(nftAddr string, tokenId uint, destinationChainSelector uint64, receiver common.Address, isCrossChain bool) bool {
	tx, txErr := nftMarketPlaceContract.Instance.Txn("buyItem", common.HexToAddress(nftAddr), tokenId, destinationChainSelector, receiver, isCrossChain)
	if txErr != nil {
		log.Error("Error buying item", "err", txErr)
		return false
	}
	err := tx.Do()
	if err != nil {
		log.Error("Error executing transaction", "err", err)
		return false
	}
	txInfo, waitErr := tx.Wait()
	if waitErr != nil {
		log.Error("Error waiting for transaction", "err", waitErr)
		return false
	}
	log.Info("Buying item transaction info", "txInfo", txInfo)
	return true
}

// 添加其他函数调用类似上面的方式

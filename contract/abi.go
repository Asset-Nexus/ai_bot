package contract

var (
	NFTMarketPlaceABI = []string{
		"function listItem(address,uint256,uint256,uint256,bool)",
		"function buyItem(address,uint256,uint64,address,bool)",
		"function cancelListing(address,uint256)",
		"function updateListing(address,uint256,uint256)",
		"function nftContractsByNames(string) view returns (address)",

		// 添加其他主要函数的签名
	}
)

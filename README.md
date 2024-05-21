# NFT Marketplace Bot with Go and Wit.ai

This project demonstrates how to use Go and Wit.ai to interact with an NFT Marketplace smart contract. The bot can list NFTs, buy NFTs, and cancel listings based on user commands.
![image](https://github.com/Asset-Nexus/ai_bot/assets/84974164/473844ad-d735-45fc-bb84-57dceed46061)


## Technologies Used

- **Go**: The primary programming language used to build the bot.
- **Wit.ai**: Natural Language Processing (NLP) platform to interpret user commands.
- **Ethereum**: Blockchain network where the NFT Marketplace smart contract is deployed.
- **ethgo**: Go library for interacting with Ethereum smart contracts.
- **gjson**: Go library for JSON parsing.
- **godotenv**: Go library for loading environment variables from a `.env` file.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/nft-marketplace-bot.git
    cd nft-marketplace-bot
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

3. Create a `.env` file in the root directory with the necessary environment variables:

    ```plaintext
    PRIVATE_KEY=your_private_key
    WIT_AI_TOKEN=your_wit_ai_token
    ETHEREUM_RPC_URL=your_ethereum_rpc_url
    NFT_MARKETPLACE_CONTRACT_ADDRESS=your_contract_address
    ```

## Wit.ai Setup

1. Create a new Wit.ai app.

2. Add intents for the NFT actions:
    - `listNFT`
    - `buyNFT`
    - `cancelListing`

3. Add entities to capture relevant data:
    - `nftSeries` (e.g., CryptoKitties)
    - `tokenId` (e.g., 99)
    - `price` (e.g., 10)
    - `chainId` (e.g., 1)

4. Train the Wit.ai app with sample utterances such as:
    - "I want to list my CryptoKitties NFT with ID 99 for 10 tokens."
    - "Buy the CryptoKitties NFT with ID 99."
    - "Cancel my listing for CryptoKitties NFT with ID 99."

## Example Usage

After setting up the environment and Wit.ai intents/entities, run the bot:

```bash
go run main.go
```

The bot will listen for commands and interact with the NFT Marketplace smart contract accordingly.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

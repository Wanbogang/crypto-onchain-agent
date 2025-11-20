# Crypto Onchain Intelligence V2 Agent

[![Built with Teneo Protocol](https://img.shields.io/badge/Built%20with-Teneo%20Protocol-blue)](https://github.com/TeneoProtocolAI/teneo-agent-sdk)
[![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

Real-time cryptocurrency price tracking, wallet balance checking, and smart contract scanning agent powered by **Teneo Protocol**.

## 🚀 Features

- ✅ **Real-time Price Tracking** - Fetch latest crypto prices from CoinGecko
- ✅ **Multi-Chain Support** - Ethereum, BSC, Polygon, Arbitrum, Base
- ✅ **Wallet Analytics** - Check balance across multiple chains
- ✅ **Smart Contract Detection** - Scan and analyze contract bytecode
- ✅ **24/7 Operation** - Built-in tmux support for continuous operation
- ✅ **Teneo Protocol Integration** - Full SDK implementation

## 🛠️ Tech Stack

- **Framework**: [Teneo Agent SDK](https://github.com/TeneoProtocolAI/teneo-agent-sdk)
- **Blockchain**: [go-ethereum](https://github.com/ethereum/go-ethereum)
- **Price Data**: [CoinGecko API](https://www.coingecko.com/api)
- **Language**: Go 1.18+

## 📋 Commands

### Price Check

    price eth

    price btc

    price bnb

    price ethereum

### Wallet Balance

    wallet eth 0x742d35Cc6634C0532925a3b844Bc9e7595f42471

    wallet bsc 0x742d35Cc6634C0532925a3b844Bc9e7595f42471

    wallet polygon 0x742d35Cc6634C0532925a3b844Bc9e7595f42471

    wallet arbitrum 0x742d35Cc6634C0532925a3b844Bc9e7595f42471

    wallet base 0x742d35Cc6634C0532925a3b844Bc9e7595f42471

### Contract Scanning

    scan_contract eth 0x2170Ed0880ac9A755fd29B2688956BD959F933F8

    scan_contract bsc 0x55d398326f99059fF775485246999027B3197955

    scan_contract polygon 0x2791Bca1f2de4661ED88A30C99a7a9449Aa84174

## 🔧 Installation

### Prerequisites
- Go 1.18 or higher
- Teneo Agent SDK
- Active RPC endpoints

### Setup

1. **Clone repository**
```bash
git clone https://github.com/YOUR_USERNAME/crypto-onchain-agent.git
cd crypto-onchain-agent
```

2. **Install dependencies**
```bash
go mod download
go mod tidy
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your RPC URLs and Teneo credentials
nano .env
```

4. **Run agent**
```bash
go run main.go
```

## 📚 Example Test Cases

### Test Price Command
price eth
price btc
price bnb

Expected Output:
Price for ethereum: 2543.500000 USD ($2543.5)

### Test Wallet Command
wallet eth 0x742d35Cc6634C0532925a3b844Bc9e7595f42471

Expected Output:
[ETH] Balance 0x742d35Cc6634C0532925a3b844Bc9e7595f42471: 5.123456789012345678 ETH

### Test Contract Scan
scan_contract eth 0x2170Ed0880ac9A755fd29B2688956BD959F933F8

Expected Output:
[ETH] 0x2170Ed0880ac9A755fd29B2688956BD959F933F8 is a contract (code size = 8192 bytes)
## 🏗️ Architecture

### Main Components

1. **Agent Handler** - Processes Teneo tasks
2. **Price Module** - Integrates with CoinGecko API
3. **Wallet Module** - Queries blockchain balances
4. **Contract Scanner** - Detects and analyzes smart contracts

### Supported Chains

| Chain | Status | RPC Provider |
|-------|--------|--------------|
| Ethereum | ✅ Active | eth.merkle.io |
| BSC | ✅ Active | bsc-rpc.publicnode.com |
| Polygon | ✅ Active | polygon-rpc.com |
| Arbitrum | ✅ Active | arb1.arbitrum.io |
| Base | ✅ Active | mainnet.base.org |

## 🔐 Environment Variables

RPC Endpoints (Required)

    ETH_RPC_URL=              # Ethereum mainnet RPC
    BSC_RPC_URL=              # Binance Smart Chain RPC
    POLYGON_RPC_URL=          # Polygon mainnet RPC
    ARBITRUM_RPC_URL=         # Arbitrum One RPC
    BASE_RPC_URL=             # Base mainnet RPC

Teneo Configuration (Required)

    PRIVATE_KEY=              # Your private key
    NFT_TOKEN_ID=             # NFT token ID
    OWNER_ADDRESS=            # Owner wallet address

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For issues and support:
- [Teneo Protocol](https://teneo-protocol.ai/)
- [Teneo Chat Room](https://developer.chatroom.teneo-protocol.ai/)
- [Teneo Agent SDK](https://github.com/TeneoProtocolAI/teneo-agent-sdk)

## 🙏 Acknowledgments

This project is built with:
- [Teneo Agent SDK](https://github.com/TeneoProtocolAI/teneo-agent-sdk) - Core agent framework
- [go-ethereum](https://github.com/ethereum/go-ethereum) - Ethereum client
- [CoinGecko API](https://www.coingecko.com/api) - Price data

---

**Built by:** Yunaiko  
**Framework:** Teneo Protocol  
**Last Updated:** 2025-11-20

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TeneoProtocolAI/teneo-agent-sdk/pkg/agent"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type CryptoOnchainIntelV2Agent struct {
	lastPrices map[string]float64
}

var symbolToID = map[string]string{
	"eth":   "ethereum",
	"btc":   "bitcoin",
	"bnb":   "binancecoin",
	"matic": "matic-network",
	"arb":   "arbitrum",
	"base":  "base",
}

func coingeckoPrice(id string) (float64, string, error) {
	// Normalize input
	id = strings.ToLower(strings.TrimSpace(id))
	
	// Map symbol to proper ID
	if v, ok := symbolToID[id]; ok {
		id = v
	}
	
	u := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", id)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "teneo-agent/crypto-intel")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("failed fetch price: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return 0, "", fmt.Errorf("coingecko status %d: %s", resp.StatusCode, string(body))
	}

	var data map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Handle empty response
	if data == nil || len(data) == 0 {
		return 0, "", fmt.Errorf("empty response from coingecko for '%s'", id)
	}

	// Find the price (case-insensitive)
	var price float64
	var found bool
	for _, val := range data {
		if val != nil {
			if p, ok := val["usd"]; ok {
				price = p
				found = true
				break
			}
		}
	}

	if !found {
		j, _ := json.Marshal(data)
		return 0, "", fmt.Errorf("USD price not found for '%s' (raw: %s)", id, string(j))
	}

	return price, fmt.Sprintf("$%.6g", price), nil
}

func ethBalance(ctx context.Context, rpc string, address string) (string, error) {
	if rpc == "" {
		return "", fmt.Errorf("RPC URL is empty")
	}

	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		return "", fmt.Errorf("failed connect rpc: %w", err)
	}
	defer client.Close()
	
	addr := common.HexToAddress(address)
	bal, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return "", fmt.Errorf("failed get balance: %w", err)
	}
	
	fbal := new(big.Float).SetInt(bal)
	ethValue := new(big.Float).Quo(fbal, big.NewFloat(1e18))
	return fmt.Sprintf("%s ETH", ethValue.Text('f', 18)), nil
}

func isContract(ctx context.Context, rpc string, address string) (bool, int, error) {
	if rpc == "" {
		return false, 0, fmt.Errorf("RPC URL is empty")
	}

	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		return false, 0, fmt.Errorf("failed connect rpc: %w", err)
	}
	defer client.Close()
	
	addr := common.HexToAddress(address)
	code, err := client.CodeAt(ctx, addr, nil)
	if err != nil {
		return false, 0, fmt.Errorf("failed get code: %w", err)
	}
	
	return len(code) > 0, len(code), nil
}

func getRPC(chain string) string {
	switch strings.ToLower(strings.TrimSpace(chain)) {
	case "eth", "ethereum":
		return os.Getenv("ETH_RPC_URL")
	case "bsc", "binance":
		return os.Getenv("BSC_RPC_URL")
	case "polygon", "matic":
		return os.Getenv("POLYGON_RPC_URL")
	case "arbitrum":
		return os.Getenv("ARBITRUM_RPC_URL")
	case "base":
		return os.Getenv("BASE_RPC_URL")
	default:
		return ""
	}
}

func (a *CryptoOnchainIntelV2Agent) ProcessTask(ctx context.Context, task string) (string, error) {
	log.Printf("Processing task: %s", task)
	
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	task = strings.TrimSpace(task)
	task = strings.TrimPrefix(task, "/")
	parts := strings.Fields(task)
	
	if len(parts) == 0 {
		return "No command provided. Available commands: price, wallet, scan_contract", nil
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch cmd {
	case "price":
		if len(args) == 0 {
			return "Usage: price <coin_id_or_symbol>. Example: price ethereum", nil
		}
		id := strings.Join(args, "-")
		priceVal, priceStr, err := coingeckoPrice(id)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Price for %s: %.6f USD (%s)", id, priceVal, priceStr), nil

	case "wallet":
		var chain, addr string
		if len(args) == 1 {
			chain = "eth"
			addr = args[0]
		} else if len(args) >= 2 {
			chain = strings.ToLower(args[0])
			addr = args[1]
		} else {
			return "Usage: wallet <chain> <0xaddress>. Example: wallet eth 0x123...", nil
		}
		
		rpc := getRPC(chain)
		if rpc == "" {
			return fmt.Sprintf("RPC URL for chain '%s' not configured in .env", chain), nil
		}
		
		if !common.IsHexAddress(addr) {
			return "Invalid address format. Expected 0x...", nil
		}
		
		bal, err := ethBalance(ctx, rpc, addr)
		if err != nil {
			return "", err
		}
		
		return fmt.Sprintf("[%s] Balance %s: %s", strings.ToUpper(chain), addr, bal), nil

	case "scan_contract":
		var chain, addr string
		if len(args) == 1 {
			chain = "eth"
			addr = args[0]
		} else if len(args) >= 2 {
			chain = strings.ToLower(args[0])
			addr = args[1]
		} else {
			return "Usage: scan_contract <chain> <0xaddress>", nil
		}
		
		rpc := getRPC(chain)
		if rpc == "" {
			return fmt.Sprintf("RPC URL for chain '%s' not configured in .env", chain), nil
		}
		
		if !common.IsHexAddress(addr) {
			return "Invalid address format. Expected 0x...", nil
		}
		
		contract, size, err := isContract(ctx, rpc, addr)
		if err != nil {
			return "", err
		}
		
		status := "is a contract"
		if !contract {
			status = "is NOT a contract"
		}
		
		return fmt.Sprintf("[%s] %s %s (code size = %d bytes)", strings.ToUpper(chain), addr, status, size), nil

	default:
		return fmt.Sprintf("Unknown command '%s'. Available commands: price, wallet, scan_contract", cmd), nil
	}
}

func main() {
	_ = godotenv.Load()
	
	config := agent.DefaultConfig()
	config.Name = "Crypto Onchain Intelligence V2"
	config.Description = "Providing fast on-chain intelligence with price checks, multi-chain wallet balance, pump/dump alerts, and smart contract risk scans."
	config.Capabilities = []string{"price_check", "wallet_analysis", "contract_risk_scan"}
	config.PrivateKey = os.Getenv("PRIVATE_KEY")
	config.NFTTokenID = os.Getenv("NFT_TOKEN_ID")
	config.OwnerAddress = os.Getenv("OWNER_ADDRESS")

	enhancedAgent, err := agent.NewEnhancedAgent(&agent.EnhancedAgentConfig{
		Config:       config,
		AgentHandler: &CryptoOnchainIntelV2Agent{},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting Crypto Onchain Intelligence V2...")
	enhancedAgent.Run()
}

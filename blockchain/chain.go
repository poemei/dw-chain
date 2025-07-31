// dw-chain/blockchain/chain.go

package blockchain

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var Blockchain []Block
var TxPool *txPool

const chainFile = "data/chain.json"

func init() {
	TxPool = &txPool{pool: []Transaction{}}
	loadChain()
}

func loadChain() {
	data, err := os.ReadFile(chainFile)
	if err != nil || len(data) == 0 {
		log.Println("[Chain] No chain found, creating genesis block...")
		genesis := Block{
			Index:        0,
			Timestamp:    TimestampNow(),
			Transactions: []Transaction{},
			PrevHash:     "",
			Nonce:        0,
			MerkleRoot:   "",
		}
		genesis.Hash = genesis.CalculateHash()
		Blockchain = []Block{genesis}
		saveChain()
		return
	}
	err = json.Unmarshal(data, &Blockchain)
	if err != nil {
		log.Fatalf("[Chain] Failed to parse chain: %v", err)
	}
	log.Printf("[Chain] Loaded %d blocks", len(Blockchain))
}

func saveChain() {
	data, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		log.Printf("[Chain] Failed to encode chain: %v", err)
		return
	}
	err = os.WriteFile(chainFile, data, 0644)
	if err != nil {
		log.Printf("[Chain] Failed to save chain: %v", err)
	}
}

func AddBlock(b Block) {
	Blockchain = append(Blockchain, b)
	saveChain()
}

func GetLastBlock() Block {
	if len(Blockchain) == 0 {
		return Block{}
	}
	return Blockchain[len(Blockchain)-1]
}

func GetChain() []Block {
	return Blockchain
}

func GetChainStats() map[string]interface{} {
	latest := GetLastBlock()
	return map[string]interface{}{
		"height":     len(Blockchain),
		"latestHash": latest.Hash,
		"timestamp":  latest.Timestamp,
		"threats":    countThreats(),
	}
}

func countThreats() int {
	total := 0
	for _, block := range Blockchain {
		total += len(block.Transactions)
	}
	return total
}

func TimestampNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}

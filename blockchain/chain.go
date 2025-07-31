// dw-chain/blockchain/chain.go

package blockchain

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Constants used across the chain
const (
	chainFile          = "./data/chain.json"
	initialDifficulty  = 3
	maxDifficulty      = 6
	difficultyWindow   = 5
	targetBlockSeconds = 10
)

// Global blockchain state
var Chain []Block

// InitChain loads or creates the blockchain on node startup.
func InitChain() {
	if _, err := os.Stat(chainFile); os.IsNotExist(err) {
		log.Println("[Chain] No chain found, creating genesis block")
		genesis := Block{
			Index:        0,
			Timestamp:    TimestampNow(),
			Transactions: []Transaction{},
			PrevHash:     "",
			Nonce:        0,
			MerkleRoot:   "",
		}
		genesis.Hash = genesis.CalculateHash()
		Chain = []Block{genesis}
		SaveChain()
	} else {
		LoadChain()
	}
}

// AddBlock adds a mined block to the chain after validation.
func AddBlock(newBlock Block) {
	lastBlock := GetLastBlock()
	if IsValidNewBlock(newBlock, lastBlock) {
		Chain = append(Chain, newBlock)
		SaveChain()
		log.Printf("[Chain] Block #%d added", newBlock.Index)
	} else {
		log.Printf("[Chain] Rejected invalid block #%d", newBlock.Index)
	}
}

// IsValidNewBlock checks if a new block is valid against the last block.
func IsValidNewBlock(newBlock, prevBlock Block) bool {
	if newBlock.Index != prevBlock.Index+1 {
		return false
	}
	if newBlock.PrevHash != prevBlock.Hash {
		return false
	}
	if newBlock.Hash != newBlock.CalculateHash() {
		return false
	}
	return true
}

// GetLastBlock returns the most recent block in the chain.
func GetLastBlock() Block {
	return Chain[len(Chain)-1]
}

// SaveChain persists the chain to disk.
func SaveChain() {
	data, err := json.MarshalIndent(Chain, "", "  ")
	if err != nil {
		log.Printf("[Chain] Failed to marshal chain: %v", err)
		return
	}
	err = os.WriteFile(chainFile, data, 0644)
	if err != nil {
		log.Printf("[Chain] Failed to write chain file: %v", err)
	}
}

// LoadChain reads the chain from disk.
func LoadChain() {
	data, err := os.ReadFile(chainFile)
	if err != nil {
		log.Fatalf("[Chain] Failed to read chain file: %v", err)
	}
	err = json.Unmarshal(data, &Chain)
	if err != nil {
		log.Fatalf("[Chain] Failed to parse chain: %v", err)
	}
}

// TimestampNow returns the current time in RFC3339 format.
func TimestampNow() string {
	return time.Now().Format(time.RFC3339)
}

// GetChain returns the full chain (for APIs or inspection).
func GetChain() []Block {
	return Chain
}

// GetChainStats exposes summary stats for external queries.
func GetChainStats() map[string]interface{} {
	last := GetLastBlock()
	return map[string]interface{}{
		"length":        len(Chain),
		"latest_index":  last.Index,
		"latest_hash":   last.Hash,
		"latest_time":   last.Timestamp,
		"threats_total": countThreats(),
	}
}

// countThreats counts the number of transactions (threats) in the chain.
func countThreats() int {
	count := 0
	for _, b := range Chain {
		count += len(b.Transactions)
	}
	return count
}

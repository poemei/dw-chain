package blockchain

import (
	"encoding/json"
	"os"
)

type Blockchain struct {
	Blocks []Block
}

// Config holds dynamic runtime configuration
type Config struct {
	MinerThreads       int `json:"miner_threads"`
	BlockTargetSeconds int `json:"block_target_seconds"`
	DifficultyWindow   int `json:"difficulty_window"`
}

var AppConfig Config

func InitBlockchain() Blockchain {
	loadConfig()

	if _, err := os.Stat("data/chain.json"); err == nil {
		data, _ := os.ReadFile("data/chain.json")
		var chain Blockchain
		json.Unmarshal(data, &chain)
		return chain
	}

	genesis := Block{
		Index:        0,
		Timestamp:    Now(),
		Transactions: []Transaction{},
		PrevHash:     "0",
		Difficulty:   3,
		Miner:        "GENESIS",
	}
	MineBlock(&genesis, AppConfig.MinerThreads)
	chain := Blockchain{Blocks: []Block{genesis}}
	chain.Save()
	return chain
}

func (bc *Blockchain) Save() {
	data, _ := json.MarshalIndent(bc, "", "  ")
	os.WriteFile("data/chain.json", data, 0644)
}

func (bc *Blockchain) LatestBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func loadConfig() {
	file, err := os.ReadFile("data/config.json")
	if err != nil {
		AppConfig = Config{MinerThreads: 3, BlockTargetSeconds: 600, DifficultyWindow: 10}
		return
	}
	json.Unmarshal(file, &AppConfig)
}

func (bc *Blockchain) AdjustDifficulty() int {
	if len(bc.Blocks) < AppConfig.DifficultyWindow+1 {
		return bc.LatestBlock().Difficulty
	}

	start := bc.Blocks[len(bc.Blocks)-AppConfig.DifficultyWindow-1].Timestamp
	end := bc.LatestBlock().Timestamp
	actual := end - start
	expected := int64(AppConfig.BlockTargetSeconds * AppConfig.DifficultyWindow)

	if actual < expected/2 {
		return bc.LatestBlock().Difficulty + 1
	} else if actual > expected*2 {
		if bc.LatestBlock().Difficulty > 1 {
			return bc.LatestBlock().Difficulty - 1
		}
	}
	return bc.LatestBlock().Difficulty
}
package main

import (
	"fmt"
	"os"
	"time"

	"dw-chain/api"
	"dw-chain/blockchain"
)

const dataDir = "data/"

func main() {
	fmt.Println("?? Booting ThreatChain Miner Node...")

	chain := blockchain.InitBlockchain()
	blockchain.LoadPeers()

	// Start the API server
	api.StartServer(&chain)

	for {
		fmt.Println("? Checking for threats...")

		txs := blockchain.LoadTransactions()
		if len(txs) == 0 {
			fmt.Println("?? No pending threats. Sleeping...")
			time.Sleep(15 * time.Second)
			continue
		}

		newBlock := blockchain.Block{
			Index:        len(chain.Blocks),
			Timestamp:    blockchain.Now(),
			Threats:      txs,
			PrevHash:     chain.LatestBlock().Hash,
			Difficulty:   chain.AdjustDifficulty(),
			Miner:        "ThreatNode",
		}

		fmt.Println("??  Mining block...")
		blockchain.MineBlock(&newBlock)

		fmt.Printf("? Block mined: %s\n", newBlock.Hash)

		chain.Blocks = append(chain.Blocks, newBlock)
		chain.Save()

		// Clear the transaction queue
		os.WriteFile(dataDir+"transactions.json", []byte("[]"), 0644)

		fmt.Println("?? Block committed. Broadcasting...")
		blockchain.BroadcastBlock(newBlock)
		time.Sleep(15 * time.Second)
	}
}
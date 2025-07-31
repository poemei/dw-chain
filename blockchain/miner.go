// dw-chain/blockchain/miner.go

package blockchain

import (
	"log"
	"time"
)

// MinerLoop runs continuously to mine blocks when threats are in the pool
func MinerLoop() {
	log.Println("[Miner] Miner started...")

	for {
		// Wait if mempool is empty
		if len(TxPool.pool) == 0 {
			log.Println("[Miner] Mempool empty, sleeping...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Get the last block's timestamp
		lastBlock := GetLastBlock()
		if lastBlock.Timestamp == "" {
			log.Println("[Miner] Last block has no timestamp, using now.")
			lastBlock.Timestamp = TimestampNow()
		}

		parsedTime, err := time.Parse(time.RFC3339, lastBlock.Timestamp)
		if err != nil {
			log.Printf("[Miner] Time parse error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Check if we're mining too fast
		if time.Since(parsedTime) < blockInterval {
			log.Println("[Miner] Too early to mine, waiting...")
			time.Sleep(3 * time.Second)
			continue
		}

		// Fetch transactions and create new block
		txs := TxPool.FetchAll()
		newBlock := CreateBlock(txs)

		log.Printf("[Miner] Mining block #%d with difficulty %d", newBlock.Index, currentDifficulty)
		start := time.Now()
		newBlock = proofOfWork(newBlock, currentDifficulty)
		duration := time.Since(start)

		AddBlock(newBlock)
		log.Printf("[Miner] Block #%d mined in %s", newBlock.Index, duration)
	}
}

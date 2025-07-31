// dw-chain/blockchain/miner.go

package blockchain

import (
	"log"
	"strings"
	"time"
)

const (
	blockInterval      = 10 * time.Second
	initialDifficulty  = 3
	maxDifficulty      = 6
	difficultyWindow   = 5
)

// StartMining continuously checks for pending transactions and mines blocks
func StartMining() {
	log.Println("[Miner] Miner loop started...")

	for {
		if len(TxPool.Pending()) == 0 {
			log.Println("[Miner] Mempool empty, waiting...")
			time.Sleep(3 * time.Second)
			continue
		}

		lastBlock := Blockchain[len(Blockchain)-1]

		lastBlockTime, err := time.Parse(time.RFC3339, lastBlock.Timestamp)
		if err != nil {
			log.Printf("[Miner] Invalid timestamp on last block: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if time.Since(lastBlockTime) < blockInterval {
			time.Sleep(2 * time.Second)
			continue
		}

		transactions := TxPool.FetchAll()
		newBlock := CreateBlock(transactions)
		newBlock = ProofOfWork(newBlock, CalculateDifficulty())

		Blockchain = append(Blockchain, newBlock)
		log.Printf("[Miner] New block mined: #%d — %s", newBlock.Index, newBlock.Hash)
	}
}

// ProofOfWork runs simple leading-zero hash PoW
func ProofOfWork(block Block, difficulty int) Block {
	prefix := strings.Repeat("0", difficulty)

	for {
		hash := block.CalculateHash()
		if strings.HasPrefix(hash, prefix) {
			block.Hash = hash
			break
		}
		block.Nonce++
	}
	return block
}

// CalculateDifficulty adjusts based on previous N blocks
func CalculateDifficulty() int {
	chainLen := len(Blockchain)
	if chainLen <= difficultyWindow {
		return initialDifficulty
	}

	start := Blockchain[chainLen-difficultyWindow-1]
	end := Blockchain[chainLen-1]

	startTime, err1 := time.Parse(time.RFC3339, start.Timestamp)
	endTime, err2 := time.Parse(time.RFC3339, end.Timestamp)
	if err1 != nil || err2 != nil {
		log.Println("[Miner] Error parsing block timestamps, using default difficulty")
		return initialDifficulty
	}

	avgTime := endTime.Sub(startTime) / time.Duration(difficultyWindow)

	switch {
	case avgTime < blockInterval/2:
		return min(maxDifficulty, initialDifficulty+1)
	case avgTime > blockInterval*2:
		return max(1, initialDifficulty-1)
	default:
		return initialDifficulty
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
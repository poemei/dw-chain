package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
	Difficulty   int
	Miner        string
}

func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index) +
		strconv.FormatInt(block.Timestamp, 10) +
		block.PrevHash +
		strconv.Itoa(block.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

func MineBlock(block *Block, threads int) {
	start := time.Now()
	log := fmt.Sprintf("?? Mining block #%d with difficulty %d...\n", block.Index, block.Difficulty)
	appendLog(log)

	var wg sync.WaitGroup
	found := false
	var mu sync.Mutex

	prefix := strings.Repeat("0", block.Difficulty)

	for t := 0; t < threads; t++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			nonce := offset
			for !found {
				block.Nonce = nonce
				hash := CalculateHash(*block)

				if strings.HasPrefix(hash, prefix) {
					mu.Lock()
					if !found {
						found = true
						block.Hash = hash
						appendLog(fmt.Sprintf("? Block mined by thread %d at nonce %d\n", offset, nonce))
					}
					mu.Unlock()
					return
				}
				nonce += threads
			}
		}(t)
	}

	wg.Wait()
	duration := time.Since(start).Seconds()
	appendLog(fmt.Sprintf("?? Mining time: %.2f seconds | Hash: %s\n\n", duration, block.Hash))
}
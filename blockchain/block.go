// dw-chain/blockchain/block.go

package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// Block defines the structure of each block in the chain
type Block struct {
	Index        int           `json:"index"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	PrevHash     string        `json:"prev_hash"`
	Hash         string        `json:"hash"`
	Nonce        int           `json:"nonce"`
	MerkleRoot   string        `json:"merkle_root"`
}

// CalculateHash computes the SHA-256 hash of a block's core contents
func (b *Block) CalculateHash() string {
	record := string(b.Index) + b.Timestamp + b.PrevHash + b.MerkleRoot + string(b.Nonce)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

// GenerateMerkleRoot creates a simple Merkle Root hash from transactions
func GenerateMerkleRoot(txs []Transaction) string {
	if len(txs) == 0 {
		return ""
	}

	var hashes []string
	for _, tx := range txs {
		hashes = append(hashes, tx.Hash())
	}

	for len(hashes) > 1 {
		var temp []string
		for i := 0; i < len(hashes); i += 2 {
			if i+1 < len(hashes) {
				combined := hashes[i] + hashes[i+1]
				hash := sha256.Sum256([]byte(combined))
				temp = append(temp, hex.EncodeToString(hash[:]))
			} else {
				// Duplicate last hash if odd
				hash := sha256.Sum256([]byte(hashes[i] + hashes[i]))
				temp = append(temp, hex.EncodeToString(hash[:]))
			}
		}
		hashes = temp
	}
	return hashes[0]
}

// HashMatchesDifficulty checks if a block hash satisfies difficulty criteria
func HashMatchesDifficulty(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
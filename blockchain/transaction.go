// dw-chain/blockchain/transaction.go

package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Transaction defines a threat record structure in the DW-Chain.
type Transaction struct {
	Type      string `json:"type"`
	IP        string `json:"ip"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

// Hash generates a SHA-256 hash for a transaction (used in Merkle root).
func (tx Transaction) Hash() string {
	record := fmt.Sprintf("%s|%s|%s|%s", tx.Type, tx.IP, tx.Reason, tx.Timestamp)
	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

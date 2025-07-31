// dw-chain/blockchain/transaction.go

package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Transaction defines a basic threat report structure for the chain.
type Transaction struct {
	Type      string `json:"type"`
	IP        string `json:"ip"`
	Reason    string `json:"reason"`
	Timestamp string `json:"timestamp"`
}

// TimestampNow returns the current UTC timestamp in RFC3339 format.
func TimestampNow() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Hash generates a SHA-256 hash of the transaction's fields.
func (t Transaction) Hash() string {
	data := t.Type + t.IP + t.Reason + t.Timestamp
	sum := sha256.Sum256([]byte(data))
	return hex.EncodeToString(sum[:])
}

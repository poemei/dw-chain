// tx_pool.go
package blockchain

import (
	"encoding/json"
	"os"
	"time"
)

//const dataDir = "data/"

var ThreatPool []ThreatTransaction


func LoadTransactions() []Transaction {
	data, _ := os.ReadFile(dataDir+"transactions.json")
	var txs []Transaction
	json.Unmarshal(data, &txs)
	return txs
}

func SaveTransactions(txs []Transaction) {
	data, _ := json.MarshalIndent(txs, "", "  ")
	os.WriteFile(dataDir+"transactions.json", data, 0644)
}

func NewTransaction(ip string, reason string, timestamp string) Transaction {
	// If no timestamp provided, generate one
	if timestamp == "" {
		timestamp = time.Now().Format(time.RFC3339)
	}

	return Transaction{
		IP:        ip,
		Reason:    reason,
		Timestamp: timestamp,
	}
}

// AddThreat queues a new threat to the mempool.
func AddThreat(tx ThreatTransaction) {
	ThreatPool = append(ThreatPool, tx)
}

// FetchThreats returns a copy of the threat pool and resets it.
func FetchThreats() []ThreatTransaction {
	txs := ThreatPool
	ThreatPool = []ThreatTransaction{}
	return txs
}
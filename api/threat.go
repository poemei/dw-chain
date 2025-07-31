// dw-chain/api/threat.go

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"dw-chain/blockchain"
)

func ThreatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var tx blockchain.Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if tx.Timestamp == "" {
		tx.Timestamp = time.Now().Format(time.RFC3339)
	}

	blockchain.TxPool.Add(tx)
	log.Printf("[Threat] New threat queued: %+v", tx)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "threat accepted"})
}

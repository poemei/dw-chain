// dw-chain/api/router.go

package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"dw-chain/blockchain"
)

// InitRouter sets up the HTTP routing for the blockchain API.
func InitRouter() {
	http.HandleFunc("/threat", ThreatHandler)
	http.HandleFunc("/chain/stats", ChainStatsHandler)
	http.HandleFunc("/peers", PeersHandler)
	log.Println("[API] Router initialized")
}

// ThreatHandler receives and queues threat data to the mempool
func ThreatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx blockchain.Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, "Invalid transaction payload", http.StatusBadRequest)
		return
	}

	if tx.Timestamp == "" {
		tx.Timestamp = blockchain.TimestampNow()
	}

	blockchain.TxPool.Add(tx)
	log.Printf("[TX] Queued transaction: %+v", tx)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Threat queued"))
}

// ChainStatsHandler returns basic metadata about the current chain state
func ChainStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := blockchain.GetChainStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// PeersHandler handles peer listing and registration
func PeersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		peers := blockchain.GetPeers()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(peers)

	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			http.Error(w, "Invalid peer payload", http.StatusBadRequest)
			return
		}

		address := strings.TrimSpace(string(body))
		if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
			http.Error(w, "Peer must begin with http:// or https://", http.StatusBadRequest)
			return
		}

		blockchain.AddPeer(address)
		timestamp := time.Now().Format(time.RFC3339)
		log.Printf("[PEER] %s registered at %s", address, timestamp)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Peer registered"))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// dw-chain/main.go

package main

import (
	"log"
	"net/http"

	"dw-chain/api"
	"dw-chain/blockchain"
)

func main() {
	log.Println("[Main] Starting DW-Chain (Threat Index)")

	// Initialize blockchain and mempool
	blockchain.InitChain()

	// Initialize HTTP router
	api.InitRouter()

	log.Println("[Main] Listening on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("[Main] Server error: %v", err)
	}
}
//dw-chain/main.go
package main

import (
	"log"
	"net/http"

	"dw-chain/api"
)

func main() {
	log.Println("[Core] dw-chain threat index starting...")
	api.InitRouter()

	log.Println("[Core] Server listening on :8333")
	err := http.ListenAndServe(":8333", nil)
	if err != nil {
		log.Fatalf("[Core] Failed to start server: %v", err)
	}
}
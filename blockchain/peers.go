// dw-chain/blockchain/peers.go

package blockchain

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var peerLock sync.Mutex
var peerList []string

const peersFile = "data/peers.json"

// LoadPeers initializes the peer list from file
func LoadPeers() {
	peerLock.Lock()
	defer peerLock.Unlock()

	data, err := os.ReadFile(peersFile)
	if err != nil {
		log.Println("[Peers] No existing peer file found")
		return
	}

	err = json.Unmarshal(data, &peerList)
	if err != nil {
		log.Println("[Peers] Failed to parse peers.json:", err)
	}
}

// SavePeers writes the current peer list to file
func SavePeers() {
	peerLock.Lock()
	defer peerLock.Unlock()

	data, err := json.MarshalIndent(peerList, "", "  ")
	if err != nil {
		log.Println("[Peers] Failed to encode peers:", err)
		return
	}

	err = os.WriteFile(peersFile, data, 0644)
	if err != nil {
		log.Println("[Peers] Failed to save peers:", err)
	}
}

// AddPeer registers a new peer if it's not already present
func AddPeer(address string) {
	peerLock.Lock()
	defer peerLock.Unlock()

	for _, p := range peerList {
		if p == address {
			return
		}
	}
	peerList = append(peerList, address)
	log.Println("[Peers] Added new peer:", address)
	SavePeers()
}

// GetPeers returns a copy of the peer list
func GetPeers() []string {
	peerLock.Lock()
	defer peerLock.Unlock()

	copied := make([]string, len(peerList))
	copy(copied, peerList)
	return copied
}

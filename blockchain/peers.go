package blockchain

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "bytes"
)

const PeersFile = "data/peers.json"

var KnownPeers []string

func LoadPeers() {
    data, err := os.ReadFile(PeersFile)
    if err != nil {
        KnownPeers = []string{}
        return
    }
    json.Unmarshal(data, &KnownPeers)
}

func SavePeers() {
    data, _ := json.MarshalIndent(KnownPeers, "", "  ")
    os.WriteFile(PeersFile, data, 0644)
}

func AddPeer(url string) {
    for _, p := range KnownPeers {
        if p == url {
            return
        }
    }
    KnownPeers = append(KnownPeers, url)
    SavePeers()
}

func BroadcastBlock(block Block) {
    for _, peer := range KnownPeers {
        go func(p string) {
            jsonData, _ := json.Marshal(block)
            resp, err := http.Post(p+"/block", "application/json", bytes.NewReader(jsonData))
            if err != nil {
                fmt.Printf("? Failed to broadcast to %s: %v\n", p, err)
                return
            }
            defer resp.Body.Close()
            body, _ := io.ReadAll(resp.Body)
            fmt.Printf("?? Sent block to %s: %s\n", p, string(body))
        }(peer)
    }
}
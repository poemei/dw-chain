//dw-chain/api/server.go

package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "log"
    "time"
    "os"
    "dw-chain/blockchain"
)

type Status struct {
    LatestBlock      blockchain.Block            `json:"latest_block"`
    ChainLength      int                         `json:"chain_length"`
    CurrentBalance   blockchain.Balances         `json:"balances"`
    TransactionPool  []blockchain.Transaction    `json:"pending_transactions"`
    Config           blockchain.Config           `json:"config"`
    Heartbeat        string                      `json:"heartbeat"`
}

var chainRef *blockchain.Blockchain
var balancesRef *blockchain.Balances


func StartServer(chain *blockchain.Blockchain, balances *blockchain.Balances) {
    chainRef = chain
    balancesRef = balances

    http.HandleFunc("/status", handleStatus)
    http.HandleFunc("/chain", handleChain)
    http.HandleFunc("/block", handleBlock)
    http.HandleFunc("/peers", handlePeers)
    http.HandleFunc("/threat", handleThreat)

    // http file server
    //http.ListenAndServe(":8081", http.FileServer(http.Dir(".")))

    go func() {
        fmt.Println("?? API server running at http://localhost:8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            fmt.Println("? API server error:", err)
        }
    }()
}

// GET /status
func handleStatus(w http.ResponseWriter, r *http.Request) {
    if chainRef == nil || balancesRef == nil {
        http.Error(w, "Chain or balances not available", http.StatusInternalServerError)
        return
    }

    status := Status{
        LatestBlock:     chainRef.LatestBlock(),
        ChainLength:     len(chainRef.Blocks),
        CurrentBalance:  *balancesRef,
        TransactionPool: blockchain.LoadTransactions(),
        Config:          blockchain.AppConfig,
        Heartbeat:       "alive",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

// GET /chain
func handleChain(w http.ResponseWriter, r *http.Request) {
    if chainRef == nil {
        http.Error(w, "Chain not available", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(chainRef.Blocks)
}

// POST /block
func handleBlock(w http.ResponseWriter, r *http.Request) {
    var incomingBlock blockchain.Block
    err := json.NewDecoder(r.Body).Decode(&incomingBlock)
    if err != nil {
        http.Error(w, "Invalid block data", http.StatusBadRequest)
        return
    }

    last := chainRef.LatestBlock()

    if incomingBlock.Index == last.Index+1 && incomingBlock.PrevHash == last.Hash {
        chainRef.Blocks = append(chainRef.Blocks, incomingBlock)
        chainRef.Save()
        fmt.Println("?? Block accepted from peer.")
        w.WriteHeader(http.StatusOK)
    } else {
        http.Error(w, "Block rejected (invalid index or hash)", http.StatusConflict)
    }
}

// GET or POST /peers
func handlePeers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(blockchain.KnownPeers)
    case http.MethodPost:
        var peer struct {
            URL string `json:"url"`
        }
        if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
            http.Error(w, "Invalid JSON", http.StatusBadRequest)
            return
        }
        blockchain.AddPeer(peer.URL)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Peer added."))
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// Threats
var threatsFile = "data/threats.json"

func handleThreat(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
        return
    }

    var t blockchain.Threat
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    if t.Timestamp == "" {
        t.Timestamp = time.Now().Format(time.RFC3339)
    }

    log.Printf("[THREAT] IP: %s | Reason: %s | Timestamp: %s", t.IP, t.Reason, t.Timestamp)
    saveThreat(t)

    w.WriteHeader(http.StatusAccepted)
    w.Write([]byte("Threat received"))
}

func saveThreat(t blockchain.Threat) {
    var threats []blockchain.Threat

    data, err := os.ReadFile(threatsFile)
    if err == nil {
        json.Unmarshal(data, &threats)
    }

    threats = append(threats, t)
    out, _ := json.MarshalIndent(threats, "", "  ")
    os.WriteFile(threatsFile, out, 0644)
}

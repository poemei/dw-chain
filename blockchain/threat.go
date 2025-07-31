// dw-chain/blockchain/threat.go
package blockchain

type Threat struct {
    IP        string `json:"ip"`
    Reason    string `json:"reason"`
    Timestamp string `json:"timestamp"`
}
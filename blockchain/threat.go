// dw-chain/blockchain/threat.go
package blockchain

type Threat struct {
    IP        string `json:"ip"`
    Reason    string `json:"reason"`
    Timestamp string `json:"timestamp"`
}

// ThreatTransaction defines the structure of a threat report.
type ThreatTransaction struct {
	IP        string `json:"ip"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

// NewThreat creates a new ThreatTransaction with the current timestamp.
func NewThreat(ip, reason string) ThreatTransaction {
	return ThreatTransaction{
		IP:        ip,
		Reason:    reason,
		Timestamp: time.Now().Unix(),
	}
}
// dw-chain/blockchain/tx_pool.go

package blockchain

import "sync"

// txPool holds unconfirmed transactions before they're mined.
type txPool struct {
	mu    sync.Mutex
	pool  []Transaction
}

var TxPool = &txPool{}

// Add inserts a transaction into the pool.
func (tp *txPool) Add(tx Transaction) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.pool = append(tp.pool, tx)
}

// All returns all pending transactions.
func (tp *txPool) All() []Transaction {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	return append([]Transaction(nil), tp.pool...) // returns a copy
}

// Clear empties the transaction pool after mining.
func (tp *txPool) Clear() {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.pool = nil
}

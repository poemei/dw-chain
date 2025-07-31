// dw-chain/blockchain/tx_pool.go

package blockchain

import "sync"

// txPool manages the in-memory pool of threat transactions waiting to be mined.
type txPool struct {
	mu    sync.Mutex
	pool  []Transaction
}

// TxPool is the global transaction pool used by miner and API.
var TxPool = &txPool{}

// Add inserts a transaction into the mempool.
func (tp *txPool) Add(tx Transaction) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.pool = append(tp.pool, tx)
}

// FetchAll returns a copy of all transactions currently in the pool.
func (tp *txPool) FetchAll() []Transaction {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	clone := make([]Transaction, len(tp.pool))
	copy(clone, tp.pool)
	return clone
}

// Clear empties the transaction pool after block creation.
func (tp *txPool) Clear() {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.pool = []Transaction{}
}

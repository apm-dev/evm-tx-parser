package repo

import (
	"sync"

	"github.com/apm-dev/evm-tx-parser/src/domain"
)

type transactionRepo struct {
	lock         *sync.RWMutex
	data         map[string]domain.Transaction
	addressIndex map[string][]domain.Transaction
}

func NewTransactionRepo() domain.TransactionRepo {
	return &transactionRepo{
		lock:         &sync.RWMutex{},
		data:         make(map[string]domain.Transaction),
		addressIndex: make(map[string][]domain.Transaction),
	}
}

func (r *transactionRepo) SaveMany(txs []domain.Transaction) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, tx := range txs {
		// skip if tx exists
		if _, ok := r.data[tx.Hash]; !ok {
			r.data[tx.Hash] = tx
			// index tx for it's `from` address
			if _, ok := r.addressIndex[tx.From]; !ok {
				r.addressIndex[tx.From] = make([]domain.Transaction, 0)
			}
			r.addressIndex[tx.From] = append(r.addressIndex[tx.From], tx)
			// index tx for it's `to` address
			if _, ok := r.addressIndex[tx.To]; !ok {
				r.addressIndex[tx.To] = make([]domain.Transaction, 0)
			}
			r.addressIndex[tx.To] = append(r.addressIndex[tx.To], tx)
		}
	}
	return nil
}

func (r *transactionRepo) FindByAddress(address string) ([]domain.Transaction, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.addressIndex[address], nil
}

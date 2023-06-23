package repo

import (
	"sync"

	"github.com/apm-dev/evm-tx-parser/src/domain"
)

type addressRepo struct {
	lock *sync.RWMutex
	data map[string]struct{}
}

func NewAddressRepo() domain.AddressRepo {
	return &addressRepo{
		lock: &sync.RWMutex{},
		data: make(map[string]struct{}),
	}
}

func (r *addressRepo) Save(address string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data[address] = struct{}{}
	return nil
}

func (r *addressRepo) Exist(address string) bool {
	r.lock.RLock()
	defer r.lock.RUnlock()
	_, ok := r.data[address]
	return ok
}

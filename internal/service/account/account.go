package account

import (
	"sync"

	"github.com/k1nky/gophermart/internal/entity/user"
)

const (
	DefMaxRows = 100
)

type Service struct {
	store            storage
	transactionLocks map[user.ID]sync.Mutex
}

func New(store storage) *Service {
	return &Service{
		store:            store,
		transactionLocks: make(map[user.ID]sync.Mutex),
	}
}

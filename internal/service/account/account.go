package account

import (
	"context"
	"errors"
	"os/user"

	"github.com/k1nky/gophermart/internal/entity/order"
)

//go:generate mockgen -source=account.go -destination=mock/storage.go -package=mock Storage
type Storage interface {
	NewOrder(ctx context.Context, u user.User, o order.Order) (*order.Order, error)
	IsUniqueViolation(err error) bool
}

type Service struct {
	store Storage
}

var (
	ErrDuplicateOrderError = errors.New("order already exists")
)

func New(store Storage) *Service {
	return &Service{
		store: store,
	}
}

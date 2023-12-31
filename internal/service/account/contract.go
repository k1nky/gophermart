package account

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/k1nky/gophermart/internal/entity/withdraw"
)

//go:generate mockgen -source=contract.go -destination=mock/storage.go -package=mock storage
type storage interface {
	NewOrder(ctx context.Context, newOrder order.Order) (*order.Order, error)
	GetOrderByNumber(ctx context.Context, number order.OrderNumber) (*order.Order, error)
	GetOrdersByUserID(ctx context.Context, userID user.ID, maxRows uint) ([]*order.Order, error)
	GetBalanceByUser(ctx context.Context, userID user.ID) (user.Balance, error)
	GetWithdrawalsByUserID(ctx context.Context, userID user.ID, maxRows uint) ([]*withdraw.Withdraw, error)
	NewWithdraw(ctx context.Context, w withdraw.Withdraw) (*withdraw.Withdraw, error)
}

type logger interface {
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

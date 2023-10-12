package account

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

//go:generate mockgen -source=contract.go -destination=mock/storage.go -package=mock storage
type storage interface {
	NewOrder(ctx context.Context, newOrder order.Order) (*order.Order, error)
	GetOrderByNumber(ctx context.Context, number order.OrderNumber) (*order.Order, error)
	GetOrdersByUserID(ctx context.Context, userID user.ID) ([]*order.Order, error)
}

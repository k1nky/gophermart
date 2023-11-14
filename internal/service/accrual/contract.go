package accural

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
)

type logger interface {
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type store interface {
	GetOrdersByStatus(ctx context.Context, statuses []order.OrderStatus, maxRows uint) ([]*order.Order, error)
	UpdateOrder(ctx context.Context, o order.Order) error
}

type orderAccrual interface {
	FetchOrder(ctx context.Context, number order.OrderNumber) (*order.Order, error)
}

package http

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

//go:generate mockgen -source=contract.go -destination=mock/auth.go -package=mock authService
type authService interface {
	Register(ctx context.Context, u user.User) (string, error)
	Login(ctx context.Context, u user.User) (string, error)
	Authorize(token string) (user.PrivateClaims, error)
}

type accountService interface {
	NewOrder(ctx context.Context, o order.Order) (*order.Order, error)
	GetUserOrders(ctx context.Context, userID user.ID) ([]*order.Order, error)
}

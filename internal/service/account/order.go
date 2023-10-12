package account

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

func (s *Service) NewOrder(ctx context.Context, newOrder order.Order) (*order.Order, error) {
	o, err := s.store.GetOrderByNumber(ctx, newOrder.Number)
	if err != nil {
		return nil, err
	}
	if o != nil {
		if o.UserID != newOrder.UserID {
			return nil, order.ErrOrderBelongsToAnotherUser
		}
		return nil, order.ErrDuplicateOrder
	}
	o, err = s.store.NewOrder(ctx, newOrder)
	return o, err
}

func (s *Service) GetUserOrders(ctx context.Context, userID user.ID) ([]*order.Order, error) {
	orders, err := s.store.GetOrdersByUserID(ctx, userID)
	return orders, err
}

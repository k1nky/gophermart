package account

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
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

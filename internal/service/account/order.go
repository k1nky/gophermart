package account

import (
	"context"
	"fmt"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

// Регистрирует новый заказ
func (s *Service) NewOrder(ctx context.Context, newOrder order.Order) (*order.Order, error) {
	fail := func(err error) (*order.Order, error) {
		wrapped := fmt.Errorf("account: new order: %w", err)
		s.log.Errorf("%s", wrapped.Error())
		return nil, wrapped
	}
	o, err := s.store.GetOrderByNumber(ctx, newOrder.Number)
	if err != nil {
		return fail(err)
	}
	if o != nil {
		if o.UserID != newOrder.UserID {
			return nil, order.ErrBelongsToAnotherUser
		}
		return nil, order.ErrDuplicated
	}
	o, err = s.store.NewOrder(ctx, newOrder)
	if err != nil {
		return fail(err)
	}
	return o, nil
}

// Возвращает спикок заказов пользователя
func (s *Service) GetUserOrders(ctx context.Context, userID user.ID) ([]*order.Order, error) {
	orders, err := s.store.GetOrdersByUserID(ctx, userID, DefaultMaxRows)
	if err != nil {
		wrapped := fmt.Errorf("account: get new orders: %w", err)
		s.log.Errorf("%s", wrapped)
	}
	return orders, err
}

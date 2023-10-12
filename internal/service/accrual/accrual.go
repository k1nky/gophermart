package accural

import (
	"context"
	"time"

	"github.com/k1nky/gophermart/internal/entity/order"
)

type Store interface {
	GetOrdersByStatus(ctx context.Context, statuses []order.OrderStatus) ([]*order.Order, error)
	UpdateOrder(ctx context.Context, o order.Order) error
}

type OrderAccrual interface {
	FetchOrder(ctx context.Context, number order.OrderNumber) (*order.Order, error)
}

type Service struct {
	store        Store
	orderAccrual OrderAccrual
}

func New(store Store, orderAccrual OrderAccrual) *Service {
	return &Service{
		store:        store,
		orderAccrual: orderAccrual,
	}
}

func (s *Service) getNewOrders(ctx context.Context) <-chan *order.Order {
	ordersCh := make(chan *order.Order, 2)
	go func() {
		t := time.NewTicker(5 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				orders, err := s.store.GetOrdersByStatus(ctx, []order.OrderStatus{order.StatusNew, order.StatusProcessing})
				if err != nil {
					// TODO: handle error
					continue
				}
				for _, o := range orders {
					ordersCh <- o
				}
			case <-ctx.Done():
				close(ordersCh)
				return
			}
		}
	}()
	return ordersCh
}

func (s *Service) updateOrder(ctx context.Context, orders <-chan *order.Order) <-chan *order.Order {
	ordersCh := make(chan *order.Order, 2)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ordersCh)
				return
			case o := <-orders:
				got, err := s.orderAccrual.FetchOrder(ctx, o.Number)
				if err != nil {
					// TODO: handle error
					continue
				}
				if got == nil {
					continue
				}
				if got.Status == order.StatusRegistered {
					got.Status = order.StatusProcessing
				}
				if o.Status != got.Status {
					o.Status = got.Status
					if o.Status == order.StatusProcessed {
						o.Accrual = got.Accrual
					}
					ordersCh <- o
				}
			}
		}
	}()
	return ordersCh
}

func (s *Service) Process(ctx context.Context) {
	go func() {

		for o := range s.updateOrder(ctx, s.getNewOrders(ctx)) {
			if err := s.store.UpdateOrder(ctx, *o); err != nil {
				//
				continue
			}
		}
	}()
}

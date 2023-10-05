package accural

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
)

type Store interface {
	GetNewOrders() ([]*order.Order, error)
}

type OrderAccrual interface {
	Fetch(number order.OrderNumber) (*order.Order, error)
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

func (s *Service) Process(ctx context.Context) {
	ordersCh := make(chan *order.Order, 2)
	go s.processOrder(ordersCh)
	orders, _ := s.store.GetNewOrders()
	for _, o := range orders {
		ordersCh <- o
	}

}

func (s *Service) processOrder(orders <-chan *order.Order) {
	order := <-orders
	s.orderAccrual.Fetch(order.Number)
}

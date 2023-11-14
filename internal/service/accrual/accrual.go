package accural

import (
	"context"
	"time"

	"github.com/k1nky/gophermart/internal/entity/order"
)

const (
	// обрабатываем за раз не более DefaultMaxRows заказов
	DefaultMaxRows = 100
	// максимальный размер очереди заказов на проверку начислений
	DefaultMaxOrderQueueSize = 2
	// интервал обновления заказов
	DefaultUpdateInterval = 5 * time.Second
)

type Service struct {
	store        store
	orderAccrual orderAccrual
	log          logger
}

func New(store store, orderAccrual orderAccrual, l logger) *Service {
	return &Service{
		log:          l,
		orderAccrual: orderAccrual,
		store:        store,
	}
}

func (s *Service) getNewOrders(ctx context.Context) <-chan *order.Order {
	ordersCh := make(chan *order.Order, DefaultMaxOrderQueueSize)
	go func() {
		t := time.NewTicker(DefaultUpdateInterval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				orders, err := s.store.GetOrdersByStatus(ctx, []order.OrderStatus{order.StatusNew, order.StatusProcessing}, DefaultMaxRows)
				s.log.Debugf("accrual: got %d new orders", len(orders))
				if err != nil {
					s.log.Errorf("accrual: %v", err)
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
	ordersCh := make(chan *order.Order, DefaultMaxOrderQueueSize)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ordersCh)
				return
			case o := <-orders:
				s.log.Debugf("accrual: fetch order #%s", o.Number)
				got, err := s.orderAccrual.FetchOrder(ctx, o.Number)
				if err != nil {
					s.log.Errorf("accrual: failed fetching order #%s: %v", o.Number, err)
					continue
				}
				if got == nil {
					// заказ не зарегистрирован в системе начислений
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
		// У сервиса accrual есть ограничение по количеству запросов. Адаптер этого сервиса умеет
		// повторять запрос с ожидаением Retry-After. В этом случае getNewOrders также будет ожидать и
		// не добавлять в очередь новые запросы для проверки начислений.
		for o := range s.updateOrder(ctx, s.getNewOrders(ctx)) {
			if err := s.store.UpdateOrder(ctx, *o); err != nil {
				s.log.Errorf("accrual: poll order #%s: %v", o.Number, err)
				continue
			}
		}
	}()
}

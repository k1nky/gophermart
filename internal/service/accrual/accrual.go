package accural

import (
	"context"
	"time"

	"github.com/k1nky/gophermart/internal/entity/order"
)

type Store interface {
	GetOrdersByFilter(ctx context.Context, filter string, args ...interface{}) ([]*order.Order, error)
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

func (s *Service) Process(ctx context.Context) {
	t := time.NewTicker(15 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			// - `PROCESSING` — вознаграждение за заказ рассчитывается;
			// - `INVALID` — система расчёта вознаграждений отказала в расчёте;
			// - `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.

			orders, err := s.store.GetOrdersByFilter(ctx, "status IN ($1, $2)", order.StatusNew, order.StatusProcessing)
			if err != nil {
				// TODO: handle err
			}
			for _, o := range orders {
				got, err := s.orderAccrual.FetchOrder(ctx, o.Number)
				if err != nil {
					// TODO: handle err
				}
				if o.Status != got.Status {
					o.Status = got.Status
					if o.Status == order.StatusProcessed {
						o.Accrual = got.Accrual
					}
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

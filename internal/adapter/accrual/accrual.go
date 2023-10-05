package accrual

import (
	"github.com/go-resty/resty/v2"
	"github.com/k1nky/gophermart/internal/entity/order"
)

type Adapter struct {
}

func (a *Adapter) FetchOrder() *order.Order {
	cli := resty.New()
	cli.R().Get("<accrual>")
	return &order.Order{}
}

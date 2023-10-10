package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/k1nky/gophermart/internal/entity/order"
)

const (
	DefRequestTimeout = 5 * time.Second
)

var (
	ErrUnexpectedResponse = errors.New("unexpected reposonse")
)

type orderResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *uint  `json:"accrual,omitempty"`
}

type Adapter struct {
	cli *resty.Client
	url string
}

func New(url string) *Adapter {
	return &Adapter{
		url: url,
		cli: resty.New().SetTimeout(DefRequestTimeout),
	}
}

func (a *Adapter) newRequest() *resty.Request {
	return a.cli.R()
}

// Заказ может быть взят в расчёт в любой момент после его совершения.
// Время выполнения расчёта системой не регламентировано. Статусы `INVALID` и `PROCESSED` являются окончательными.
// Общее количество запросов информации о начислении не ограничено.
func (a *Adapter) FetchOrder(ctx context.Context, number order.OrderNumber) (*order.Order, error) {
	req := a.newRequest()
	// Получение информации о расчёте начислений баллов лояльности.
	url, err := url.JoinPath(a.url, "/api/order", string(number))
	if err != nil {
		return nil, err
	}
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		responseData := orderResponse{}
		if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
			return nil, err
		}
		o := order.Order{
			Number:  order.OrderNumber(responseData.Order),
			Accrual: responseData.Accrual,
			Status:  order.OrderStatus(responseData.Status),
		}
		return &o, nil
	case http.StatusNoContent:
		// заказ не зарегистрирован в системе расчета.
		return nil, nil
	case http.StatusTooManyRequests:
		// TODO: TooManyRequest, place retrier here
		// превышено количество запросов к сервису.
		return nil, nil
	}
	return nil, ErrUnexpectedResponse
}

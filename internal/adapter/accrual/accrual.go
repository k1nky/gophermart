package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/k1nky/gophermart/internal/entity/order"
)

const (
	DefRequestTimeout   = 5 * time.Second
	DefRetryCount       = 1
	DefRetryMaxWaitTime = 10 * time.Second
)

var (
	ErrUnexpectedResponse = errors.New("unexpected response")
)

type orderResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float32 `json:"accrual,omitempty"`
}

type Adapter struct {
	cli *resty.Client
	url string
}

func New(url string) *Adapter {
	return &Adapter{
		url: url,
		cli: resty.New().
			SetTimeout(DefRequestTimeout).
			SetRetryCount(DefRetryCount).
			SetRetryMaxWaitTime(DefRetryMaxWaitTime).
			SetRetryAfter(func(c *resty.Client, r *resty.Response) (time.Duration, error) {
				// если в ответе был заголовок Retry-After в формате <seconds>
				// то постараемся придерживаться его, но не дольше MaxWaitTime
				s := r.Header().Get("Retry-After")
				if s == "" {
					return 0, nil
				}
				sec, err := strconv.Atoi(s)
				if err != nil {
					return 0, nil
				}
				return time.Duration(sec * int(time.Second)), nil
			}),
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
	req.AddRetryCondition(func(response *resty.Response, err error) bool {
		// если на запрос придет Too many requests, то сделаем retry в соответствии с настройками клиента
		// описанными в конструкторе
		return response.StatusCode() == http.StatusTooManyRequests
	})
	// Получение информации о расчёте начислений баллов лояльности.
	url, err := url.JoinPath(a.url, "/api/orders", string(number))
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
	}
	return nil, fmt.Errorf("%d %s: %w", resp.StatusCode(), resp.Body(), ErrUnexpectedResponse)
}

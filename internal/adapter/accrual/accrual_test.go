package accrual

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/stretchr/testify/assert"
)

func TestTooManyRequestWithRetryAfter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Retry-After", "3")
		rw.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()
	c := New(ts.URL)
	_, err := c.FetchOrder(context.TODO(), "1")
	assert.ErrorIs(t, err, ErrUnexpectedResponse)
}

func TestTooManyRequestWithoutRetryAfter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()
	c := New(ts.URL)
	_, err := c.FetchOrder(context.TODO(), "1")
	assert.ErrorIs(t, err, ErrUnexpectedResponse)
}

func TestTooManyRequestSuccess(t *testing.T) {
	isFirst := true
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if isFirst {
			rw.WriteHeader(http.StatusTooManyRequests)
			isFirst = false
		} else {
			rw.Write([]byte(`{"order":"1", "status":"INVALID"}`))
			rw.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()
	c := New(ts.URL)
	_, err := c.FetchOrder(context.TODO(), "1")
	assert.NoError(t, err)
}

func TestNoContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	c := New(ts.URL)
	o, err := c.FetchOrder(context.TODO(), "1")
	assert.NoError(t, err)
	assert.Nil(t, o)
}

func TestFetchOrderWithAccrual(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"order":"1", "status":"PROCESSED", "accrual": 123.0}`))
		rw.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c := New(ts.URL)
	o, err := c.FetchOrder(context.TODO(), "1")
	assert.NoError(t, err)
	v := float32(123.0)
	assert.Equal(t, &order.Order{
		Number:  "1",
		Status:  order.StatusProcessed,
		Accrual: &v,
	}, o)
}

func TestFetchOrderWithoutAccrual(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"order":"1", "status":"PROCESSED"}`))
		rw.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	c := New(ts.URL)
	o, err := c.FetchOrder(context.TODO(), "1")
	assert.NoError(t, err)
	assert.Equal(t, &order.Order{
		Number:  "1",
		Status:  order.StatusProcessed,
		Accrual: nil,
	}, o)
}

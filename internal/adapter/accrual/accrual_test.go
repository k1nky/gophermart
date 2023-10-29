package accrual

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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

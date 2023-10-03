package order

import (
	"time"
)

type OrderStatus string
type OrderNumber string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

//go:generate easyjson order.go
//easyjson:json
type Order struct {
	Number     OrderNumber `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *uint       `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func (n OrderNumber) IsValid() bool {
	l := len(n)
	digits := make([]byte, 0, l)
	for i := 0; i < len(n); i++ {
		digits = append(digits, byte(n[i])-48)
	}
	sum := 0
	for i := 0; i < l; i++ {
		d := digits[l-1-i]
		if i%2 == 1 {
			d = d * 2
			if d > 9 {
				d = d - 9
			}
		}
		sum = sum + int(d)
	}
	return (sum % 10) == 0
}

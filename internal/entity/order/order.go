package order

import (
	"strconv"
	"time"
)

type OrderStatus string
type OrderNumber string
type ID uint64

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

//go:generate easyjson order.go
//easyjson:json
type Order struct {
	ID         ID
	Number     OrderNumber `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *uint       `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func (n OrderNumber) IsValid() bool {
	digits := make([]int, 0, len(n))
	for _, c := range n {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			return false
		}
		digits = append(digits, v)
	}
	checksum := luhnChecksum(digits)
	return (checksum % 10) == 0
}

func luhnChecksum(digits []int) int {
	sum := 0
	l := len(digits)
	for i := 1; i <= l; i++ {
		d := digits[l-i]
		if i%2 == 0 {
			d = d * 2
			if d > 9 {
				d = d - 9
			}
		}
		sum = sum + int(d)
	}
	return sum
}

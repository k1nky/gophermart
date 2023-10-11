package order

import "errors"

var (
	ErrDuplicateOrderError = errors.New("order already exists")
)

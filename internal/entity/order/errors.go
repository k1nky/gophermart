package order

import "errors"

var (
	ErrDuplicated           = errors.New("order already exists")
	ErrInvalidNumberFormat  = errors.New("invalid order number format")
	ErrBelongsToAnotherUser = errors.New("order belongs to another user")
	ErrAlreadyProcessed     = errors.New("order has already processed")
)

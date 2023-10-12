package order

import "errors"

var (
	ErrDuplicateOrder            = errors.New("order already exists")
	ErrInvalidOrderNumberFormat  = errors.New("invalid order number format")
	ErrOrderBelongsToAnotherUser = errors.New("order belongs to another user")
)

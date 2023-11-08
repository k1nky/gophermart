package database

import (
	"fmt"
)

type ExecutingQueryError struct {
	Err error
}

func NewExecutingQueryError(err error) error {
	return &ExecutingQueryError{
		Err: err,
	}
}

func (e *ExecutingQueryError) Error() string {
	return fmt.Sprintf("failed executing query: %v", e.Err)
}

func (e *ExecutingQueryError) Unwrap() error {
	return e.Err
}

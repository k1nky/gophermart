package user

import "errors"

var (
	ErrDuplicateLoginError = errors.New("login already exists")
	ErrInvalidCredentials  = errors.New("login or password is not correct")
)

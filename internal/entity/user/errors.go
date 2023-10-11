package user

import "errors"

var (
	ErrDuplicateLogin           = errors.New("login already exists")
	ErrInvalidCredentials       = errors.New("login or password is not correct")
	ErrUnathorized              = errors.New("user is not authorized")
	ErrCredentialsInvalidFormat = errors.New("login or password has invalid format")
)

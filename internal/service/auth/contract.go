package auth

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/user"
)

//go:generate mockgen -source=contract.go -destination=mock/storage.go -package=mock storage
type storage interface {
	GetUserByLogin(ctx context.Context, login string) (*user.User, error)
	NewUser(ctx context.Context, u user.User) (*user.User, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}

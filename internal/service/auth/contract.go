package auth

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/user"
)

//go:generate mockgen -source=contract.go -destination=mock/storage.go -package=mock storage
type storage interface {
	GetUser(ctx context.Context, login string) (*user.User, error)
	NewUser(ctx context.Context, u user.User) (*user.User, error)
}

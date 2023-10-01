package auth

import (
	"context"
	"time"

	"github.com/k1nky/gophermart/internal/entity/user"
)

type Storage interface {
	GetUser(ctx context.Context, login string) (*user.User, error)
	NewUser(ctx context.Context, u user.User) (*user.User, error)
}

type Service struct {
	secret          []byte
	tokenExpiration time.Duration
	store           Storage
}

func New(secret string, tokenExpiration time.Duration, store Storage) *Service {
	s := &Service{
		secret:          []byte(secret),
		tokenExpiration: tokenExpiration,
	}
	return s
}

func (s *Service) Register(ctx context.Context, newUser user.User) (string, error) {
	u, err := s.store.GetUser(ctx, newUser.Login)
	if err != nil {
		return "", err
	}
	if u != nil {
		// TODO: user already exist
	}
	if u, err = s.store.NewUser(ctx, newUser); err != nil {
		return "", err
	}
	return s.GenerateToken(TokenData{Login: u.Login})
}

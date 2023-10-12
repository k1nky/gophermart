package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/k1nky/gophermart/internal/entity/user"
)

type Service struct {
	secret          []byte
	tokenExpiration time.Duration
	store           storage
}

func New(secret string, tokenExpiration time.Duration, store storage) *Service {
	s := &Service{
		secret:          []byte(secret),
		tokenExpiration: tokenExpiration,
		store:           store,
	}
	return s
}

func (s *Service) Register(ctx context.Context, newUser user.User) (token string, err error) {
	var u *user.User
	if newUser.Password, err = user.HashPassword(newUser.Password); err != nil {
		return "", fmt.Errorf("auth: %w", err)
	}
	if u, err = s.store.NewUser(ctx, newUser); err != nil {
		return "", fmt.Errorf("auth: %w", err)
	}
	// TODO: handle token error
	token, err = s.GenerateToken(user.NewPrivateClaims(*u))
	return token, err
}

func (s *Service) Login(ctx context.Context, credentials user.User) (string, error) {
	u, err := s.store.GetUser(ctx, credentials.Login)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", fmt.Errorf("auth: %s %w", credentials.Login, user.ErrInvalidCredentials)
	}
	if err := u.CheckPassword(credentials.Password); err != nil {
		return "", fmt.Errorf("auth: %s %w", credentials.Login, user.ErrInvalidCredentials)
	}
	// TODO: handle token error
	token, err := s.GenerateToken(user.NewPrivateClaims(*u))
	return token, err
}

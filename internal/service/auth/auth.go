package auth

import (
	"context"
	"errors"
	"fmt"
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

var (
	ErrDuplicateError     = errors.New("login already exists")
	ErrInvalidCredentials = errors.New("login or password is not correct")
)

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
		return "", fmt.Errorf("%s %w", newUser.Login, ErrDuplicateError)
	}
	if u, err = s.store.NewUser(ctx, newUser); err != nil {
		return "", err
	}
	return s.GenerateToken(PrivateClaims{Login: u.Login})
}

func (s *Service) Login(ctx context.Context, credentials user.User) (string, error) {
	u, err := s.store.GetUser(ctx, credentials.Login)
	if err != nil {
		if u == nil {
			return "", fmt.Errorf("%s %w", credentials.Login, ErrInvalidCredentials)
		}
		return "", err
	}
	if err := u.CheckPassword(credentials.Password); err != nil {
		return "", fmt.Errorf("%s %w", credentials.Login, ErrInvalidCredentials)
	}
	return s.GenerateToken(PrivateClaims{Login: u.Login})
}

func (s *Service) IsDuplicateLogin(err error) bool {
	return errors.Is(err, ErrDuplicateError)
}

func (s *Service) IsIncorrectCredentials(err error) bool {
	return errors.Is(err, ErrInvalidCredentials)
}

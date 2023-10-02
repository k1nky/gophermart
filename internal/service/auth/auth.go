package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/k1nky/gophermart/internal/entity"
)

//go:generate mockgen -source=auth.go -destination=mock/storage.go -package=mock Storage
type Storage interface {
	GetUser(ctx context.Context, login string) (*entity.User, error)
	NewUser(ctx context.Context, u entity.User) (*entity.User, error)
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
		store:           store,
	}
	return s
}

func (s *Service) Register(ctx context.Context, newUser entity.User) (string, error) {
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
	token, err := s.GenerateToken(PrivateClaims{Login: u.Login})
	return token, err
}

func (s *Service) Login(ctx context.Context, credentials entity.User) (string, error) {
	u, err := s.store.GetUser(ctx, credentials.Login)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", fmt.Errorf("%s %w", credentials.Login, ErrInvalidCredentials)
	}
	if err := u.CheckPassword(credentials.Password); err != nil {
		return "", fmt.Errorf("%s %w", credentials.Login, ErrInvalidCredentials)
	}
	token, err := s.GenerateToken(PrivateClaims{Login: u.Login})
	return token, err
}

func (s *Service) IsDuplicateLogin(err error) bool {
	return errors.Is(err, ErrDuplicateError)
}

func (s *Service) IsIncorrectCredentials(err error) bool {
	return errors.Is(err, ErrInvalidCredentials)
}

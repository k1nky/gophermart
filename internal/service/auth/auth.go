package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/k1nky/gophermart/internal/entity/user"
)

type Service struct {
	secret          []byte
	tokenExpiration time.Duration
	store           storage
	log             logger
}

func New(secret string, tokenExpiration time.Duration, store storage, log logger) *Service {
	s := &Service{
		secret:          []byte(secret),
		tokenExpiration: tokenExpiration,
		store:           store,
		log:             log,
	}
	return s
}

// Регистрирует нового пользователя и возвращает jwt токен.
func (s *Service) Register(ctx context.Context, newUser user.User) (token string, err error) {
	var u *user.User

	fail := func(err error) (string, error) {
		wrapped := fmt.Errorf("auth: register failed: %w", err)
		s.log.Errorf("%s", wrapped.Error())
		return "", wrapped
	}

	if newUser.Password, err = user.HashPassword(newUser.Password); err != nil {
		return fail(err)
	}
	if u, err = s.store.NewUser(ctx, newUser); err != nil {
		return fail(err)
	}
	token, err = s.GenerateToken(user.NewPrivateClaims(*u))
	if err != nil {
		return fail(err)
	}
	return token, nil
}

// Аутентифицирует пользователя пользователя и возвращает jwt токен в случае успеха.
func (s *Service) Login(ctx context.Context, credentials user.User) (string, error) {
	fail := func(err error) (string, error) {
		wrapped := fmt.Errorf("auth: login failed for %s: %w", credentials.Login, err)
		if errors.Is(wrapped, user.ErrInvalidCredentials) {
			s.log.Debugf("%s", wrapped.Error())
		} else {
			s.log.Errorf("%s", wrapped.Error())
		}
		return "", wrapped
	}
	u, err := s.store.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return fail(err)
	}
	if u == nil {
		return fail(user.ErrInvalidCredentials)
	}
	if err := u.CheckPassword(credentials.Password); err != nil {
		return fail(user.ErrInvalidCredentials)
	}
	token, err := s.GenerateToken(user.NewPrivateClaims(*u))
	if err != nil {
		return fail(err)
	}
	return token, err
}

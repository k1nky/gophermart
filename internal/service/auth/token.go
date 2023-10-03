package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k1nky/gophermart/internal/entity/user"
)

type Claims struct {
	jwt.RegisteredClaims
	user.PrivateClaims
}

var (
	ErrInvalidToken = errors.New("invalid token")
)

func (s *Service) GenerateToken(claims user.PrivateClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiration)),
		},
		PrivateClaims: claims,
	})

	return token.SignedString(s.secret)
}

func (s *Service) ParseToken(signedToken string) (user.PrivateClaims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return user.PrivateClaims{}, err
	}
	if !token.Valid {
		return user.PrivateClaims{}, fmt.Errorf("token invalid")
	}
	return claims.PrivateClaims, nil
}

func (s *Service) IsInvalidToken(err error) bool {
	return errors.Is(err, ErrInvalidToken)
}

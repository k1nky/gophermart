package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k1nky/gophermart/internal/entity/user"
)

type Claims struct {
	jwt.RegisteredClaims
	user.PrivateClaims
}

func (s *Service) GenerateToken(claims user.PrivateClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiration)),
		},
		PrivateClaims: claims,
	})

	return token.SignedString(s.secret)
}

func (s *Service) parseToken(signedToken string) (user.PrivateClaims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return user.PrivateClaims{}, err
	}
	if !token.Valid {
		return user.PrivateClaims{}, err
	}
	return claims.PrivateClaims, nil
}

func (s *Service) Authorize(token string) (user.PrivateClaims, error) {
	claims, err := s.parseToken(token)
	if err != nil {
		return claims, fmt.Errorf("auth: invalid token: %w", user.ErrUnathorized)
	}
	return claims, nil
}

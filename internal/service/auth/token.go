package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PrivateClaims struct {
	Login string
}

type Claims struct {
	jwt.RegisteredClaims
	PrivateClaims
}

func (s *Service) GenerateToken(d PrivateClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiration)),
		},
		PrivateClaims: d,
	})

	return token.SignedString(s.secret)
}

func (s *Service) ValidateToken(signedToken string) (PrivateClaims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return PrivateClaims{}, err
	}
	if !token.Valid {
		return PrivateClaims{}, fmt.Errorf("token invalid")
	}
	return claims.PrivateClaims, nil
}

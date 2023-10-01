package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenData struct {
	Login string
}

type Claims struct {
	jwt.RegisteredClaims
	TokenData
}

func (s *Service) GenerateToken(d TokenData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiration)),
		},
		TokenData: d,
	})

	return token.SignedString(s.secret)
}

func (s *Service) ValidateToken(signedToken string) (TokenData, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return TokenData{}, err
	}
	if !token.Valid {
		return TokenData{}, fmt.Errorf("token invalid")
	}
	return claims.TokenData, nil
}

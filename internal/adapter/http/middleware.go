package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/k1nky/gophermart/internal/entity/user"
)

type contextKey int

const (
	keyUserClaims contextKey = iota
)

func AuthorizeMiddleware(auth authService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			claims, err := auth.Authorize(token)
			if err != nil {
				if errors.Is(err, user.ErrUnathorized) {
					http.Error(w, "", http.StatusUnauthorized)
				} else {
					http.Error(w, "", http.StatusInternalServerError)
				}
				return
			}
			ctx := context.WithValue(r.Context(), keyUserClaims, claims)
			newRequest := r.WithContext(ctx)
			next.ServeHTTP(w, newRequest)
		})
	}
}

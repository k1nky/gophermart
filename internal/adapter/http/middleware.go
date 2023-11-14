package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/k1nky/gophermart/internal/entity/user"
)

type contextKey int

const (
	keyUserClaims contextKey = iota
)

type loggingWriter struct {
	http.ResponseWriter
	code int
}

func (bw *loggingWriter) WriteHeader(statusCode int) {
	bw.code = statusCode
	bw.ResponseWriter.WriteHeader(statusCode)
}

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

func LoggingMiddleware(l logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			bw := &loggingWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(bw, r)
			l.Infof("%s %s status %d duration %s", r.Method, r.RequestURI, bw.code, time.Since(start))
		})
	}
}

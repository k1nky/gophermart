package http

import (
	"net/http"
)

func AuthorizeMiddleware(auth AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			_, err := auth.ParseToken(token)
			if err != nil {
				if auth.IsInvalidToken(err) {
					http.Error(w, "", http.StatusUnauthorized)
				} else {
					http.Error(w, "", http.StatusInternalServerError)
				}
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

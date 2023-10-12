package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Adapter struct {
	auth    authService
	account accountService
}

func New(ctx context.Context, address string, port int, auth authService, account accountService) *Adapter {
	a := &Adapter{
		auth:    auth,
		account: account,
	}

	srv := &http.Server{
		Handler:      a.buildRouter(),
		Addr:         fmt.Sprintf("%s:%d", address, port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		// TODO: graceful shutdown
		if err := srv.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()

	return a
}

func (a *Adapter) buildRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.Register)
		r.Post("/login", a.Login)
	})

	r.With(AuthorizeMiddleware(a.auth)).Route("/api/user/orders", func(r chi.Router) {
		r.Get("/", a.GetOrder)
		r.Post("/", a.NewOrder)
	})
	return r
}

func (a *Adapter) writeJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

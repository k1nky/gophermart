package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Adapter struct {
	auth    authService
	account accountService
}

func New(ctx context.Context, listen string, auth authService, account accountService) *Adapter {
	a := &Adapter{
		auth:    auth,
		account: account,
	}

	srv := &http.Server{
		Handler:      a.buildRouter(),
		Addr:         listen,
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
		r.With(AuthorizeMiddleware(a.auth)).Get("/balance", a.GetBalance)
		r.With(AuthorizeMiddleware(a.auth)).Get("/orders", a.GetOrder)
		r.With(AuthorizeMiddleware(a.auth)).Post("/orders", a.NewOrder)
		// r.With(AuthorizeMiddleware(a.auth)).Get("/balance/withdrawals", a.GetWithdrawals)
		r.With(AuthorizeMiddleware(a.auth)).Get("/withdrawals", a.GetWithdrawals)
		r.With(AuthorizeMiddleware(a.auth)).Post("/balance/withdraw", a.NewWithdraw)
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

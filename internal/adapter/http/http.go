package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	DefaultReadTimeout  = 10 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	DefaultCloseTimeout = 5 * time.Second
)

type Adapter struct {
	auth    authService
	account accountService
	log     logger
}

func New(auth authService, account accountService, log logger) *Adapter {
	a := &Adapter{
		auth:    auth,
		account: account,
		log:     log,
	}

	return a
}

func (a *Adapter) ListenAndServe(ctx context.Context, addr string) {
	srv := &http.Server{
		Handler:      a.buildRouter(),
		Addr:         addr,
		WriteTimeout: DefaultReadTimeout,
		ReadTimeout:  DefaultWriteTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.log.Debugf("http server was closed")
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.Errorf("unexpected server closing: %v", err)
			}
		}
	}()
	go func() {
		<-ctx.Done()
		a.log.Debugf("closing http server")
		c, cancel := context.WithTimeout(context.Background(), DefaultCloseTimeout)
		defer cancel()
		srv.Shutdown(c)
	}()
}

func (a *Adapter) buildRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(LoggingMiddleware(a.log))
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.Register)
		r.Post("/login", a.Login)
		r.With(AuthorizeMiddleware(a.auth)).Get("/balance", a.GetBalance)
		r.With(AuthorizeMiddleware(a.auth)).Get("/orders", a.GetOrder)
		r.With(AuthorizeMiddleware(a.auth)).Post("/orders", a.NewOrder)
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

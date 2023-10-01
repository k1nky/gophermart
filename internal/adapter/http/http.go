package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/k1nky/gophermart/internal/entity/user"
)

type AuthService interface {
	Register(ctx context.Context, credentials user.User) error
	Login(ctx context.Context, credentials user.User) error
}

type Adapter struct {
	auth AuthService
}

func New(ctx context.Context, address string, port int) *Adapter {
	a := &Adapter{}

	r := chi.NewRouter()
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.Register)
		r.Post("/login", a.Login)
	})

	srv := &http.Server{
		Handler:      r,
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

// Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
// После успешной регистрации должна происходить автоматическая аутентификация пользователя.
//
// POST /api/user/register HTTP/1.1
// Content-Type: application/json
// ...
//
//	{
//		"login": "<login>",
//		"password": "<password>"
//	}
//
// Возможные коды ответа:
// - `200` — пользователь успешно зарегистрирован и аутентифицирован;
// - `400` — неверный формат запроса;
// - `409` — логин уже занят;
// - `500` — внутренняя ошибка сервера.
func (a *Adapter) Register(w http.ResponseWriter, r *http.Request) {
	credentials := user.User{}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.auth.Register(r.Context(), credentials)
	if err != nil {
		// TODO: если логин уже заняты
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Аутентификация производится по паре логин/пароль.
// Формат запроса:
//
//	{
//		"login": "<login>",
//		"password": "<password>"
//	}
//
// ```
// Возможные коды ответа:
// - `200` — пользователь успешно аутентифицирован;
// - `400` — неверный формат запроса;
// - `401` — неверная пара логин/пароль;
// - `500` — внутренняя ошибка сервера.
func (a *Adapter) Login(w http.ResponseWriter, r *http.Request) {
	credentials := user.User{}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := a.auth.Login(r.Context(), credentials)
	if err != nil {
		// TODO: неверная пара логин/пароль
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

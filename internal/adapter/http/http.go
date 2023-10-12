package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
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
		r.Get("/", nil)
		r.Post("/", a.NewOrder)
	})
	return r
}

// Регистрация пользователя. Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
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
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := credentials.IsValid(); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	signedToken, err := a.auth.Register(r.Context(), credentials)
	if err != nil {
		if errors.Is(err, user.ErrDuplicateLogin) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Authorization", signedToken)
	w.WriteHeader(http.StatusOK)
}

// Аутентификация пользователя. Аутентификация производится по паре логин/пароль.
// Формат запроса:
//
//	 ```
//		{
//			"login": "<login>",
//			"password": "<password>"
//		}
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
	if err := credentials.IsValid(); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	signedToken, err := a.auth.Login(r.Context(), credentials)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Authorization", signedToken)
	w.WriteHeader(http.StatusOK)
}

// Загрузка номера заказа. Хендлер доступен только аутентифицированным пользователям. Номером заказа является последовательность цифр произвольной длины.
// Номер заказа может быть проверен на корректность ввода с помощью [алгоритма Луна](https://ru.wikipedia.org/wiki/Алгоритм_Луна){target="_blank"}.
// Формат запроса:
// ```
// POST /api/user/orders HTTP/1.1
// Content-Type: text/plain
// ...
// 12345678903
// ```
// Возможные коды ответа:
// - `200` — номер заказа уже был загружен этим пользователем;
// - `202` — новый номер заказа принят в обработку;
// - `400` — неверный формат запроса;
// - `401` — пользователь не аутентифицирован;
// - `409` — номер заказа уже был загружен другим пользователем;
// - `422` — неверный формат номера заказа;
// - `500` — внутренняя ошибка сервера.
func (a *Adapter) NewOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	newOrder := order.Order{
		UserID: claims.ID,
	}
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	newOrder.Number = order.OrderNumber(buf.String())
	if !newOrder.Number.IsValid() {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}
	if _, err := a.account.NewOrder(r.Context(), newOrder); err != nil {
		if errors.Is(err, order.ErrDuplicateOrder) {
			http.Error(w, "", http.StatusOK)
		} else if errors.Is(err, order.ErrOrderBelongsToAnotherUser) {
			http.Error(w, "", http.StatusConflict)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// Получение списка загруженных номеров заказов. Хендлер доступен только авторизованному пользователю. Номера заказа в выдаче должны быть отсортированы по времени загрузки от самых старых к самым новым. Формат даты — RFC3339.
// Доступные статусы обработки расчётов:
// - `NEW` — заказ загружен в систему, но не попал в обработку;
// - `PROCESSING` — вознаграждение за заказ рассчитывается;
// - `INVALID` — система расчёта вознаграждений отказала в расчёте;
// - `PROCESSED` — данные по заказу проверены и информация о расчёте успешно получена.
// Формат запроса:
// ```
// GET /api/user/orders HTTP/1.1
// Content-Length: 0
// ```
// Возможные коды ответа:
//   - `200` — успешная обработка запроса.
//     Формат ответа:
//     ```
//     200 OK HTTP/1.1
//     Content-Type: application/json
//     ...
//     [
//     {
//     "number": "9278923470",
//     "status": "PROCESSED",
//     "accrual": 500,
//     "uploaded_at": "2020-12-10T15:15:45+03:00"
//     },
//     {
//     "number": "12345678903",
//     "status": "PROCESSING",
//     "uploaded_at": "2020-12-10T15:12:01+03:00"
//     },
//     {
//     "number": "346436439",
//     "status": "INVALID",
//     "uploaded_at": "2020-12-09T16:09:53+03:00"
//     }
//     ]
//     ```
//   - `204` — нет данных для ответа.
//   - `401` — пользователь не авторизован.
//   - `500` — внутренняя ошибка сервера.
func (a *Adapter) GetOrder(w http.ResponseWriter, r *http.Request) {}

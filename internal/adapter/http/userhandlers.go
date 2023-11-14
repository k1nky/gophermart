package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/k1nky/gophermart/internal/entity/user"
)

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

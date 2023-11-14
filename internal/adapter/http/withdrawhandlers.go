package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/k1nky/gophermart/internal/entity/withdraw"
)

// Получение текущего баланса пользователя
// Хендлер доступен только авторизованному пользователю. В ответе должны содержаться данные о текущей сумме баллов лояльности, а также сумме использованных за весь период регистрации баллов.
// Формат запроса:
// ```
// GET /api/user/balance HTTP/1.1
// Content-Length: 0
// ```
// Возможные коды ответа:
//   - `200` — успешная обработка запроса.
//     Формат ответа:
//     ```
//     200 OK HTTP/1.1
//     Content-Type: application/json
//     ...
//     {
//     "current": 500.5,
//     "withdrawn": 42
//     }
//     ```
//   - `401` — пользователь не авторизован.
//   - `500` — внутренняя ошибка сервера.
func (a *Adapter) GetBalance(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	balance, err := a.account.GetUserBalance(r.Context(), claims.ID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err := a.writeJSON(w, balance); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// Получение информации о выводе средств.
// Хендлер доступен только авторизованному пользователю. Факты выводов в выдаче должны быть отсортированы по времени вывода от самых старых к самым новым. Формат даты — RFC3339.
// Формат запроса:
// ```
// GET /api/user/withdrawals HTTP/1.1
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
//     "order": "2377225624",
//     "sum": 500,
//     "processed_at": "2020-12-09T16:09:57+03:00"
//     }
//     ]
//     ```
//   - `204` - нет ни одного списания.
//   - `401` — пользователь не авторизован.
//   - `500` — внутренняя ошибка сервера.
func (a *Adapter) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	withdrawals, err := a.account.GetUserWithdrawals(r.Context(), claims.ID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := a.writeJSON(w, withdrawals); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// Запрос на списание средств
// Хендлер доступен только авторизованному пользователю. Номер заказа представляет собой гипотетический номер нового заказа пользователя в счет оплаты которого списываются баллы.
// Примечание: для успешного списания достаточно успешной регистрации запроса, никаких внешних систем начисления не предусмотрено и не требуется реализовывать.
// Формат запроса:
// ```
// POST /api/user/balance/withdraw HTTP/1.1
// Content-Type: application/json

//	{
//		"order": "2377225624",
//	    "sum": 751
//	}
//
// ```
// Здесь `order` — номер заказа, а `sum` — сумма баллов к списанию в счёт оплаты.
// Возможные коды ответа:
// - `200` — успешная обработка запроса;
// - `401` — пользователь не авторизован;
// - `402` — на счету недостаточно средств;
// - `422` — неверный номер заказа;
// - `500` — внутренняя ошибка сервера.
func (a *Adapter) NewWithdraw(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	newWithdraw := withdraw.Withdraw{}
	if err := json.NewDecoder(r.Body).Decode(&newWithdraw); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	newWithdraw.UserID = claims.ID
	if err := a.account.NewWithdraw(r.Context(), newWithdraw); err != nil {
		if errors.Is(err, withdraw.ErrInsufficientBalance) {
			http.Error(w, "", http.StatusPaymentRequired)
		} else if errors.Is(err, order.ErrInvalidNumberFormat) {
			http.Error(w, "", http.StatusUnprocessableEntity)
		} else {
			http.Error(w, "", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

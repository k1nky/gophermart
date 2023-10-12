package http

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

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
func (a *Adapter) GetOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(keyUserClaims).(user.PrivateClaims)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	orders, err := a.account.GetUserOrders(r.Context(), claims.ID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := a.writeJSON(w, orders); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

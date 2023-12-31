package database

import (
	"context"
	"fmt"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

func (a *Adapter) selectOrders(ctx context.Context, where string, limit uint, args ...interface{}) ([]*order.Order, error) {
	query := fmt.Sprintf(`SELECT order_id, number, status, accrual, uploaded_at, user_id FROM orders WHERE %s`, where)
	if limit > 0 {
		query = query + fmt.Sprintf(" LIMIT %d", limit)
	}
	orders := make([]*order.Order, 0)
	rows, err := a.QueryContext(ctx, query, args...)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	for rows.Next() {
		o := &order.Order{}
		if err := rows.Scan(&o.ID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.UserID); err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil
}

// Возвращает заказ по номеру
func (a *Adapter) GetOrderByNumber(ctx context.Context, number order.OrderNumber) (*order.Order, error) {
	orders, err := a.selectOrders(ctx, "number = $1", 1, number)
	if err != nil {
		return nil, NewExecutingQueryError(err)
	}
	if len(orders) == 0 {
		return nil, nil
	}
	return orders[0], err
}

// Возвращает не более maxRows заказов с заданными статусами
func (a *Adapter) GetOrdersByStatus(ctx context.Context, statuses []order.OrderStatus, maxRows uint) ([]*order.Order, error) {
	args := make([]string, 0, len(statuses))
	// преобразуем в совместимый с postgres тип
	for _, v := range statuses {
		args = append(args, string(v))
	}
	orders, err := a.selectOrders(ctx, "status = any($1::order_status[])", maxRows, args)
	if err != nil {
		err = NewExecutingQueryError(err)
	}
	return orders, err
}

// Возвращает не более maxRows заказов для указанного пользователя в порядке возрастания даты загрузки
func (a *Adapter) GetOrdersByUserID(ctx context.Context, userID user.ID, maxRows uint) ([]*order.Order, error) {
	orders, err := a.selectOrders(ctx, "user_id = $1 ORDER BY uploaded_at ASC", maxRows, userID)
	if err != nil {
		err = NewExecutingQueryError(err)
	}
	return orders, err
}

// Создает новый заказ и возвращает его
func (a *Adapter) NewOrder(ctx context.Context, o order.Order) (*order.Order, error) {
	const query = `
		INSERT INTO orders AS o (user_id, number, status)
		VALUES ($1, $2, 'NEW')
		RETURNING o.order_id, o.uploaded_at
	`
	row := a.QueryRowContext(ctx, query, o.UserID, o.Number)
	if err := row.Err(); err != nil {
		if a.hasUniqueViolationError(err) {
			return nil, fmt.Errorf("%s %w", o.Number, order.ErrDuplicated)
		}
		return nil, NewExecutingQueryError(err)
	}
	if err := row.Scan(&o.ID, &o.UploadedAt); err != nil {
		return nil, NewExecutingQueryError(err)
	}
	o.Status = order.StatusNew
	return &o, nil

}

// Обновляет заказ
func (a *Adapter) UpdateOrder(ctx context.Context, o order.Order) error {
	// не допускаем обновление уже обработанного заказ
	const updateOrderQuery = `
		UPDATE orders 
		SET status = $1, accrual = $2
		WHERE order_id = $3 AND status <> 'PROCESSED'
	`
	tx, err := a.BeginTx(ctx, nil)
	if err != nil {
		return NewExecutingQueryError(err)
	}
	defer tx.Rollback()

	if r, err := tx.ExecContext(ctx, updateOrderQuery, o.Status, o.Accrual, o.ID); err != nil {
		return NewExecutingQueryError(err)
	} else {
		if rows, err := r.RowsAffected(); rows == 0 || err != nil {
			return fmt.Errorf("%s %w", o.Number, order.ErrAlreadyProcessed)
		}
	}
	// добавляем соответствующую транзакцию
	if o.Accrual != nil && o.Status == order.StatusProcessed {
		// получаем последнюю транзакцию пользователя
		// для новой транзакции увеличиваем последовательный номер транзакции пользователя на 1 и баланс на размер начисления
		// последовательный номер уникальный для каждого пользователя
		const transactionQuery = `
			WITH user_balance AS (
				SELECT coalesce(sum(user_transaction_seq), 0) seq, coalesce(sum(balance), 0) balance FROM transactions WHERE user_id = $1 ORDER BY seq DESC LIMIT 1
			)
			INSERT INTO transactions(
				user_id,
				user_transaction_seq,
				source_id, source_type,
				balance
			) VALUES (
				$1,
				(SELECT seq FROM user_balance) + 1,
				$2, 'ACCRUAL',
				(SELECT balance FROM user_balance) + $3
			)
		`
		if _, err := tx.ExecContext(ctx, transactionQuery, o.UserID, o.ID, o.Accrual); err != nil {
			return NewExecutingQueryError(err)
		}
	}
	if err = tx.Commit(); err != nil {
		return NewExecutingQueryError(err)
	}
	return nil
}

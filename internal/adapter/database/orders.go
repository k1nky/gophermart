package database

import (
	"context"
	"fmt"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

func (a *Adapter) selectOrders(ctx context.Context, where string, limit uint, args ...interface{}) ([]*order.Order, error) {
	// TODO: join on transactions
	query := fmt.Sprintf(`SELECT order_id, number, status, accrual, uploaded_at, user_id, updated_at FROM orders WHERE %s`, where)
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
		if err := rows.Scan(&o.ID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.UserID, &o.UpdatedAt); err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil
}

func (a *Adapter) GetOrderByNumber(ctx context.Context, number order.OrderNumber) (*order.Order, error) {
	orders, err := a.selectOrders(ctx, "number = $1", 1, number)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, nil
	}
	return orders[0], err
}

func (a *Adapter) GetOrdersByStatus(ctx context.Context, statuses []order.OrderStatus) ([]*order.Order, error) {
	// TODO: pass limit as an argument
	args := make([]string, 0, len(statuses))
	for _, v := range statuses {
		args = append(args, string(v))
	}
	orders, err := a.selectOrders(ctx, "status = any($1::order_status[])", 100, args)
	return orders, err
}

func (a *Adapter) GetOrdersByUserID(ctx context.Context, userID user.ID) ([]*order.Order, error) {
	// TODO: pass limit as an argument
	orders, err := a.selectOrders(ctx, "user_id = $1 ORDER BY uploaded_at ASC", 100, userID)
	return orders, err
}

func (a *Adapter) NewOrder(ctx context.Context, o order.Order) (*order.Order, error) {
	const query = `
		INSERT INTO orders AS o (user_id, number, status)
		VALUES ($1, $2, 'NEW')
		RETURNING o.order_id, o.uploaded_at
	`
	row := a.QueryRowContext(ctx, query, o.UserID, o.Number)
	if err := row.Err(); err != nil {
		if a.hasUniqueViolationError(err) {
			return nil, fmt.Errorf("%s %w", o.Number, order.ErrDuplicateOrder)
		}
		return nil, err
	}
	if err := row.Scan(&o.ID, &o.UploadedAt); err != nil {
		return nil, err
	}
	o.Status = order.StatusNew
	return &o, nil

}

func (a *Adapter) UpdateOrder(ctx context.Context, o order.Order) error {
	const updateOrderQuery = `
		UPDATE orders 
		SET status = $1
		WHERE order_id = $2
	`
	tx, err := a.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, updateOrderQuery, o.Status, o.ID); err != nil {
		return err
	}
	if o.Accrual != nil {
		const transactionQuery = `
			INSERT INTO transactions(source_id, balance, amount, normal)
			VALUES ($1, (SELECT balance FROM transactions WHERE transaction_id=(SELECT MAX(transaction_id) FROM transactions WHERE user_id=$3)) + $2, $2, 1)
		`
	}
	err = tx.Commit()
	return err
}

// func (a *Adapter) UpdateBalance(ctx context.Context, userID user.ID) error {
// 	const query = `
// 	insert into balance (user_id, value, updated_at) values(1, (select (select sum(accrual) as s from orders where user_id=1 and updated_at > b.updated_at) - (select sum(sum) s from withdrawals where user_id=1 and processed_at is null) from balance b where user_id=1), NOW()) on conflict on constraint balance_user_i
// 	d_key do update set value = excluded.value;
// 	`
// }

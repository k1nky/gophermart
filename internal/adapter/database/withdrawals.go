package database

import (
	"context"
	"fmt"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/k1nky/gophermart/internal/entity/withdraw"
)

func (a *Adapter) selectWithdrawals(ctx context.Context, where string, limit uint, args ...interface{}) ([]*withdraw.Withdraw, error) {
	query := fmt.Sprintf(`SELECT withdraw_id, user_id, amount, order_number, processed_at FROM withdrawals WHERE %s`, where)
	if limit > 0 {
		query = query + fmt.Sprintf(" LIMIT %d", limit)
	}
	withdrawals := make([]*withdraw.Withdraw, 0)
	rows, err := a.QueryContext(ctx, query, args...)
	if err != nil {
		return withdrawals, err
	}
	defer rows.Close()
	for rows.Next() {
		w := &withdraw.Withdraw{}
		if err := rows.Scan(&w.ID, &w.UserID, &w.Sum, &w.Number, &w.ProcessedAt); err != nil {
			return withdrawals, err
		}
		withdrawals = append(withdrawals, w)
	}
	if err := rows.Err(); err != nil {
		return withdrawals, err
	}

	return withdrawals, nil
}

// Возвращает список запросов на списание для указанного пользователя
func (a *Adapter) GetWithdrawalsByUserID(ctx context.Context, userID user.ID) ([]*withdraw.Withdraw, error) {
	// TODO: pass limit as an argument
	withdrawals, err := a.selectWithdrawals(ctx, "user_id = $1 ORDER BY processed_at ASC", 100, userID)
	return withdrawals, err
}

// Возврашает баланс указанного пользователя
func (a *Adapter) GetBalanceByUser(ctx context.Context, userID user.ID) (user.Balance, error) {
	balance := user.Balance{}
	// получаем баланс из последней транзакции пользователя
	const query = `
		SELECT COALESCE(SUM(balance), 0) 
		FROM transactions 
		WHERE transaction_id=(SELECT MAX(transaction_id) FROM transactions WHERE user_id=$1)
	`
	row := a.QueryRowContext(ctx, query, userID)
	if err := row.Err(); err != nil {
		return balance, err
	}
	if err := row.Scan(&balance.Current); err != nil {
		return balance, err
	}
	// получаем сумму всех списаний
	row = a.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount),0) FROM withdrawals WHERE user_id=$1 AND processed_at IS NOT NULL`, userID)
	if err := row.Err(); err != nil {
		return balance, err
	}
	if err := row.Scan(&balance.Withdrawn); err != nil {
		return balance, err
	}
	return balance, nil
}

// Создает новое списание и возвращает его
func (a *Adapter) NewWithdraw(ctx context.Context, w withdraw.Withdraw) (*withdraw.Withdraw, error) {
	tx, err := a.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// создаем новое списание
	const newWithdrawQuery = `
		INSERT INTO withdrawals (user_id, amount, order_number, processed_at) 
		VALUES($1, $2, $3, NOW())
		RETURNING withdraw_id, processed_at
	`
	row := a.QueryRowContext(ctx, newWithdrawQuery, w.UserID, w.Sum, w.Number)
	if err := row.Err(); err != nil {
		if a.hasUniqueViolationError(err) {
			return nil, fmt.Errorf("order %s %w", w.Number, order.ErrDuplicated)
		}
		return nil, err
	}
	if err := row.Scan(&w.ID, &w.ProcessedAt); err != nil {
		return nil, err
	}

	// добавляем новую транзакцию на списание
	const newTransactionQuery = `
		WITH user_balance AS (
			SELECT COALESCE(SUM(user_transaction_seq), 0) seq, COALESCE(SUM(balance), 0) balance FROM transactions WHERE user_id = $1 ORDER BY seq DESC LIMIT 1
		)
		INSERT INTO transactions(
			user_id,
			user_transaction_seq,
			source_id, source_type,
			balance
		) VALUES (
			$1,
			(SELECT seq FROM user_balance) + 1,
			$2, 'WITHDRAW',
			(SELECT balance FROM user_balance) - $3
		)
		RETURNING balance
	`
	row = tx.QueryRowContext(ctx, newTransactionQuery, w.UserID, w.ID, w.Sum)
	if err := row.Err(); err != nil {
		return nil, err
	}
	// если в результате списание баланс отрицательный, то откатываем транзакцию
	var balance float32
	if err := row.Scan(&balance); err != nil {
		return nil, err
	}
	if balance < 0 {
		return nil, withdraw.ErrInsufficientBalance
	}
	err = tx.Commit()
	return &w, err
}

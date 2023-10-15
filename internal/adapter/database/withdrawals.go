package database

import (
	"context"
	"fmt"

	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/k1nky/gophermart/internal/entity/withdraw"
)

func (a *Adapter) selectWithdrawals(ctx context.Context, where string, limit uint, args ...interface{}) ([]*withdraw.Withdraw, error) {
	query := fmt.Sprintf(`
		SELECT withdraw_id, user_id, amount, order_number, processed_at FROM withdrawals WHERE %s
	`, where)
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

func (a *Adapter) GetWithdrawalsByUserID(ctx context.Context, userID user.ID) ([]*withdraw.Withdraw, error) {
	// TODO: pass limit as an argument
	withdrawals, err := a.selectWithdrawals(ctx, "user_id = $1 ORDER BY processed_at ASC", 100, userID)
	return withdrawals, err
}

func (a *Adapter) GetBalanceByUser(ctx context.Context, userID user.ID) (user.Balance, error) {
	balance := user.Balance{}
	const query = `
		SELECT coalesce(sum(balance), 0) 
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
	row = a.QueryRowContext(ctx, `select coalesce(sum(amount),0) from withdrawals where user_id=$1 and processed_at is not null`, userID)
	if err := row.Err(); err != nil {
		return balance, err
	}
	if err := row.Scan(&balance.Withdrawn); err != nil {
		return balance, err
	}
	return balance, nil
}

func (a *Adapter) NewWithdraw(ctx context.Context, w withdraw.Withdraw) (*withdraw.Withdraw, error) {
	const newWithdrawQuery = `
		INSERT INTO withdrawals (user_id, amount, order_number, processed_at) 
		VALUES($1, $2, $3, NOW())
		RETURNING withdraw_id, processed_at
	`
	tx, err := a.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	row := a.QueryRowContext(ctx, newWithdrawQuery, w.UserID, w.Sum, w.Number)
	if err := row.Err(); err != nil {
		if a.hasUniqueViolationError(err) {
			// TODO: wrap error
			return nil, fmt.Errorf("duplicated withdraw")
			// return nil, fmt.Errorf("%s %w", u.Login, user.ErrDuplicateLogin)
		}
		return nil, err
	}
	if err := row.Scan(&w.ID, &w.ProcessedAt); err != nil {
		return nil, err
	}

	const newTransactionQuery = `
		INSERT INTO transactions(user_id, source_id, balance, source_type)
		VALUES (
			$1, $2, (
				SELECT coalesce(sum(balance), 0) 
				FROM transactions 
				WHERE transaction_id=(
					SELECT MAX(transaction_id) FROM transactions WHERE user_id=$1)
				) - $3, 'WITHDRAW'
		)
	`
	if _, err := tx.ExecContext(ctx, newTransactionQuery, w.UserID, w.ID, w.Sum); err != nil {
		return nil, err
	}
	err = tx.Commit()
	return &w, err
}

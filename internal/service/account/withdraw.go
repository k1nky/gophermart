package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/k1nky/gophermart/internal/entity/withdraw"
)

// Возвращает баланс пользователя
func (s *Service) GetUserBalance(ctx context.Context, userID user.ID) (user.Balance, error) {
	b, err := s.store.GetBalanceByUser(ctx, userID)
	if err != nil {
		err = fmt.Errorf("account: get user balance: %w", err)
		s.log.Errorf("%s", err)
	}
	return b, err
}

// Возвращает списания пользователя
func (s *Service) GetUserWithdrawals(ctx context.Context, userID user.ID) ([]*withdraw.Withdraw, error) {
	withdrawals, err := s.store.GetWithdrawalsByUserID(ctx, userID, DefaultMaxRows)
	if err != nil {
		err = fmt.Errorf("account: get user withdrawals: %w", err)
		s.log.Errorf("%s", err)
	}
	return withdrawals, err
}

// Проводит новое списание
func (s *Service) NewWithdraw(ctx context.Context, w withdraw.Withdraw) error {
	_, err := s.store.NewWithdraw(ctx, w)
	if err != nil && !errors.Is(err, withdraw.ErrInsufficientBalance) {
		err = fmt.Errorf("account: new withdrawals: %w", err)
		s.log.Errorf("%s", err)
	}
	return err
}

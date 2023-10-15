package account

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/k1nky/gophermart/internal/entity/withdraw"
)

func (s *Service) GetUserBalance(ctx context.Context, userID user.ID) (user.Balance, error) {
	b, err := s.store.GetBalanceByUser(ctx, userID)
	return b, err
}

func (s *Service) GetUserWithdrawals(ctx context.Context, userID user.ID) ([]*withdraw.Withdraw, error) {
	withdrawals, err := s.store.GetWithdrawalsByUserID(ctx, userID)
	return withdrawals, err
}

func (s *Service) NewWithdraw(ctx context.Context, w withdraw.Withdraw) error {
	_, err := s.store.NewWithdraw(ctx, w)
	return err
	// o, err := s.store.GetOrderByNumber(ctx, newOrder.Number)
	// if err != nil {
	// 	return nil, err
	// }
	// if o != nil {
	// 	if o.UserID != newOrder.UserID {
	// 		return nil, order.ErrOrderBelongsToAnotherUser
	// 	}
	// 	return nil, order.ErrDuplicateOrder
	// }
	// o, err = s.store.NewOrder(ctx, newOrder)
	// return o, err
}

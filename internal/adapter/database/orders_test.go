package database

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/stretchr/testify/suite"
)

type ordersTestSuite struct {
	suite.Suite
	a *Adapter
}

func (suite *ordersTestSuite) SetupTest() {
	if shouldSkipDBTest(suite.T()) {
		return
	}
	var err error
	if suite.a, err = openTestDB(); err != nil {
		suite.FailNow(err.Error())
		return
	}
	if _, err := suite.a.Exec(`
		DELETE FROM transactions CASCADE;
		DELETE FROM orders CASCADE;
		DELETE FROM users CASCADE;

		INSERT INTO users(user_id, login, password) 
			VALUES (1, 'u1', 'p1'), 
					(2, 'u2', 'p2');
		INSERT INTO orders(order_id, user_id, number, status)
			VALUES (1, 1, '100', 'NEW'), (2, 1, '200', 'NEW'), (3, 1, '300', 'INVALID'), (4, 1, '400', 'PROCESSED');				
	`); err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *ordersTestSuite) TestNewOrder() {
	o := order.Order{
		Number: "999",
		UserID: user.ID(2),
	}
	newOrder, err := suite.a.NewOrder(context.TODO(), o)
	suite.NoError(err)
	suite.Equal(o.Number, newOrder.Number)
	suite.Equal(order.StatusNew, newOrder.Status)
	suite.NotEqual(0, newOrder.ID)
}

func (suite *ordersTestSuite) TestNewOrderDuplicate() {
	o := order.Order{
		Number: "100",
		UserID: user.ID(2),
	}
	got, err := suite.a.NewOrder(context.TODO(), o)
	suite.ErrorIs(err, order.ErrDuplicated)
	suite.Nil(got)
}

func (suite *ordersTestSuite) TestUpdateOrder() {
	var v float32 = 120.0

	o := order.Order{
		ID:      1,
		Number:  "100",
		Status:  order.StatusProcessed,
		Accrual: &v,
		UserID:  user.ID(1),
	}
	err := suite.a.UpdateOrder(context.TODO(), o)
	suite.NoError(err)
}

func (suite *ordersTestSuite) TestUpdateOrderProcessed() {
	var v float32 = 120.0

	o := order.Order{
		ID:      4,
		Number:  "400",
		Status:  order.StatusInvalid,
		Accrual: &v,
		UserID:  user.ID(1),
	}
	err := suite.a.UpdateOrder(context.TODO(), o)
	suite.ErrorIs(err, order.ErrAlreadyProcessed)
}

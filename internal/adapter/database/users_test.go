package database

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/stretchr/testify/suite"
)

type usersTestSuite struct {
	suite.Suite
	a *Adapter
}

func (suite *usersTestSuite) SetupTest() {
	if shouldSkipDBTest(suite.T()) {
		return
	}
	var err error
	if suite.a, err = openTestDB(); err != nil {
		suite.FailNow(err.Error())
		return
	}
	if _, err := suite.a.Exec(`
		DELETE FROM users CASCADE;
		INSERT INTO users(user_id, login, password) 
			VALUES (1, 'u1', 'p1'), 
					(2, 'u2', 'p2');
	`); err != nil {
		suite.FailNow(err.Error())
	}

}

func (suite *usersTestSuite) TestNewUser() {
	u := user.User{
		Login:    "test_u",
		Password: "test_p",
	}
	newUser, err := suite.a.NewUser(context.TODO(), u)
	suite.NoError(err)
	suite.Equal(u.Login, newUser.Login)
	suite.Equal(u.Password, newUser.Password)
	suite.NotEqual(0, newUser.ID)
}

func (suite *usersTestSuite) TestNewUserDuplicate() {
	u := user.User{
		ID:       1,
		Login:    "u1",
		Password: "p1",
	}
	got, err := suite.a.NewUser(context.TODO(), u)
	suite.ErrorIs(err, user.ErrDuplicateLogin, "")
	suite.Nil(got, "")
}

func (suite *usersTestSuite) TestGetUserByLogin() {
	u := &user.User{
		ID:       1,
		Login:    "u1",
		Password: "p1",
	}
	got, err := suite.a.GetUserByLogin(context.TODO(), "u1")
	suite.NoError(err)
	suite.Equal(u, got)
}

func (suite *usersTestSuite) TestGetUserByLoginNotExists() {
	got, err := suite.a.GetUserByLogin(context.TODO(), "u1000")
	suite.NoError(err)
	suite.Nil(got)
}

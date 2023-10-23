package database

import (
	"context"

	"github.com/k1nky/gophermart/internal/entity/user"
)

func (suite *adapterTestSuite) TestNewUser() {
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

func (suite *adapterTestSuite) TestNewUserDuplicate() {
	u := user.User{
		ID:       1,
		Login:    "u1",
		Password: "p1",
	}
	got, err := suite.a.NewUser(context.TODO(), u)
	suite.ErrorIs(err, user.ErrDuplicateLogin, "")
	suite.Nil(got, "")
}

func (suite *adapterTestSuite) TestGetUserByLogin() {
	u := &user.User{
		ID:       1,
		Login:    "u1",
		Password: "p1",
	}
	got, err := suite.a.GetUserByLogin(context.TODO(), "u1")
	suite.NoError(err)
	suite.Equal(u, got)
}

func (suite *adapterTestSuite) TestGetUserByLoginNotExists() {
	got, err := suite.a.GetUserByLogin(context.TODO(), "u2")
	suite.NoError(err)
	suite.Nil(got)
}

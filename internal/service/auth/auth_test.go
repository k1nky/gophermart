package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophermart/internal/entity/user"
	log "github.com/k1nky/gophermart/internal/logger"
	"github.com/k1nky/gophermart/internal/service/auth/mock"
	"github.com/stretchr/testify/suite"
)

type authServiceTestSuite struct {
	suite.Suite
	store *mock.Mockstorage
	svc   *Service
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(authServiceTestSuite))
}

func (suite *authServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.Suite.T())
	suite.store = mock.NewMockstorage(ctrl)
	suite.svc = New("secret", 3*time.Hour, suite.store, &log.Blackhole{})
}

func (suite *authServiceTestSuite) TestRegisterNewUser() {
	u := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().NewUser(gomock.Any(), gomock.Any()).Return(&user.User{
		Login:    "user",
		Password: "password",
		ID:       1,
	}, nil)

	token, err := suite.svc.Register(ctx, u)
	suite.NoError(err)
	suite.NotEmpty(token)
}

func (suite *authServiceTestSuite) TestRegisterUserAlreadyExists() {
	u := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().NewUser(gomock.Any(), gomock.Any()).Return(nil, user.ErrDuplicateLogin)

	token, err := suite.svc.Register(ctx, u)
	suite.ErrorIs(err, user.ErrDuplicateLogin)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestRegisterUnexpectedError() {
	u := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().NewUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error"))

	token, err := suite.svc.Register(ctx, u)
	suite.Error(err)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestLoginCorrectCredentials() {
	credentials := user.User{
		Login:    "user",
		Password: "password",
	}
	password, _ := user.HashPassword("password")
	u := user.User{
		ID:       1,
		Login:    "user",
		Password: password,
	}
	ctx := context.TODO()

	suite.store.EXPECT().GetUserByLogin(gomock.Any(), "user").Return(&u, nil)

	token, err := suite.svc.Login(ctx, credentials)
	suite.NoError(err)
	suite.NotEmpty(token)
}

func (suite *authServiceTestSuite) TestLoginIncorrectPassword() {
	credentials := user.User{
		Login:    "user",
		Password: "password",
	}
	password, _ := user.HashPassword("password2")
	u := user.User{
		Login:    "user",
		Password: password,
	}
	ctx := context.TODO()

	suite.store.EXPECT().GetUserByLogin(gomock.Any(), "user").Return(&u, nil)

	token, err := suite.svc.Login(ctx, credentials)
	suite.ErrorIs(err, user.ErrInvalidCredentials)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestLoginUserNotExists() {
	credentials := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().GetUserByLogin(gomock.Any(), "user").Return(nil, nil)

	token, err := suite.svc.Login(ctx, credentials)
	suite.ErrorIs(err, user.ErrInvalidCredentials)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestLoginUnexpectedError() {
	credentials := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().GetUserByLogin(gomock.Any(), "user").Return(nil, errors.New("unexpected error"))

	token, err := suite.svc.Login(ctx, credentials)
	suite.Error(err)
	suite.Empty(token)
}

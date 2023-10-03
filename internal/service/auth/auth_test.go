package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophermart/internal/entity/user"
	"github.com/k1nky/gophermart/internal/service/auth/mock"
	"github.com/stretchr/testify/suite"
)

type authServiceTestSuite struct {
	suite.Suite
	store *mock.MockStorage
	svc   *Service
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(authServiceTestSuite))
}

func (suite *authServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.Suite.T())
	suite.store = mock.NewMockStorage(ctrl)
	suite.svc = New("secret", 3*time.Hour, suite.store)
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

	suite.store.EXPECT().NewUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("duplicate"))
	suite.store.EXPECT().IsUniqueViolation(gomock.Any()).Return(true)

	token, err := suite.svc.Register(ctx, u)
	suite.ErrorIs(err, ErrDuplicateLoginError)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestRegisterUnexpectedError() {
	u := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().NewUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("unexpected error"))
	suite.store.EXPECT().IsUniqueViolation(gomock.Any()).Return(false)

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

	suite.store.EXPECT().GetUser(gomock.Any(), "user").Return(&u, nil)

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

	suite.store.EXPECT().GetUser(gomock.Any(), "user").Return(&u, nil)

	token, err := suite.svc.Login(ctx, credentials)
	suite.ErrorIs(err, ErrInvalidCredentials)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestLoginUserNotExists() {
	credentials := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().GetUser(gomock.Any(), "user").Return(nil, nil)

	token, err := suite.svc.Login(ctx, credentials)
	suite.ErrorIs(err, ErrInvalidCredentials)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestLoginUnexpectedError() {
	credentials := user.User{
		Login:    "user",
		Password: "password",
	}
	ctx := context.TODO()

	suite.store.EXPECT().GetUser(gomock.Any(), "user").Return(nil, errors.New("unexpected error"))

	token, err := suite.svc.Login(ctx, credentials)
	suite.Error(err)
	suite.Empty(token)
}

func (suite *authServiceTestSuite) TestParseToken() {
	claims := user.PrivateClaims{
		Login: "user",
	}
	token, err := suite.svc.GenerateToken(claims)
	suite.NoError(err)
	got, err := suite.svc.ParseToken(token)
	suite.NoError(err)
	suite.Equal(claims, got)
}

func (suite *authServiceTestSuite) TestParseExpiredToken() {
	claims := user.PrivateClaims{
		Login: "user",
	}
	suite.svc.tokenExpiration = 1 * time.Second
	token, err := suite.svc.GenerateToken(claims)
	suite.NoError(err)
	time.Sleep(3 * time.Second)
	got, err := suite.svc.ParseToken(token)
	suite.Error(err)
	suite.Empty(got)
}

func (suite *authServiceTestSuite) TestParseInvalidToken() {
	got, err := suite.svc.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYyODIxMDgsIkxvZ2luIjoidXNlciJ9.44K4rEcXS1bvyQY8h-TomgkKCC6Yysf44nl7O3n0KUI_invalid")
	suite.Error(err)
	suite.Empty(got)
}

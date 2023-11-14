package auth

import (
	"time"

	"github.com/k1nky/gophermart/internal/entity/user"
)

func (suite *authServiceTestSuite) TestAuthorize() {
	claims := user.PrivateClaims{
		Login: "user",
	}
	token, err := suite.svc.GenerateToken(claims)
	suite.NoError(err)
	got, err := suite.svc.Authorize(token)
	suite.NoError(err)
	suite.Equal(claims, got)
}

func (suite *authServiceTestSuite) TestAuthorizeExpiredToken() {
	claims := user.PrivateClaims{
		Login: "user",
	}
	suite.svc.tokenExpiration = 1 * time.Second
	token, err := suite.svc.GenerateToken(claims)
	suite.NoError(err)
	time.Sleep(3 * time.Second)
	got, err := suite.svc.Authorize(token)
	suite.ErrorIs(err, user.ErrUnathorized)
	suite.Empty(got)
}

func (suite *authServiceTestSuite) TestAuthorizeInvalidToken() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTYyODIxMDgsIkxvZ2luIjoidXNlciJ9.44K4rEcXS1bvyQY8h-TomgkKCC6Yysf44nl7O3n0KUI_invalid"
	got, err := suite.svc.Authorize(token)

	suite.ErrorIs(err, user.ErrUnathorized)
	suite.Empty(got)
}

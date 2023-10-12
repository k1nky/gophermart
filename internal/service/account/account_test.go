package account

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophermart/internal/service/account/mock"
	"github.com/stretchr/testify/suite"
)

type accountServiceTestSuite struct {
	suite.Suite
	store *mock.Mockstorage
}

func TestAccountService(t *testing.T) {
	suite.Run(t, new(accountServiceTestSuite))
}

func (suite *accountServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.store = mock.NewMockstorage(ctrl)
}

package database

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type adapterTestSuite struct {
	suite.Suite
	a *Adapter
}

func openTestDB() (*Adapter, error) {
	a := New()
	if err := a.Open(context.TODO(), "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"); err != nil {
		return nil, err
	}
	a.Exec(`delete from users; insert into users(user_id, login, password) values (1, 'u1', 'p1')`)
	return a, nil
}

func (suite *adapterTestSuite) SetupTest() {
	if shouldSkipDBTest(suite.T()) {
		return
	}
	var err error
	if suite.a, err = openTestDB(); err != nil {
		suite.FailNow(err.Error())
		return
	}
}

func shouldSkipDBTest(t *testing.T) bool {
	return false
	if len(os.Getenv("TEST_DB_READY")) == 0 {
		t.Skip()
		return true
	}
	return false
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(adapterTestSuite))
}

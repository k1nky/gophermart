package database

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func openTestDB() (*Adapter, error) {
	a := New()
	if err := a.Open(context.TODO(), "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"); err != nil {
		return nil, err
	}
	return a, nil
}

func shouldSkipDBTest(t *testing.T) bool {
	// TODO:
	// return false
	if len(os.Getenv("TEST_DB_READY")) == 0 {
		t.Skip()
		return true
	}
	return false
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(usersTestSuite))
	suite.Run(t, new(ordersTestSuite))
}

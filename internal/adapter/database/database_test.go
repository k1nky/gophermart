package database

import (
	"context"
	"testing"
)

func TestT(t *testing.T) {
	db := New()
	db.Open(context.TODO(), "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable")
}

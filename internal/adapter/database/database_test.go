package database

import "testing"

func TestT(t *testing.T) {
	db := New()
	db.Open("postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable")
}

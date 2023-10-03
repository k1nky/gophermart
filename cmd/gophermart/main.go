package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k1nky/gophermart/internal/adapter/database"
	"github.com/k1nky/gophermart/internal/adapter/http"
	"github.com/k1nky/gophermart/internal/service/auth"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	store := database.New()
	if err := store.Open("postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"); err != nil {
		panic(err.Error())
	}
	authService := auth.New("secret", 3*time.Hour, store)
	http.New(ctx, "", 8080, authService)

	<-ctx.Done()
}

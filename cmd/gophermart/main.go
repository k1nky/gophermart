package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/k1nky/gophermart/internal/adapter/accrual"
	"github.com/k1nky/gophermart/internal/adapter/database"
	"github.com/k1nky/gophermart/internal/adapter/http"
	"github.com/k1nky/gophermart/internal/service/account"
	accural "github.com/k1nky/gophermart/internal/service/accrual"
	"github.com/k1nky/gophermart/internal/service/auth"
)

func main() {

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	store := database.New()
	if err := store.Open("postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable"); err != nil {
		panic(err.Error())
	}
	authService := auth.New("secret", 3*time.Hour, store)
	account := account.New(store)
	accrualClient := accrual.New("http://localhost:8082")
	accrual := accural.New(store, accrualClient)
	accrual.Process(ctx)
	http.New(ctx, "", 8080, authService, account)

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

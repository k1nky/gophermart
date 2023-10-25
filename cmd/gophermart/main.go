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
	"github.com/k1nky/gophermart/internal/config"
	"github.com/k1nky/gophermart/internal/logger"
	"github.com/k1nky/gophermart/internal/service/account"
	accural "github.com/k1nky/gophermart/internal/service/accrual"
	"github.com/k1nky/gophermart/internal/service/auth"
)

func main() {

	cfg := config.Config{}
	if err := config.ParseConfig(&cfg); err != nil {
		panic(err)
	}
	log := logger.New()
	log.SetLevel("DEBUG")

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	store := database.New()
	if err := store.Open(ctx, cfg.DarabaseURI); err != nil {
		panic(err.Error())
	}
	authService := auth.New("secret", 3*time.Hour, store)
	account := account.New(store)
	accrualClient := accrual.New(cfg.AccrualSystemAddress)
	accrual := accural.New(store, accrualClient, log)
	accrual.Process(ctx)
	http.New(ctx, string(cfg.RunAddress), authService, account)

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

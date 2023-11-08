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

const (
	DefaultSecret          = "secret"
	DefaultTokenExpiration = 3 * time.Hour
)

func main() {
	cfg := config.Config{}
	if err := config.Parse(&cfg); err != nil {
		panic(err)
	}
	log := logger.New()
	log.SetLevel(cfg.LogLevel)
	log.SetLevel("debug")
	log.Debugf("config: %+v", cfg)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	run(ctx, cfg, log)

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}

func run(ctx context.Context, cfg config.Config, log *logger.Logger) {
	store := database.New()
	if err := store.Open(ctx, cfg.DarabaseURI); err != nil {
		log.Errorf("failed opening db: %v", err)
		return
	}
	authService := auth.New(DefaultSecret, DefaultTokenExpiration, store, log)
	account := account.New(store, log)
	accrualClient := accrual.New(cfg.AccrualSystemAddress)
	accrual := accural.New(store, accrualClient, log)
	accrual.Process(ctx)
	http.New(ctx, string(cfg.RunAddress), authService, account, log)
}

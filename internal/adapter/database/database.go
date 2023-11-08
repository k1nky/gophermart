package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	DefaultMaxKeepaliveConnections = 10
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Adapter struct {
	*sql.DB
}

func New() *Adapter {
	return &Adapter{}
}

// Открывает новое подключение к базе
func (a *Adapter) Open(ctx context.Context, dsn string) (err error) {
	if a.DB, err = sql.Open("pgx", dsn); err != nil {
		return
	}
	a.DB.SetMaxIdleConns(DefaultMaxKeepaliveConnections)
	a.DB.SetMaxOpenConns(DefaultMaxKeepaliveConnections)
	err = a.Initialize(dsn)
	if err == nil {
		// закрывает соединение при отмене контекста
		go func() {
			<-ctx.Done()
			a.Close()
		}()
	}
	return err
}

// Применяет миграции
func (a *Adapter) Initialize(dsn string) error {
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	// отсутствие изменений при миграции не считаем ошибкой
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func (a *Adapter) hasUniqueViolationError(err error) bool {
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		if pgerrcode.IsIntegrityConstraintViolation(pgerr.Code) {
			return true
		}
	}
	return false
}

// Закрывает подключение к базе
func (a *Adapter) Close() error {
	return a.DB.Close()
}

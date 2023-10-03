package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/k1nky/gophermart/internal/entity"
)

const (
	MaxKeepaliveConnections = 10
)

var (
	ErrUniqueViolation = errors.New("duplicate key value")
)

type Adapter struct {
	*sql.DB
}

func New() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Open(dsn string) (err error) {
	if a.DB, err = sql.Open("pgx", dsn); err != nil {
		return
	}
	a.DB.SetMaxIdleConns(MaxKeepaliveConnections)
	a.DB.SetMaxOpenConns(MaxKeepaliveConnections)
	return a.Initialize(dsn)
}

func (a *Adapter) Initialize(dsn string) error {
	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	return err
}

func (a *Adapter) GetUser(ctx context.Context, login string) (*entity.User, error) {
	u := &entity.User{
		Login: login,
	}
	row := a.QueryRowContext(ctx, `SELECT password FROM users WHERE login=$1`, login)
	if err := row.Err(); err != nil {
		return nil, err
	}
	if err := row.Scan(&u.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (a *Adapter) NewUser(ctx context.Context, u *entity.User) error {

	row := a.QueryRowContext(ctx, `
		INSERT INTO users AS u (login, password)
		VALUES ($1, $2)
		RETURNING u.user_id
	`, u.Login, u.Password)
	if err := row.Err(); err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgerr.Code) {
				return fmt.Errorf("%s %w", u.Login, ErrUniqueViolation)
			}
		}
		return err
	}
	if err := row.Scan(&u.ID); err != nil {
		return err
	}
	return nil
}

func (a *Adapter) IsUniqueViolation(err error) bool {
	return errors.Is(err, ErrUniqueViolation)
}

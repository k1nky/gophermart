package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/k1nky/gophermart/internal/entity/order"
	"github.com/k1nky/gophermart/internal/entity/user"
)

const (
	MaxKeepaliveConnections = 10
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func (a *Adapter) GetUser(ctx context.Context, login string) (*user.User, error) {
	u := &user.User{
		Login: login,
	}

	const query = `SELECT password FROM users WHERE login=$1`
	row := a.QueryRowContext(ctx, query, login)
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

func (a *Adapter) NewUser(ctx context.Context, u user.User) (*user.User, error) {

	const query = `
		INSERT INTO users AS u (login, password)
		VALUES ($1, $2)
		RETURNING u.user_id
	`

	row := a.QueryRowContext(ctx, query, u.Login, u.Password)
	if err := row.Err(); err != nil {
		if a.hasUniqueViolationError(err) {
			return nil, fmt.Errorf("%s %w", u.Login, user.ErrDuplicateLogin)
		}
		return nil, err
	}
	if err := row.Scan(&u.ID); err != nil {
		return nil, err
	}
	return &u, nil
}

func (a *Adapter) NewOrder(ctx context.Context, u user.User, o order.Order) (*order.Order, error) {
	const query = `
		INSERT INTO orders AS o (user_id, number, status, uploaded_at)
		VALUES ($1, $2, 'NEW', NOW())
		RETURNING o.order_id, o.uploaded_at
	`
	row := a.QueryRowContext(ctx, query, u.ID, o.Number)
	if err := row.Err(); err != nil {
		if a.hasUniqueViolationError(err) {
			return nil, fmt.Errorf("%s %w", o.Number, order.ErrDuplicateOrderError)
		}
		return nil, err
	}
	if err := row.Scan(&o.ID, &o.UploadedAt); err != nil {
		return nil, err
	}
	o.Status = order.StatusNew
	return &o, nil

}

func (a *Adapter) GetOrdersByStatus(ctx context.Context, statuses []order.OrderStatus) ([]*order.Order, error) {
	const query = `
		SELECT
			order_id, number, status, accrual, uploaded_at
		FROM orders
		WHERE status in ($1)
	`
	rows, err := a.QueryContext(ctx, query, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]*order.Order, 0)
	for rows.Next() {
		o := &order.Order{}
		if err := rows.Scan(&o.ID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
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

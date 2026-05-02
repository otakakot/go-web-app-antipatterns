package repository

import (
	"context"
	"database/sql"

	"github.com/otakakot/go-web-app-antipatterns/internal/app"
)

func MigrateDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			points INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS user_discounts (
			user_id INT PRIMARY KEY REFERENCES users(id),
			next_order_discount INT NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS audit_log (
			id SERIAL PRIMARY KEY,
			log TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

type db interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

var _ app.UserRepository = (*PostgresUserRepository)(nil)

type PostgresUserRepository struct {
	db db
}

func NewPostgresUserRepository(db db) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) GetPoints(ctx context.Context, userID int) (int, error) {
	row := r.db.QueryRowContext(ctx, "SELECT points FROM users WHERE id = $1 FOR UPDATE", userID)

	var points int
	err := row.Scan(&points)
	if err != nil {
		return 0, err
	}

	return points, nil
}

func (r *PostgresUserRepository) TakePoints(ctx context.Context, userID int, points int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE users SET points = points - $1 WHERE id = $2", points, userID)
	return err
}

var _ app.DiscountRepository = (*PostgresDiscountRepository)(nil)

type PostgresDiscountRepository struct {
	db db
}

func NewPostgresDiscountRepository(db db) *PostgresDiscountRepository {
	return &PostgresDiscountRepository{
		db: db,
	}
}

func (r *PostgresDiscountRepository) AddDiscount(ctx context.Context, userID int, discount int) error {
	_, err := r.db.ExecContext(ctx, "UPDATE user_discounts SET next_order_discount = next_order_discount + $1 WHERE user_id = $2", discount, userID)
	return err
}

package repository

import (
	"database/sql"
	"errors"

	"github.com/otakakot/go-web-app-antipatterns/internal/app"
)

var _ app.Provider = (*TransactionProvider)(nil)

type TransactionProvider struct {
	db *sql.DB
}

func NewTransactionProvider(db *sql.DB) *TransactionProvider {
	return &TransactionProvider{
		db: db,
	}
}

func (p *TransactionProvider) Transact(txFunc func(adapters app.Adapters) error) error {
	return runInTx(p.db, func(tx *sql.Tx) error {
		adapters := app.Adapters{
			UserRepository:     NewPostgresUserRepository(tx),
			DiscountRepository: NewPostgresDiscountRepository(tx),
		}

		return txFunc(adapters)
	})
}

func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}

package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type dbTransaction struct {
	db *sqlx.DB
}

// NewDBTransaction is to initialize mysql mysqlTransaction.
func NewDBTransaction(db *sqlx.DB) TransactionProvider {
	return &dbTransaction{
		db: db,
	}
}

// NewTransaction is to begin new sql transaction.
func (d *dbTransaction) NewTransaction(ctx context.Context, opts *sql.TxOptions) (TxProvider, error) {
	return d.db.BeginTxx(ctx, opts)
}

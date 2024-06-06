package repository

import (
	"context"
	"database/sql"
	"staffinc/internal/model/entity"
)

type UserRepoProvider interface {
	FindUserByEmail(ctx context.Context, email string) (user entity.User, err error)
	InsertUser(ctx context.Context, tx TxProvider, user entity.User) (int64, error)
}

type GeneratorLinkProvider interface {
	InsertGeneratorLink(ctx context.Context, tx TxProvider, userId int64) error
	GetGeneratorLinkByCode(ctx context.Context, tx TxProvider, code string) (entity.GeneratorLink, error)
	IncrementCount(ctx context.Context, tx TxProvider, code string) error
}

type TxProvider interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type QueryProvider interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

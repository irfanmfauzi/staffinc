package repository

import (
	"context"
	"staffinc/internal/model/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type generatorLink struct {
	baseRepo
}

func NewGeneratorLink(db *sqlx.DB) generatorLink {
	return generatorLink{
		baseRepo: baseRepo{db: db},
	}
}

func (g *generatorLink) InsertGeneratorLink(ctx context.Context, tx TxProvider, userId int64, code string, expiredAt time.Time) error {
	query := "INSERT INTO generator_links (user_id,code,expired_at) VALUES ($1,$2,$3)"

	_, err := g.DB(tx).ExecContext(ctx, query, userId, code, expiredAt)
	if err != nil {
		return err
	}

	return nil
}

func (g *generatorLink) GetGeneratorLinkByUserId(ctx context.Context, userId int64) ([]entity.GeneratorLink, error) {
	query := "SELECT * FROM generator_links WHERE user_id = $1"
	result := []entity.GeneratorLink{}

	err := g.db.SelectContext(ctx, &result, query, userId)

	if err != nil {
		return result, err
	}

	return result, nil

}

func (g *generatorLink) LockGetGeneratorLinkByCode(ctx context.Context, tx TxProvider, code string) (entity.GeneratorLink, error) {
	query := "SELECT * FROM generator_links WHERE code = $1 FOR UPDATE"

	result := entity.GeneratorLink{}
	err := g.DB(tx).GetContext(ctx, &result, query, code)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (g *generatorLink) IncrementCount(ctx context.Context, tx TxProvider, code string) error {
	query := "UPDATE generator_links SET count_access = count_access + 1 WHERE code = $1"

	_, err := g.DB(tx).ExecContext(ctx, query, code)
	if err != nil {
		return err
	}

	return nil
}

package repository

import (
	"context"
	"staffinc/internal/model/entity"
	"time"

	"github.com/google/uuid"
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

func (g *generatorLink) InsertGeneratorLink(ctx context.Context, tx TxProvider, userId int64) error {
	code := uuid.New()
	query := "INSERT INTO generator_links (user_id,code,expired_at) VALUES ($1,$2,$3)"

	_, err := g.DB(tx).ExecContext(ctx, query, userId, code, time.Now().Add(time.Hour*24))
	if err != nil {
		return err
	}

	return nil
}

func (g *generatorLink) GetGeneratorLinkByCode(ctx context.Context, tx TxProvider, code string) (entity.GeneratorLink, error) {
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

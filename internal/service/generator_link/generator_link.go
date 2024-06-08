package generatorlinkService

import (
	"context"
	"log/slog"
	"staffinc/internal/model/entity"
	errorX "staffinc/internal/model/error"
	"staffinc/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GeneratorLinkServiceProvider interface {
	GenerateLink(ctx context.Context) errorX.Error
	GetLink(ctx context.Context) (generatorLink []entity.GeneratorLink, err error)
}

type GeneratorLinkServiceConfig struct {
	Db                *sqlx.DB
	GeneratorLinkRepo repository.GeneratorLinkProvider
}

type generatorLinkService struct {
	db                *sqlx.DB
	generatorLinkRepo repository.GeneratorLinkProvider
}

func NewGenerateLinkService(cfg GeneratorLinkServiceConfig) generatorLinkService {
	return generatorLinkService{
		db:                cfg.Db,
		generatorLinkRepo: cfg.GeneratorLinkRepo,
	}
}

func (g *generatorLinkService) GenerateLink(ctx context.Context) errorX.Error {
	user := ctx.Value("user").(map[string]interface{})
	userId := int64(user["id"].(float64))
	role := user["role"].(string)

	if role != "generator" {
		return errorX.New(errorX.ERROR_CODE_FORBIDDEN_GENERATE_LINK)
	}

	tx, err := g.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "Error", err)
		return errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
	}
	defer tx.Rollback()

	err = g.generatorLinkRepo.InsertGeneratorLink(ctx, tx, userId, uuid.New().String(), time.Now().Add(time.Hour*24*7).UTC())
	if err != nil {
		slog.Error("failed to insert generator link", "Error", err)
		return errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
	}
	tx.Commit()
	return errorX.Error{}
}

func (g *generatorLinkService) GetLink(ctx context.Context) (generatorLink []entity.GeneratorLink, err error) {
	user := ctx.Value("user").(map[string]interface{})
	userId := int64(user["id"].(float64))

	generatorLink, err = g.generatorLinkRepo.GetGeneratorLinkByUserId(ctx, userId)
	if err != nil {
		return generatorLink, err
	}

	return generatorLink, nil

}

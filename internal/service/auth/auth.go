package authService

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"staffinc/internal/model"
	"staffinc/internal/model/entity"
	errorX "staffinc/internal/model/error"
	"staffinc/internal/model/request"
	"staffinc/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthServiceProvider interface {
	Login(ctx context.Context, request request.LoginRequest) (string, errorX.Error)

	Register(ctx context.Context, request request.RegisterRequest, code string) errorX.Error
}

type AuthServiceConfig struct {
	Db                *sqlx.DB
	UserRepo          repository.UserRepoProvider
	GeneratorLinkRepo repository.GeneratorLinkProvider
}

type authService struct {
	db                *sqlx.DB
	userRepo          repository.UserRepoProvider
	generatorLinkRepo repository.GeneratorLinkProvider
}

func NewAuthService(cfg AuthServiceConfig) authService {
	return authService{
		db:                cfg.Db,
		userRepo:          cfg.UserRepo,
		generatorLinkRepo: cfg.GeneratorLinkRepo,
	}
}

func (a *authService) Login(ctx context.Context, request request.LoginRequest) (string, errorX.Error) {
	user, err := a.userRepo.FindUserByEmail(ctx, request.Email)
	if err != nil {
		code := errorX.ERROR_CODE_INTERNAL_SERVER

		if err == sql.ErrNoRows {
			code = errorX.ERROR_CODE_NOT_AUTHENTICATED
			slog.Info("Failed to Find User By Email", "Error", err)
		} else {
			slog.Error("Failed to Find User By Email", "Error", err)
		}

		return "", errorX.New(code)
	}

	if user.Password != request.Password {
		return "", errorX.New(errorX.ERROR_CODE_NOT_AUTHENTICATED)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaim{
		User: entity.User{
			Id:    user.Id,
			Email: user.Email,
			Role:  user.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24).UTC()),
		},
	})

	tokenString, err := token.SignedString([]byte("REPLACE_THIS_WITH_ENV_VAR_SECRET"))
	if err != nil {
		slog.Error("Failed to Signed String", "Error", err)
		return "", errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
	}

	return tokenString, errorX.Error{}
}

func (a *authService) Register(ctx context.Context, req request.RegisterRequest, code string) errorX.Error {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("Failed to Begin Transaction", "Error", err)
		return errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
	}
	defer tx.Rollback()

	if code != "" {
		generatorLink, err := a.generatorLinkRepo.LockGetGeneratorLinkByCode(ctx, tx, code)
		if err != nil {
			slog.Info("Failed to Get Generator Link By Code", "Error", err)
			code := errorX.ERROR_CODE_INTERNAL_SERVER
			if err == sql.ErrNoRows {
				code = errorX.ERROR_CODE_REFERAL_NOT_FOUND_OR_EXPIRED
				err = errors.New("Code not found or expired")
			}
			return errorX.New(code)
		}
		if generatorLink.ExpiredAt.Before(time.Now().UTC()) {
			return errorX.New(errorX.ERROR_CODE_REFERAL_NOT_FOUND_OR_EXPIRED)
		}
	}

	userId, err := a.userRepo.InsertUser(ctx, tx, entity.User{Email: req.Email, Password: req.Password, Role: req.Role})
	if err != nil {
		slog.Error("Failed to Insert User", "Error", err)
		return errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
	}

	if req.Role == "generator" {
		err = a.generatorLinkRepo.InsertGeneratorLink(ctx, tx, userId, uuid.New().String(), time.Now().Add(24*time.Hour).UTC())
		if err != nil {
			slog.Error("Failed to Insert Generator Link", "Error", err)
			return errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
		}
	} else {
		err = a.generatorLinkRepo.IncrementCount(ctx, tx, code)
		if err != nil {
			return errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to Commit", "Error", err)
		return errorX.New(errorX.ERROR_CODE_INTERNAL_SERVER)
	}
	return errorX.Error{}
}

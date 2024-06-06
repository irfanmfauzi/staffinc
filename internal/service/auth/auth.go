package authService

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"staffinc/internal/model"
	"staffinc/internal/model/entity"
	"staffinc/internal/model/request"
	"staffinc/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type AuthServiceProvider interface {
	Login(ctx context.Context, request request.LoginRequest) (string, int, error)

	Register(ctx context.Context, request request.RegisterRequest, code string) (int, error)
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

func (a *authService) Login(ctx context.Context, request request.LoginRequest) (string, int, error) {
	user, err := a.userRepo.FindUserByEmail(ctx, request.Email)
	if err != nil {
		msg := ""
		code := http.StatusInternalServerError

		if err == sql.ErrNoRows {
			msg = "Email or Password is Wrong"
			code = http.StatusBadRequest
			slog.Info("Failed to Find User By Email", "Error", err)
		} else {
			msg = "Something Wrong With System"
			slog.Error("Failed to Find User By Email", "Error", err)
		}

		return "", code, errors.New(msg)
	}

	if user.Password != request.Password {
		return "", http.StatusBadRequest, errors.New("Email or Password is Wrong")
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
		return "", http.StatusInternalServerError, err
	}

	return tokenString, http.StatusOK, nil
}

func (a *authService) Register(ctx context.Context, req request.RegisterRequest, code string) (int, error) {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.Error("Failed to Begin Transaction", "Error", err)
		return http.StatusInternalServerError, err
	}
	defer tx.Rollback()

	if code != "" {
		generatorLink, err := a.generatorLinkRepo.GetGeneratorLinkByCode(ctx, tx, code)
		if err != nil {
			slog.Info("Failed to Get Generator Link By Code", "Error", err)
			httpCode := http.StatusInternalServerError
			if err == sql.ErrNoRows {
				httpCode = http.StatusNotFound
				err = errors.New("Code not found or expired")
			}
			return httpCode, err
		}
		if generatorLink.ExpiredAt.Before(time.Now().UTC()) {
			return http.StatusNotFound, errors.New("Code not found or expired")
		}
	}

	userId, err := a.userRepo.InsertUser(ctx, tx, entity.User{Email: req.Email, Password: req.Password, Role: req.Role})
	if err != nil {
		slog.Error("Failed to Insert User", "Error", err)
		return http.StatusInternalServerError, err
	}

	if req.Role == "generator" {
		err = a.generatorLinkRepo.InsertGeneratorLink(ctx, tx, userId)
		if err != nil {
			slog.Error("Failed to Insert Generator Link", "Error", err)
			return http.StatusInternalServerError, err
		}
	} else {
		err = a.generatorLinkRepo.IncrementCount(ctx, tx, code)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("Failed to Commit", "Error", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

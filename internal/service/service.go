package service

import (
	"fmt"
	"log"
	"os"
	"staffinc/internal/repository"
	authService "staffinc/internal/service/auth"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	UserService authService.AuthServiceProvider
}

var (
	database        = os.Getenv("DB_DATABASE")
	password        = os.Getenv("DB_PASSWORD")
	username        = os.Getenv("DB_USERNAME")
	port            = os.Getenv("DB_PORT")
	host            = os.Getenv("DB_HOST")
	serviceInstance Service
)

func New() Service {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepo(db)
	generatorLinkRepo := repository.NewGeneratorLink(db)
	userService := authService.NewAuthService(authService.AuthServiceConfig{
		UserRepo:          &userRepo,
		GeneratorLinkRepo: &generatorLinkRepo,
		Db:                db,
	})

	return Service{
		UserService: &userService,
	}
}

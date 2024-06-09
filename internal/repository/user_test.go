package repository_test

import (
	"context"
	"staffinc/internal/model/entity"
	"staffinc/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_userRepo_FindUserByEmail(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name     string
		args     args
		wantUser entity.User
		wantErr  error
	}{
		{
			name: "best_case",
			args: args{
				ctx:   ctx,
				email: "someemail@mail.com",
			},
			wantUser: entity.User{
				Id:       1,
				Email:    "someemail@mail.com",
				Password: "password",
				Role:     "generator",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDb := sqlx.NewDb(db, "sqlmock")

			u := repository.NewUserRepo(sqlxDb)

			query := "SELECT id,email,password,role FROM users where email = $1"

			mockeExpectExec := mock.ExpectQuery(query).WithArgs(tt.args.email)
			if tt.wantErr == nil {
				rows := sqlmock.NewRows([]string{"id", "email", "password", "role"})
				rows.AddRow(tt.wantUser.Id, tt.wantUser.Email, tt.wantUser.Password, tt.wantUser.Role)
				mockeExpectExec.WillReturnRows(rows)
			} else {
				mockeExpectExec.WillReturnError(tt.wantErr)
			}

			got, gotErr := u.FindUserByEmail(tt.args.ctx, tt.args.email)
			assert.Equal(t, tt.wantUser, got)
			assert.Equal(t, tt.wantErr, gotErr)

		})
	}
}

func Test_userRepo_InsertUser(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx  context.Context
		user entity.User
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr error
	}{
		{
			name: "best_case",
			args: args{
				ctx: ctx,
				user: entity.User{
					Email:    "mail@mail.com",
					Password: "password",
					Role:     "generator",
				},
			},
			want:    1,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			sqlxDb := sqlx.NewDb(db, "sqlmock")

			g := repository.NewUserRepo(sqlxDb)

			mock.ExpectBegin()

			query := "INSERT INTO users (email, password, role) VALUES ($1,$2,$3) RETURNING id"

			mockeExpectExec := mock.ExpectQuery(query).WithArgs(tt.args.user.Email, tt.args.user.Password, tt.args.user.Role)
			if tt.wantErr == nil {
				rows := sqlmock.NewRows([]string{"id"})
				rows.AddRow(int64(1))
				mockeExpectExec.WillReturnRows(rows)
			} else {
				mockeExpectExec.WillReturnError(tt.wantErr)
			}

			tx, err := sqlxDb.BeginTxx(tt.args.ctx, nil)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when creating transaction", err)
			}
			got, gotErr := g.InsertUser(tt.args.ctx, tx, tt.args.user)

			assert.Equal(t, tt.wantErr, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

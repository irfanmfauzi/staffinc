package repository_test

import (
	"context"
	"staffinc/internal/model/entity"
	"staffinc/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_generatorLink_LockGetGeneratorLinkByCode(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    entity.GeneratorLink
		wantErr error
	}{
		{
			name: "best_case",
			args: args{
				ctx:  ctx,
				code: "sample code",
			},
			want: entity.GeneratorLink{
				Id:          1,
				UserId:      1,
				Code:        "sample code",
				ExpiredAt:   time.Now().UTC(),
				CountAccess: 0,
			},
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

			g := repository.NewGeneratorLink(sqlxDb)

			mock.ExpectBegin()

			query := "SELECT * FROM generator_links WHERE code = $1 FOR UPDATE"

			mockeExpectExec := mock.ExpectQuery(query).WithArgs(tt.args.code)
			if tt.wantErr == nil {
				rows := sqlmock.NewRows([]string{"id", "user_id", "code", "expired_at", "count_access"})
				rows.AddRow(tt.want.Id, tt.want.UserId, tt.want.Code, tt.want.ExpiredAt, tt.want.CountAccess)
				mockeExpectExec.WillReturnRows(rows)
			} else {
				mockeExpectExec.WillReturnError(tt.wantErr)
			}

			tx, err := sqlxDb.BeginTxx(tt.args.ctx, nil)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when creating transaction", err)
			}

			got, gotErr := g.LockGetGeneratorLinkByCode(tt.args.ctx, tx, tt.args.code)

			assert.Equal(t, tt.wantErr, gotErr)
			assert.Equal(t, tt.want, got)

		})
	}
}

func Test_generatorLink_IncrementCount(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		code string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "best_case",
			args: args{
				ctx:  ctx,
				code: "sample code",
			},
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

			g := repository.NewGeneratorLink(sqlxDb)

			mock.ExpectBegin()

			query := "UPDATE generator_links SET count_access = count_access + 1 WHERE code = $1"

			mockeExpectExec := mock.ExpectExec(query).WithArgs(tt.args.code)
			if tt.wantErr == nil {
				mockeExpectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockeExpectExec.WillReturnError(tt.wantErr)
			}

			tx, err := sqlxDb.BeginTxx(tt.args.ctx, nil)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when creating transaction", err)
			}
			gotErr := g.IncrementCount(tt.args.ctx, tx, tt.args.code)

			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

func Test_generatorLink_InsertGeneratorLink(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		userId    int64
		code      string
		expiredAt time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "best_case",
			args: args{
				ctx:       ctx,
				userId:    int64(1),
				code:      "sample code",
				expiredAt: time.Now().Add(24 * time.Hour).UTC(),
			},
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

			g := repository.NewGeneratorLink(sqlxDb)

			mock.ExpectBegin()

			query := "INSERT INTO generator_links (user_id,code,expired_at) VALUES ($1,$2,$3)"

			mockeExpectExec := mock.ExpectExec(query).WithArgs(tt.args.userId, tt.args.code, tt.args.expiredAt)
			if tt.wantErr == nil {
				mockeExpectExec.WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mockeExpectExec.WillReturnError(tt.wantErr)
			}

			tx, err := sqlxDb.BeginTxx(tt.args.ctx, nil)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when creating transaction", err)
			}
			gotErr := g.InsertGeneratorLink(tt.args.ctx, tx, tt.args.userId, tt.args.code, tt.args.expiredAt)

			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}

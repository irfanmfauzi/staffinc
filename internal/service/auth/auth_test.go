package authService_test

import (
	"context"
	"staffinc/internal/model/entity"
	errorX "staffinc/internal/model/error"
	"staffinc/internal/model/request"
	"staffinc/internal/repository"
	"staffinc/internal/repository/mocks"
	authService "staffinc/internal/service/auth"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_authService_Login(t *testing.T) {
	ctx := context.Background()
	user := entity.User{
		Id:       1,
		Email:    "mail@mail.com",
		Password: "password",
		Role:     "generator",
	}

	type fields struct {
		userRepo          repository.UserRepoProvider
		generatorLinkRepo repository.GeneratorLinkProvider
	}
	type args struct {
		ctx     context.Context
		request request.LoginRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want1  errorX.Error
	}{
		{
			name: "best_case",
			fields: fields{
				userRepo: func() repository.UserRepoProvider {
					m := mocks.UserRepoProvider{}
					m.On("FindUserByEmail", mock.Anything, "mail@mail.com").Return(user, nil)
					return &m
				}(),
			},
			args: args{
				ctx: ctx,
				request: request.LoginRequest{
					Email:    "mail@mail.com",
					Password: "password",
				},
			},
			want1: errorX.Error{Error: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := authService.NewAuthService(authService.AuthServiceConfig{
				UserRepo:          tt.fields.userRepo,
				GeneratorLinkRepo: tt.fields.generatorLinkRepo,
			})

			_, errX := a.Login(tt.args.ctx, tt.args.request)
			assert.Equal(t, tt.want1, errX)
		})
	}
}

func Test_authService_Register(t *testing.T) {
	ctx := context.Background()
	userGenerator := entity.User{
		Email:    "tes@mail.com",
		Password: "password",
		Role:     "generator",
	}

	userContributor := entity.User{
		Email:    "tes@mail.com",
		Password: "password",
		Role:     "contributor",
	}

	type fields struct {
		transactionProvider repository.TransactionProvider
		userRepo            repository.UserRepoProvider
		generatorLinkRepo   repository.GeneratorLinkProvider
	}
	type args struct {
		ctx  context.Context
		req  request.RegisterRequest
		code string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errorX.Error
	}{
		{
			name: "best_case_generator",
			fields: fields{
				transactionProvider: func() *mocks.TransactionProvider {
					m := mocks.TransactionProvider{}
					tx := mocks.TxProvider{}
					m.On("NewTransaction", mock.Anything, mock.Anything).Return(&tx, nil)
					tx.On("Commit").Return(nil)
					tx.On("Rollback").Return(nil)
					return &m

				}(),
				userRepo: func() *mocks.UserRepoProvider {
					m := mocks.UserRepoProvider{}
					m.On("InsertUser", mock.Anything, mock.Anything, userGenerator).Return(int64(1), nil)

					return &m
				}(),
				generatorLinkRepo: func() *mocks.GeneratorLinkProvider {
					m := mocks.GeneratorLinkProvider{}

					m.On("InsertGeneratorLink", mock.Anything, mock.Anything, int64(1), mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("string")).Return(nil)
					return &m
				}(),
			},
			args: args{
				ctx: ctx,
				req: request.RegisterRequest{
					Email:    userGenerator.Email,
					Role:     userGenerator.Role,
					Password: userGenerator.Password,
				},
			},
			want: errorX.Error{Error: nil},
		},
		{
			name: "best_case_contributor",
			fields: fields{
				transactionProvider: func() *mocks.TransactionProvider {
					m := mocks.TransactionProvider{}
					tx := mocks.TxProvider{}
					m.On("NewTransaction", mock.Anything, mock.Anything).Return(&tx, nil)
					tx.On("Commit").Return(nil)
					tx.On("Rollback").Return(nil)
					return &m

				}(),
				userRepo: func() *mocks.UserRepoProvider {
					m := mocks.UserRepoProvider{}
					m.On("InsertUser", mock.Anything, mock.Anything, userContributor).Return(int64(1), nil)

					return &m
				}(),
				generatorLinkRepo: func() *mocks.GeneratorLinkProvider {
					m := mocks.GeneratorLinkProvider{}

					m.On("LockGetGeneratorLinkByCode", mock.Anything, mock.Anything, "sample-code").Return(entity.GeneratorLink{ExpiredAt: time.Now().Add(time.Hour * 24 * 7).UTC()}, nil)
					m.On("IncrementCount", mock.Anything, mock.Anything, "sample-code").Return(nil)
					return &m
				}(),
			},
			args: args{
				ctx:  ctx,
				code: "sample-code",
				req: request.RegisterRequest{
					Email:    userContributor.Email,
					Password: userContributor.Password,
					Role:     userContributor.Role,
				},
			},
			want: errorX.Error{Error: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := authService.NewAuthService(authService.AuthServiceConfig{
				TransactionProvider: tt.fields.transactionProvider,
				UserRepo:            tt.fields.userRepo,
				GeneratorLinkRepo:   tt.fields.generatorLinkRepo,
			})

			gotErr := a.Register(tt.args.ctx, tt.args.req, tt.args.code)
			assert.Equal(t, tt.want, gotErr)
		})
	}
}

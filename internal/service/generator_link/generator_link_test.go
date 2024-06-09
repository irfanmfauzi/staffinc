package generatorlinkService_test

import (
	"context"
	"staffinc/internal/model/entity"
	errorX "staffinc/internal/model/error"
	"staffinc/internal/repository"
	"staffinc/internal/repository/mocks"
	generatorlinkService "staffinc/internal/service/generator_link"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_generatorLinkService_GenerateLink(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		transactionProvider repository.TransactionProvider
		generatorLinkRepo   repository.GeneratorLinkProvider
	}
	type args struct {
		ctx    context.Context
		userId int64
		role   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   errorX.Error
	}{
		{
			name: "best_case",
			args: args{
				ctx:    ctx,
				userId: int64(1),
				role:   "generator",
			},
			fields: fields{
				transactionProvider: func() repository.TransactionProvider {
					m := mocks.TransactionProvider{}
					tx := mocks.TxProvider{}
					m.On("NewTransaction", mock.Anything, mock.Anything).Return(&tx, nil)
					tx.On("Rollback").Return(nil)
					tx.On("Commit").Return(nil)
					return &m
				}(),
				generatorLinkRepo: func() repository.GeneratorLinkProvider {
					m := mocks.GeneratorLinkProvider{}

					m.On("InsertGeneratorLink", mock.Anything, mock.Anything, int64(1), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
					return &m
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generatorlinkService.NewGenerateLinkService(generatorlinkService.GeneratorLinkServiceConfig{
				TransactionProvider: tt.fields.transactionProvider,
				GeneratorLinkRepo:   tt.fields.generatorLinkRepo,
			})

			got := g.GenerateLink(tt.args.ctx, tt.args.userId, tt.args.role)
			assert.Equal(t, tt.want, got)

		})
	}
}

func Test_generatorLinkService_GetLink(t *testing.T) {
	ctx := context.Background()

	generatorLinkExample := []entity.GeneratorLink{
		{
			Id:          1,
			UserId:      1,
			Code:        "sample-code",
			ExpiredAt:   time.Now().UTC(),
			CountAccess: 1,
		},
	}

	type fields struct {
		generatorLinkRepo repository.GeneratorLinkProvider
	}
	type args struct {
		ctx    context.Context
		userId int64
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantGeneratorLink []entity.GeneratorLink
		wantErr           error
	}{
		{
			name: "best_case",
			args: args{
				ctx:    ctx,
				userId: int64(1),
			},
			fields: fields{
				generatorLinkRepo: func() repository.GeneratorLinkProvider {
					m := mocks.GeneratorLinkProvider{}

					m.On("GetGeneratorLinkByUserId", mock.Anything, int64(1)).Return(generatorLinkExample, nil)
					return &m
				}(),
			},
			wantGeneratorLink: generatorLinkExample,
			wantErr:           nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := generatorlinkService.NewGenerateLinkService(generatorlinkService.GeneratorLinkServiceConfig{
				GeneratorLinkRepo: tt.fields.generatorLinkRepo,
			})
			gotGeneratorLink, err := g.GetLink(tt.args.ctx, tt.args.userId)

			assert.Equal(t, tt.wantGeneratorLink, gotGeneratorLink)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

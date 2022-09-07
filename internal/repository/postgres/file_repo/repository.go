package file_repo

import (
	"context"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	"github.com/abdivasiyev/project_template/pkg/helpers"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type repo struct {
	querier storage.Querier
	log     logger.Logger
}

type Params struct {
	fx.In
	Querier storage.Querier
	Log     logger.Logger
}

func New(params Params) repository.File {
	return &repo{
		querier: params.Querier,
		log:     params.Log,
	}
}

func (r *repo) Create(ctx context.Context, request models.GetFileResponse) error {
	query := `insert into file (id, name, url, created_at) values ($1, $2, $3, now())`

	_, err := r.querier.Exec(ctx, query, request.FileID, request.FileName, request.FileURL)

	return helpers.ToCustomError(err)
}

func (r *repo) Get(ctx context.Context, id string) (models.GetFileResponse, error) {
	var response models.GetFileResponse

	query := `select id, name, url from file where deleted_at is null and id = $1`

	err := r.querier.QueryRow(ctx, query, id).Scan(
		&response.FileID,
		&response.FileName,
		&response.FileURL,
	)

	return response, helpers.ToCustomError(err)
}

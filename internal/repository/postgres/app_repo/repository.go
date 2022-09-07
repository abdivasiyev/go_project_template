package app_repo

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

func New(params Params) repository.App {
	return &repo{
		querier: params.Querier,
		log:     params.Log,
	}
}

func (r repo) Get(ctx context.Context) (models.GetAppVersionResponse, error) {
	query := `select version, title, description, force_update from app_version order by created_at desc limit 1`

	var response models.GetAppVersionResponse

	if err := r.querier.QueryRow(ctx, query).Scan(
		&response.Version,
		&response.Title,
		&response.Description,
		&response.ForceUpdate,
	); err != nil {
		return models.GetAppVersionResponse{}, helpers.ToCustomError(err)
	}

	return response, nil
}

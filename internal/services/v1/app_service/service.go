package app_service

import (
	"context"
	"errors"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewService)

type service struct {
	environment string
	log         logger.Logger
	sentry      sentry.Handler
	appRepo     repository.App
}

type Params struct {
	fx.In
	Config  config.Config
	Log     logger.Logger
	Sentry  sentry.Handler
	AppRepo repository.App
}

func NewService(params Params) v1.AppServiceV1 {
	return &service{
		environment: params.Config.GetString(config.EnvironmentKey),
		log:         params.Log,
		sentry:      params.Sentry,
		appRepo:     params.AppRepo,
	}
}

func (s service) GetVersion(ctx context.Context) (models.GetAppVersionResponse, error) {
	response, err := s.appRepo.Get(ctx)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.log.Error("could not get app version", zap.Error(err))
			s.sentry.HandleError(err)
		}
	}

	return response, err
}

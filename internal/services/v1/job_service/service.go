package job_service

import (
	"context"
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewService)

type service struct {
	environment string
	uploadPath  string
	logger      logger.Logger
	sentry      sentry.Handler
}

type Params struct {
	fx.In
	Config config.Config
	Logger logger.Logger
	Sentry sentry.Handler
}

func NewService(params Params) v1.JobServiceV1 {
	s := &service{
		environment: params.Config.GetString(config.EnvironmentKey),
		uploadPath:  params.Config.GetString(config.UploadPathKey),
		logger:      params.Logger,
		sentry:      params.Sentry,
	}

	return s
}

func (s *service) ExampleJob(ctx context.Context) error {
	fmt.Println("I am example job")

	return nil
}

package main

import (
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/handler"
	"github.com/abdivasiyev/project_template/internal/job"
	"github.com/abdivasiyev/project_template/internal/middleware"
	postgresRepo "github.com/abdivasiyev/project_template/internal/repository/postgres"
	"github.com/abdivasiyev/project_template/internal/server"
	"github.com/abdivasiyev/project_template/internal/services"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/mailer"
	"github.com/abdivasiyev/project_template/pkg/router"
	"github.com/abdivasiyev/project_template/pkg/security"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"github.com/abdivasiyev/project_template/pkg/storage/postgres"
	"github.com/abdivasiyev/project_template/pkg/storage/redis"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		config.Module,
		logger.Module,
		router.Module,
		security.Module,
		sentry.Module,
		postgres.Module,
		redis.Module,
		handler.Module,
		postgresRepo.Module,
		server.Module,
		services.Module,
		middleware.Module,
		job.Module,
		mailer.Module,
	)

	app.Run()
}

package services

import (
	appV1 "github.com/abdivasiyev/project_template/internal/services/v1/app_service"
	authV1 "github.com/abdivasiyev/project_template/internal/services/v1/auth_service"
	fileV1 "github.com/abdivasiyev/project_template/internal/services/v1/file_service"
	jobV1 "github.com/abdivasiyev/project_template/internal/services/v1/job_service"
	middlewareV1 "github.com/abdivasiyev/project_template/internal/services/v1/middleware_service"
	roleV1 "github.com/abdivasiyev/project_template/internal/services/v1/role_service"
	userV1 "github.com/abdivasiyev/project_template/internal/services/v1/user_service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	authV1.Module,
	fileV1.Module,
	jobV1.Module,
	middlewareV1.Module,
	roleV1.Module,
	userV1.Module,
	appV1.Module,
)

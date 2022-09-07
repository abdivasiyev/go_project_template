package handler

import (
	handlerV1 "github.com/abdivasiyev/project_template/internal/handler/v1"
	appV1 "github.com/abdivasiyev/project_template/internal/handler/v1/app"
	authV1 "github.com/abdivasiyev/project_template/internal/handler/v1/auth"
	docV1 "github.com/abdivasiyev/project_template/internal/handler/v1/doc"
	fileV1 "github.com/abdivasiyev/project_template/internal/handler/v1/file"
	pprofV1 "github.com/abdivasiyev/project_template/internal/handler/v1/pprof"
	roleV1 "github.com/abdivasiyev/project_template/internal/handler/v1/role"
	userV1 "github.com/abdivasiyev/project_template/internal/handler/v1/user"
	"go.uber.org/fx"
)

var Module = fx.Options(
	authV1.Module,
	docV1.Module,
	fileV1.Module,
	pprofV1.Module,
	roleV1.Module,
	userV1.Module,
	appV1.Module,
	handlerV1.Module,
)

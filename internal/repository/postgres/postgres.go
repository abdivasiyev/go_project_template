package postgres

import (
	"github.com/abdivasiyev/project_template/internal/repository/postgres/app_repo"
	"github.com/abdivasiyev/project_template/internal/repository/postgres/file_repo"
	"github.com/abdivasiyev/project_template/internal/repository/postgres/permission_repo"
	"github.com/abdivasiyev/project_template/internal/repository/postgres/role_repo"
	"github.com/abdivasiyev/project_template/internal/repository/postgres/user_repo"
	"go.uber.org/fx"
)

var Module = fx.Options(
	file_repo.Module,
	permission_repo.Module,
	role_repo.Module,
	user_repo.Module,
	app_repo.Module,
)

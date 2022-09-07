package permission_repo

import (
	"context"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	"github.com/abdivasiyev/project_template/internal/types"
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

func New(params Params) repository.Permission {
	return &repo{
		querier: params.Querier,
		log:     params.Log,
	}
}

func (r *repo) GetPermissionByUserAndPathAndMethod(ctx context.Context, userID, path, method string) (models.GetPermissionResponse, error) {
	permissions, err := r.findBy(ctx, "(:user_id=ANY(select ur.user_id from user_role ur where ur.role_id = rp.role_id) and p.path = :path and p.method = :method) or p.allow_all = true", types.M{
		"user_id": userID,
		"path":    path,
		"method":  method,
	})

	if err != nil {
		return models.GetPermissionResponse{}, helpers.ToCustomError(err)
	}

	if len(permissions) > 0 {
		return permissions[0], nil
	}

	return models.GetPermissionResponse{}, models.ErrNotFound
}

func (r *repo) GetByRole(ctx context.Context, roleID string) ([]models.GetPermissionResponse, error) {
	return r.findBy(ctx, "r.id = :role_id", types.M{
		"role_id": roleID,
	})
}

func (r *repo) GetByUser(ctx context.Context, userID string) ([]models.GetPermissionResponse, error) {
	return r.findBy(ctx, ":user_id=ANY(select ur.user_id from user_role ur where ur.role_id = rp.role_id)", types.M{
		"user_id": userID,
	})
}

func (r *repo) findBy(ctx context.Context, statement string, params types.M) ([]models.GetPermissionResponse, error) {
	var permissions []models.GetPermissionResponse

	query := `
		select p.id,
			   p.alias,
			   p.name,
			   p.path,
			   p.method,
			   coalesce(p.query_param, ''),
			   coalesce(p.query_param_value, '')
		from permission p
				 join role_permission rp on p.id = rp.permission_id
				 join role r on rp.role_id = r.id and r.deleted_at is null
		where p.deleted_at is null and ` + statement

	stmt, err := r.querier.PrepareNamed(ctx, query)
	if err != nil {
		return nil, helpers.ToCustomError(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(params)
	if err != nil {
		return nil, helpers.ToCustomError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var perm models.GetPermissionResponse

		if err = rows.Scan(
			&perm.ID,
			&perm.Alias,
			&perm.Name,
			&perm.Path,
			&perm.Method,
			&perm.QueryParam,
			&perm.QueryParamValue,
		); err != nil {
			return nil, helpers.ToCustomError(err)
		}

		permissions = append(permissions, perm)
	}

	return permissions, nil
}

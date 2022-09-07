package role_repo

import (
	"context"
	"database/sql"
	"github.com/abdivasiyev/project_template/internal/repository"
	"go.uber.org/fx"
	"sync"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/pkg/helpers"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"github.com/lib/pq"
	"go.uber.org/zap"
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

func New(params Params) repository.Role {
	return &repo{
		querier: params.Querier,
		log:     params.Log,
	}
}

func (r *repo) Delete(ctx context.Context, id string) error {
	query := `update role set deleted_at=current_timestamp where is_basic=false and id=$1`

	result, err := r.querier.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return models.ErrNotFound
	}

	return nil
}

func (r *repo) Create(ctx context.Context, req models.CreateRoleRequest) error {
	tx, err := r.querier.Begin(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := `insert into role (id, alias, name, description, created_at) values ($1, $2, $3, $4, current_timestamp)`

	if _, err = tx.Exec(query, req.ID, req.Alias, req.Name, req.Description); err != nil {
		_ = tx.Rollback()
		return helpers.ToCustomError(err)
	}

	query = `with d as (delete from role_permission where role_id = $1) insert into role_permission (role_id, permission_id) (select $1, unnest($2::uuid[]))`

	if _, err = tx.Exec(query, req.ID, pq.Array(req.Permissions)); err != nil {
		_ = tx.Rollback()
		return helpers.ToCustomError(err)
	}

	return helpers.ToCustomError(tx.Commit())
}

func (r *repo) Update(ctx context.Context, req models.CreateRoleRequest) error {
	tx, err := r.querier.Begin(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	query := `update role set alias = $2, name = $3, description = $4, updated_at = current_timestamp where id = $1`

	if _, err = tx.Exec(query, req.ID, req.Alias, req.Name, req.Description); err != nil {
		_ = tx.Rollback()
		return helpers.ToCustomError(err)
	}

	query = `with d as (delete from role_permission where role_id = $1) insert into role_permission (role_id, permission_id) (select $1, unnest($2::uuid[]))`

	if _, err = tx.Exec(query, req.ID, pq.Array(req.Permissions)); err != nil {
		_ = tx.Rollback()
		return helpers.ToCustomError(err)
	}

	return helpers.ToCustomError(tx.Commit())
}

func (r *repo) GetAll(ctx context.Context, req models.GetAllRoleRequest) (models.GetAllRoleResponse, error) {
	var response models.GetAllRoleResponse

	queryCount := `
		select
			count(1)
		from role where deleted_at is null
	`

	if err := r.querier.QueryRow(ctx, queryCount).Scan(&response.Count); err != nil {
		return models.GetAllRoleResponse{}, helpers.ToCustomError(err)
	}

	query := `
		select
			id,
			name,
			alias,
			description
		from role where deleted_at is null
		offset $1
		limit $2
	`

	offset, limit := helpers.NormalizePagination(req.Page, req.Limit)

	rows, err := r.querier.Query(ctx, query, offset, limit)
	if err != nil {
		return models.GetAllRoleResponse{}, helpers.ToCustomError(err)
	}

	for rows.Next() {
		var (
			role        models.GetRoleResponse
			description sql.NullString
		)

		if err = rows.Scan(&role.ID, &role.Name, &role.Alias, &description); err != nil {
			return models.GetAllRoleResponse{}, helpers.ToCustomError(err)
		}

		role.Description = description.String

		response.Roles = append(response.Roles, role)
	}

	return response, nil
}

func (r *repo) IsAdmin(ctx context.Context, userID string) (bool, error) {
	var count int
	query := `select count(1) from user_role where user_id=$1 and role_id='e715df60-2384-4f6a-bbd6-65126b14f6b2'`

	if err := r.querier.QueryRow(ctx, query, userID).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *repo) Get(ctx context.Context, id string) (models.GetRoleResponse, error) {
	var (
		role        models.GetRoleResponse
		description sql.NullString
	)

	query := `
		select
			id,
			name,
			alias,
			description
		from role where id=$1 and deleted_at is null
	`

	if err := r.querier.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Alias, &description); err != nil {
		return models.GetRoleResponse{}, helpers.ToCustomError(err)
	}

	role.Description = description.String

	query = `
		select p.id,
			   p.alias,
			   p.name,
			   p.path,
			   p.method
		from permission p
				 join role_permission rp on p.id = rp.permission_id
		where p.deleted_at is null
		  and rp.role_id = $1
		order by p.sequence
	`

	rows, err := r.querier.Query(ctx, query, id)
	if err != nil {
		return models.GetRoleResponse{}, helpers.ToCustomError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var perm models.GetPermissionResponse

		if err = rows.Scan(&perm.ID, &perm.Alias, &perm.Name, &perm.Path, &perm.Method); err != nil {
			return models.GetRoleResponse{}, helpers.ToCustomError(err)
		}

		role.Permissions = append(role.Permissions, perm)
	}

	return role, nil
}

func (r *repo) GetModules(ctx context.Context) (models.GetModulesResponse, error) {
	var (
		mu          sync.Mutex
		wg          sync.WaitGroup
		modules     map[string]models.GetModuleResponse
		groups      map[string][]models.GetPermissionGroupResponse
		permissions map[string][]models.GetPermissionResponse
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		modulesMap, err := r.getModulesMap(ctx)
		if err != nil {
			r.log.Error("could not get modules map", zap.Error(err))
		}

		mu.Lock()
		modules = modulesMap
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		groupsMap, err := r.getGroupsMap(ctx)
		if err != nil {
			r.log.Error("could not get groups map", zap.Error(err))
		}

		mu.Lock()
		groups = groupsMap
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()
		permissionsMap, err := r.getPermissionsMap(ctx)
		if err != nil {
			r.log.Error("could not get permissions map", zap.Error(err))
		}

		mu.Lock()
		permissions = permissionsMap
		mu.Unlock()
	}()

	wg.Wait()

	for moduleID, module := range modules {
		moduleGroups := groups[moduleID]

		for i := range moduleGroups {
			moduleGroups[i].Permissions = permissions[moduleGroups[i].ID]
		}
		groups[moduleID] = moduleGroups
		module.Groups = moduleGroups
		modules[moduleID] = module
	}

	var response = models.GetModulesResponse{Count: len(modules)}

	for _, modul := range modules {
		response.Modules = append(response.Modules, modul)
	}

	return response, nil
}

func (r *repo) getModulesMap(ctx context.Context) (map[string]models.GetModuleResponse, error) {
	modules := make(map[string]models.GetModuleResponse)

	query := `
		select id,
			   name,
			   alias
		from permission_module
		where deleted_at is null
	`

	rows, err := r.querier.Query(ctx, query)
	if err != nil {
		return modules, helpers.ToCustomError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var module models.GetModuleResponse

		if err = rows.Scan(
			&module.ID,
			&module.Name,
			&module.Alias,
		); err != nil {
			return modules, helpers.ToCustomError(err)
		}

		modules[module.ID] = module
	}

	return modules, nil
}

func (r *repo) getGroupsMap(ctx context.Context) (map[string][]models.GetPermissionGroupResponse, error) {
	groups := make(map[string][]models.GetPermissionGroupResponse)

	query := `
		select id,
			   name,
			   alias,
               module_id
		from permission_group
		where deleted_at is null
	`

	rows, err := r.querier.Query(ctx, query)
	if err != nil {
		return groups, helpers.ToCustomError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var group models.GetPermissionGroupResponse

		if err = rows.Scan(
			&group.ID,
			&group.Name,
			&group.Alias,
			&group.ModuleID,
		); err != nil {
			return groups, helpers.ToCustomError(err)
		}

		groups[group.ModuleID] = append(groups[group.ModuleID], group)
	}

	return groups, nil
}

func (r *repo) getPermissionsMap(ctx context.Context) (map[string][]models.GetPermissionResponse, error) {
	permissions := make(map[string][]models.GetPermissionResponse)

	query := `
		select p.id,
			   p.name,
			   p.alias,
			   p.path,
			   p.method,
			   coalesce(p.query_param, ''),
			   coalesce(p.query_param_value, ''),
			   pgr.group_id
		from permission p
		join permission_group_relation pgr on p.id = pgr.permission_id
		where p.deleted_at is null and p.allow_all = false
	`

	rows, err := r.querier.Query(ctx, query)
	if err != nil {
		return permissions, helpers.ToCustomError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var permission models.GetPermissionResponse

		if err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Alias,
			&permission.Path,
			&permission.Method,
			&permission.QueryParam,
			&permission.QueryParamValue,
			&permission.GroupID,
		); err != nil {
			return permissions, helpers.ToCustomError(err)
		}

		permissions[permission.GroupID] = append(permissions[permission.GroupID], permission)
	}

	return permissions, nil
}

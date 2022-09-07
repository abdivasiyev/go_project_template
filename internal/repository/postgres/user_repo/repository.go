package user_repo

import (
	"context"
	"database/sql"
	"github.com/abdivasiyev/project_template/internal/repository"
	"go.uber.org/fx"
	"time"

	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/pkg/helpers"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"github.com/pkg/errors"
)

const (
	RoleHR                  = "hr"
	RoleAdmin               = "admin"
	RoleSafety              = "safety"
	RoleDispatchOrientation = "dispatch_orientation"
	RoleFleet               = "fleet"
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

func New(params Params) repository.User {
	return &repo{
		querier: params.Querier,
		log:     params.Log,
	}
}

func (r *repo) Create(ctx context.Context, req models.CreateUserRequest) error {
	tx, err := r.querier.Begin(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "could not begin transaction")
	}

	query := `insert into "user" (id, company_id, username, password_hash, first_name, last_name, phone, created_at) values ($1, $2, $3, $4, $5, $6, $7, current_timestamp)`

	_, err = tx.Exec(
		query,
		req.ID,
		helpers.ToNullString(req.CompanyID),
		req.Username,
		req.Password,
		helpers.ToNullString(req.FirstName),
		helpers.ToNullString(req.LastName),
		helpers.ToNullString(req.Phone),
	)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "could not create user")
	}

	query = `with d as (delete from user_role where user_id = $1) insert into user_role (user_id, role_id) values ($1, $2)`
	_, err = tx.Exec(
		query,
		req.ID,
		req.RoleID,
	)

	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "could not set user role")
	}

	return tx.Commit()
}

func (r *repo) Update(ctx context.Context, req models.UpdateUserRequest) error {
	tx, err := r.querier.Begin(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "could not begin transaction")
	}

	query := `
		update "user" set 
			company_id = $2,
			username = $3,
			first_name = $4,
			last_name = $5,
			password_hash = coalesce(nullif($6, ''), password_hash),
			phone = $7,
			updated_at = current_timestamp
		where id = $1 and deleted_at is null
	`

	_, err = tx.Exec(
		query,
		req.ID,
		helpers.ToNullString(req.CompanyID),
		req.Username,
		helpers.ToNullString(req.FirstName),
		helpers.ToNullString(req.LastName),
		req.NewPassword,
		helpers.ToNullString(req.Phone),
	)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "could not update user")
	}

	query = `with d as (delete from user_role where user_id = $1) insert into user_role (user_id, role_id) values ($1, $2)`
	_, err = tx.Exec(
		query,
		req.ID,
		req.RoleID,
	)

	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "could not set user role")
	}

	return tx.Commit()
}

func (r *repo) UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) error {
	tx, err := r.querier.Begin(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "could not begin transaction")
	}

	query := `
		update "user" set 
			username = $2,
			first_name = $3,
			last_name = $4,
			password_hash = coalesce(nullif($5, ''), password_hash),
			image_id = $6,
			updated_at = current_timestamp
		where id = $1 and deleted_at is null
	`

	_, err = tx.Exec(
		query,
		req.ID,
		req.Username,
		helpers.ToNullString(req.FirstName),
		helpers.ToNullString(req.LastName),
		req.NewPassword,
		helpers.ToNullString(req.ImageID),
	)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "could not update user")
	}

	return tx.Commit()
}

func (r *repo) Delete(ctx context.Context, id string) error {
	query := `
		update "user" set deleted_at = current_timestamp where id = $1 and deleted_at is null
	`

	_, err := r.querier.Exec(
		ctx,
		query,
		id,
	)
	return errors.Wrap(err, "could not delete user")
}

func (r *repo) GetByUsername(ctx context.Context, username string) (models.GetUserResponse, error) {
	return r.findByOne(ctx, "WHERE u.username = :username AND u.deleted_at is null", map[string]interface{}{
		"username": username,
	})
}

func (r *repo) Get(ctx context.Context, id string) (models.GetUserResponse, error) {
	return r.findByOne(ctx, "WHERE u.id = :id AND u.deleted_at is null", map[string]interface{}{
		"id": id,
	})
}

func (r *repo) GetAll(ctx context.Context, req models.GetAllUsersRequest) (models.GetAllUsersResponse, error) {
	var (
		statement = `WHERE u.deleted_at is null`
		params    = make(map[string]interface{})
	)

	if req.Search != "" {
		params["search"] = req.Search

		statement += ` AND u.username ilike '%' || :search || '%'`
	}

	params["offset"], params["limit"] = helpers.NormalizePagination(req.Page, req.Limit)
	return r.find(ctx, statement, params)
}

func (r *repo) find(ctx context.Context, statement string, params map[string]interface{}) (models.GetAllUsersResponse, error) {
	var response models.GetAllUsersResponse

	queryCount := `
		SELECT
			count(1)
		FROM "user" u
		LEFT JOIN company c ON c.id = u.company_id AND c.deleted_at is null
	` + statement

	stmtCount, err := r.querier.PrepareNamed(ctx, queryCount)
	if err != nil {
		return response, errors.Wrap(err, "could not prepare named context")
	}
	defer stmtCount.Close()

	if err = stmtCount.QueryRow(params).Scan(&response.Count); err != nil {
		return response, helpers.ToCustomError(err)
	}

	query := `
		WITH r AS ( select id, name, alias, description from "role" where deleted_at is null),
			ur AS (select role_id, user_id from user_role)
		SELECT
			u.id,
			u.username,
			u.password_hash,
			u.first_name,
			u.last_name,
			u.created_at,
			u.updated_at,
			c.id,
			c.name,
			c.created_at,
			c.updated_at,
			(select r.id from r where r.id=(select ur.role_id from ur where ur.user_id=u.id limit 1)),
			(select r.name from r where r.id=(select ur.role_id from ur where ur.user_id=u.id limit 1)),
			(select r.alias from r where r.id=(select ur.role_id from ur where ur.user_id=u.id limit 1)),
			(select r.description from r where r.id=(select ur.role_id from ur where ur.user_id=u.id limit 1))
		FROM "user" u
		LEFT JOIN company c ON c.id = u.company_id AND c.deleted_at is null
	` + statement + `
		ORDER BY u.created_at DESC
		OFFSET :offset LIMIT :limit
	`

	stmt, err := r.querier.PrepareNamed(ctx, query)
	if err != nil {
		return response, errors.Wrap(err, "could not prepare named context")
	}
	defer stmt.Close()

	rows, err := stmt.Query(params)
	if err != nil {
		return response, errors.Wrap(err, "could not query with params")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user                                          models.GetUserResponse
			firstName, lastName                           sql.NullString
			createdAt                                     time.Time
			updatedAt, companyCreatedAt, companyUpdatedAt sql.NullTime
			companyID, companyName                        sql.NullString
			roleID, roleName, roleAlias, roleDescription  sql.NullString
		)

		if err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.PasswordHash,
			&firstName,
			&lastName,
			&createdAt,
			&updatedAt,
			&companyID,
			&companyName,
			&companyCreatedAt,
			&companyUpdatedAt,
			&roleID,
			&roleName,
			&roleAlias,
			&roleDescription,
		); err != nil {
			return response, errors.Wrap(err, "could not scan rows")
		}

		user.FirstName = firstName.String
		user.LastName = lastName.String
		user.CreatedAt = helpers.TimeToString(createdAt, config.DateFormat, true)
		user.UpdatedAt = helpers.TimeToString(updatedAt.Time, config.DateFormat, updatedAt.Valid)
		user.Company = models.GetCompanyResponse{
			ID:        companyID.String,
			Name:      companyName.String,
			CreatedAt: helpers.TimeToString(companyCreatedAt.Time, config.DateFormat, companyCreatedAt.Valid),
			UpdatedAt: helpers.TimeToString(companyUpdatedAt.Time, config.DateFormat, companyUpdatedAt.Valid),
		}
		user.Role = models.GetRoleResponse{
			ID:          roleID.String,
			Alias:       roleAlias.String,
			Name:        roleName.String,
			Description: roleDescription.String,
		}

		response.Users = append(response.Users, user)
	}

	return response, nil
}

func (r *repo) findByOne(ctx context.Context, statement string, params map[string]interface{}) (models.GetUserResponse, error) {
	var (
		user                                          models.GetUserResponse
		firstName, lastName, imageID                  sql.NullString
		createdAt                                     time.Time
		updatedAt, companyCreatedAt, companyUpdatedAt sql.NullTime
		companyID, companyName                        sql.NullString
	)

	query := `
		SELECT
			u.id,
			u.username,
			u.password_hash,
			u.first_name,
			u.last_name,
			u.created_at,
			u.updated_at,
			u.image_id,
			c.id,
			c.name,
			c.created_at,
			c.updated_at,
			r.id,
			r.alias,
			r.name
		FROM "user" u
		LEFT JOIN company c ON c.id = u.company_id AND c.deleted_at is null
		JOIN user_role ur ON ur.user_id = u.id
		JOIN "role" r ON r.id = ur.role_id
		` + statement

	stmt, err := r.querier.PrepareNamed(ctx, query)
	if err != nil {
		return user, errors.Wrap(err, "could not prepare named context")
	}
	defer stmt.Close()

	if err = stmt.QueryRow(params).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&firstName,
		&lastName,
		&createdAt,
		&updatedAt,
		&imageID,
		&companyID,
		&companyName,
		&companyCreatedAt,
		&companyUpdatedAt,
		&user.Role.ID,
		&user.Role.Alias,
		&user.Role.Name,
	); err != nil {
		return user, helpers.ToCustomError(err)
	}

	user.FirstName = firstName.String
	user.LastName = lastName.String
	user.ImageID = imageID.String
	user.CreatedAt = helpers.TimeToString(createdAt, config.DateFormat, true)
	user.UpdatedAt = helpers.TimeToString(updatedAt.Time, config.DateFormat, updatedAt.Valid)
	user.Company = models.GetCompanyResponse{
		ID:        companyID.String,
		Name:      companyName.String,
		CreatedAt: helpers.TimeToString(companyCreatedAt.Time, config.DateFormat, companyCreatedAt.Valid),
		UpdatedAt: helpers.TimeToString(companyUpdatedAt.Time, config.DateFormat, companyUpdatedAt.Valid),
	}

	return user, nil
}

package repository

import (
	"context"

	"github.com/abdivasiyev/project_template/internal/models"
)

type App interface {
	Get(ctx context.Context) (models.GetAppVersionResponse, error)
}

type File interface {
	Create(ctx context.Context, request models.GetFileResponse) error
	Get(ctx context.Context, id string) (models.GetFileResponse, error)
}

type Permission interface {
	GetPermissionByUserAndPathAndMethod(ctx context.Context, userID, path, method string) (models.GetPermissionResponse, error)
	GetByRole(ctx context.Context, roleID string) ([]models.GetPermissionResponse, error)
	GetByUser(ctx context.Context, userID string) ([]models.GetPermissionResponse, error)
}

type Role interface {
	Delete(ctx context.Context, id string) error
	IsAdmin(ctx context.Context, userID string) (bool, error)
	Get(ctx context.Context, id string) (models.GetRoleResponse, error)
	GetAll(ctx context.Context, req models.GetAllRoleRequest) (models.GetAllRoleResponse, error)
	Update(ctx context.Context, req models.CreateRoleRequest) error
	Create(ctx context.Context, req models.CreateRoleRequest) error
	GetModules(ctx context.Context) (models.GetModulesResponse, error)
}

// User provides user database functions
type User interface {
	Create(ctx context.Context, req models.CreateUserRequest) error
	Update(ctx context.Context, req models.UpdateUserRequest) error
	UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) error
	GetByUsername(ctx context.Context, username string) (models.GetUserResponse, error)
	Get(ctx context.Context, id string) (models.GetUserResponse, error)
	GetAll(ctx context.Context, req models.GetAllUsersRequest) (models.GetAllUsersResponse, error)
	Delete(ctx context.Context, id string) error
}

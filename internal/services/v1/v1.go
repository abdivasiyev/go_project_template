package v1

import (
	"context"
	"github.com/abdivasiyev/project_template/internal/models"
	"mime/multipart"
	"net/http"
	"time"
)

type AppServiceV1 interface {
	GetVersion(ctx context.Context) (models.GetAppVersionResponse, error)
}

type JobServiceV1 interface {
	ExampleJob(ctx context.Context) error
}

type MiddlewareServiceV1 interface {
	RecoverPanic(ctx context.Context, request *http.Request, err error, stack bool)
	Log(ctx context.Context, statusCode int, request *http.Request, clientIP string, startTime, endTime time.Time, errors []error, timeFormat string)
	HasAccess(ctx context.Context, userID, path, method string, fn func(queryParam string) string) error
	CheckAuth(ctx context.Context, token string) (models.GetUserResponse, error)
}

type UserServiceV1 interface {
	Create(ctx context.Context, req models.CreateUserRequest) (models.GetUserResponse, error)
	Update(ctx context.Context, req models.UpdateUserRequest) (models.GetUserResponse, error)
	UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) (models.GetUserResponse, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (models.GetUserResponse, error)
	GetAll(ctx context.Context, request models.GetAllUsersRequest) (models.GetAllUsersResponse, error)
}

type RoleServiceV1 interface {
	Update(ctx context.Context, req models.CreateRoleRequest) (models.GetRoleResponse, error)
	Create(ctx context.Context, req models.CreateRoleRequest) (models.GetRoleResponse, error)
	GetModules(ctx context.Context) (models.GetModulesResponse, error)
	Get(ctx context.Context, id string) (models.GetRoleResponse, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, req models.GetAllRoleRequest) (models.GetAllRoleResponse, error)
}

type AuthServiceV1 interface {
	ResetPassword(ctx context.Context, req models.ResetPasswordRequest) (models.SuccessResponse, error)
	Login(ctx context.Context, request models.LoginRequest) (models.AuthenticationResponse, error)
	Refresh(ctx context.Context, request models.RefreshTokenRequest) (models.AuthenticationResponse, error)
}

type FileServiceV1 interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (models.GetFileResponse, error)
	GetFile(ctx context.Context, id string) (models.GetFileResponse, string, error)
}

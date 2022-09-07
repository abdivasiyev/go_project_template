package role_service

import (
	"context"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"go.uber.org/zap"
	"time"

	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"go.uber.org/fx"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var Module = fx.Provide(NewService)

type service struct {
	environment          string
	log                  logger.Logger
	sentry               sentry.Handler
	roleRepository       repository.Role
	permissionRepository repository.Permission
	cache                storage.Cacher
}

type Params struct {
	fx.In
	Config               config.Config
	Log                  logger.Logger
	Sentry               sentry.Handler
	RoleRepository       repository.Role
	PermissionRepository repository.Permission
	Cache                storage.Cacher
}

func NewService(params Params) v1.RoleServiceV1 {
	return &service{
		environment:          params.Config.GetString(config.EnvironmentKey),
		log:                  params.Log,
		sentry:               params.Sentry,
		roleRepository:       params.RoleRepository,
		permissionRepository: params.PermissionRepository,
		cache:                params.Cache,
	}
}

func (s *service) Get(ctx context.Context, id string) (models.GetRoleResponse, error) {
	response, err := s.roleRepository.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get role", zap.Error(err), zap.String("roleID", id))
		}
		return models.GetRoleResponse{}, err
	}

	response.Permissions, err = s.permissionRepository.GetByRole(ctx, response.ID)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get permissions", zap.Error(err), zap.String("roleID", id))
	}

	return response, err
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.roleRepository.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.ErrForbidden
		}
		s.log.Error("could not delete role", zap.String("id", id))
		s.sentry.HandleError(err)
	}
	return err
}

func (s *service) GetAll(ctx context.Context, req models.GetAllRoleRequest) (models.GetAllRoleResponse, error) {
	response, err := s.roleRepository.GetAll(ctx, req)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get all roles", zap.Error(err), zap.Any("req", req))
		}
	}

	return response, err
}

func (s *service) Update(ctx context.Context, req models.CreateRoleRequest) (models.GetRoleResponse, error) {
	if err := s.roleRepository.Update(ctx, req); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not update role", zap.Error(err), zap.Any("req", req))
		return models.GetRoleResponse{}, errors.Wrap(err, "could not create role")
	}

	response, err := s.roleRepository.Get(ctx, req.ID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get role", zap.Error(err), zap.Any("req", req))
		}
	}

	return response, err
}

func (s *service) Create(ctx context.Context, req models.CreateRoleRequest) (models.GetRoleResponse, error) {
	req.ID = uuid.New().String()

	if err := s.roleRepository.Create(ctx, req); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not create role", zap.Error(err), zap.Any("req", req))
		return models.GetRoleResponse{}, errors.Wrap(err, "could not create role")
	}

	response, err := s.roleRepository.Get(ctx, req.ID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get role", zap.Error(err), zap.Any("req", req))
		}
	}

	return response, err
}

func (s *service) GetModules(ctx context.Context) (models.GetModulesResponse, error) {
	var (
		response models.GetModulesResponse
		err      error
		key      = "role:permission:modules"
	)

	if err = s.cache.GetObj(ctx, key, &response); err == nil {
		return response, nil
	}

	response, err = s.roleRepository.GetModules(ctx)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get role modules", zap.Error(err))
		}
		return models.GetModulesResponse{}, err
	}

	if err = s.cache.SetObj(ctx, key, response, 7*24*time.Hour); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get role modules", zap.Error(err))
	}

	return response, nil
}

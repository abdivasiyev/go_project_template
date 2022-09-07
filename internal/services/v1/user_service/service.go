package user_service

import (
	"context"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"go.uber.org/zap"

	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"go.uber.org/fx"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/security"
	"github.com/abdivasiyev/project_template/pkg/validator"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var Module = fx.Provide(NewService)

type service struct {
	environment          string
	log                  logger.Logger
	sentry               sentry.Handler
	userRepository       repository.User
	security             security.Handler
	permissionRepository repository.Permission
}

type Params struct {
	fx.In
	Config               config.Config
	Log                  logger.Logger
	Sentry               sentry.Handler
	UserRepository       repository.User
	Security             security.Handler
	PermissionRepository repository.Permission
}

func NewService(params Params) v1.UserServiceV1 {
	return &service{
		environment:          params.Config.GetString(config.EnvironmentKey),
		log:                  params.Log,
		sentry:               params.Sentry,
		userRepository:       params.UserRepository,
		security:             params.Security,
		permissionRepository: params.PermissionRepository,
	}
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.userRepository.Delete(ctx, id)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not delete user", zap.Error(err), zap.Any("userID", id))
		}
	}
	return err
}

func (s *service) UpdateProfile(ctx context.Context, req models.UpdateProfileRequest) (models.GetUserResponse, error) {
	newPasswordHash, err := s.validateUserForUpdate(ctx, req.ID, req.Username, req.NewPassword, req.OldPassword)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not validate profile", zap.Error(err), zap.Any("req", req))
		return models.GetUserResponse{}, errors.Wrap(err, "could not create user")
	}
	req.NewPassword = newPasswordHash

	if err = s.userRepository.UpdateProfile(ctx, req); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not update profile", zap.Error(err), zap.Any("req", req))
		return models.GetUserResponse{}, errors.Wrap(err, "could not create user")
	}

	response, err := s.userRepository.Get(ctx, req.ID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get user", zap.Error(err), zap.Any("req", req))
		}
	}
	return response, err
}

func (s *service) Update(ctx context.Context, req models.UpdateUserRequest) (models.GetUserResponse, error) {
	newPasswordHash, err := s.validateUserForUpdate(ctx, req.ID, req.Username, req.NewPassword, req.OldPassword)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not validate user", zap.Error(err), zap.Any("req", req))
		return models.GetUserResponse{}, errors.Wrap(err, "could not create user")
	}
	req.NewPassword = newPasswordHash

	if err = s.userRepository.Update(ctx, req); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not update user", zap.Error(err), zap.Any("req", req))
		return models.GetUserResponse{}, errors.Wrap(err, "could not create user")
	}

	response, err := s.userRepository.Get(ctx, req.ID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get user", zap.Error(err), zap.Any("req", req))
		}
	}
	return response, err
}

func (s *service) validateUserForUpdate(ctx context.Context, userID, username, oldPassword, newPassword string) (string, error) {
	if err := validator.ValidatePassword(oldPassword, newPassword); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not validate password", zap.Error(err))
		return "", err
	}

	userByUsername, err := s.userRepository.GetByUsername(ctx, username)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get user", zap.Error(err), zap.String("username", username))
			return "", errors.Wrap(err, "could not check username")
		}
	}

	if oldPassword != "" && userByUsername.ID != "" {
		valid, err := s.security.CompareHash(oldPassword, userByUsername.PasswordHash)
		if err != nil {
			s.sentry.HandleError(err)
			s.log.Error("could not compare hashes", zap.Error(err))
			return "", err
		}

		if !valid {
			return "", validator.NewValidationError("old_password", "incorrect password supplied for old_password")
		}
	}

	if newPassword != "" {
		newPassword, err = s.security.GenerateHash(newPassword)
		if err != nil {
			s.sentry.HandleError(err)
			s.log.Error("could not generate hash", zap.Error(err))
			return "", err
		}
	}

	if userByUsername.ID != "" && userID != userByUsername.ID {
		return "", validator.NewValidationError("username", "username already exists")
	}

	return newPassword, nil
}

func (s *service) Create(ctx context.Context, req models.CreateUserRequest) (models.GetUserResponse, error) {
	req.ID = uuid.New().String()

	_, err := s.userRepository.GetByUsername(ctx, req.Username)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get user", zap.Error(err), zap.Any("req", req))
			return models.GetUserResponse{}, errors.Wrap(err, "could not check username")
		}
	}

	req.Password, err = s.security.GenerateHash(req.Password)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not generate hash", zap.Error(err), zap.Any("req", req))
		return models.GetUserResponse{}, errors.Wrap(err, "could not generate hash")
	}

	if err = s.userRepository.Create(ctx, req); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not create user", zap.Error(err), zap.Any("req", req))
		return models.GetUserResponse{}, errors.Wrap(err, "could not create user")
	}

	response, err := s.userRepository.Get(ctx, req.ID)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get user", zap.Error(err), zap.Any("req", req))
		}
	}
	return response, err
}

func (s *service) Get(ctx context.Context, id string) (models.GetUserResponse, error) {
	response, err := s.userRepository.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get user", zap.Error(err), zap.String("userID", id))
		}
		return models.GetUserResponse{}, err
	}

	response.Role.Permissions, err = s.permissionRepository.GetByRole(ctx, response.Role.ID)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get role permissions", zap.Error(err), zap.String("roleID", response.Role.ID), zap.String("userID", id))
	}

	return response, err
}

func (s *service) GetAll(ctx context.Context, request models.GetAllUsersRequest) (models.GetAllUsersResponse, error) {
	response, err := s.userRepository.GetAll(ctx, request)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get user", zap.Error(err), zap.Any("req", request))
		}
	}
	return response, err
}

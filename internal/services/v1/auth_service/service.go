package auth_service

import (
	"context"
	"github.com/abdivasiyev/project_template/config"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/mailer"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"go.uber.org/fx"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/security"
	"github.com/abdivasiyev/project_template/pkg/security/jwt"
	customValidator "github.com/abdivasiyev/project_template/pkg/validator"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewService)

type service struct {
	environment          string
	log                  logger.Logger
	sentry               sentry.Handler
	userRepository       repository.User
	permissionRepository repository.Permission
	cache                storage.Cacher
	mailer               mailer.Mailer
	security             security.Handler
}

type Params struct {
	fx.In
	Config               config.Config
	Log                  logger.Logger
	Sentry               sentry.Handler
	UserRepository       repository.User
	PermissionRepository repository.Permission
	Cache                storage.Cacher
	Mailer               mailer.Mailer
	Security             security.Handler
}

func NewService(params Params) v1.AuthServiceV1 {
	return &service{
		environment:          params.Config.GetString(config.EnvironmentKey),
		log:                  params.Log,
		sentry:               params.Sentry,
		userRepository:       params.UserRepository,
		permissionRepository: params.PermissionRepository,
		security:             params.Security,
		cache:                params.Cache,
		mailer:               params.Mailer,
	}
}

func (s *service) Login(ctx context.Context, request models.LoginRequest) (models.AuthenticationResponse, error) {
	user, err := s.userRepository.GetByUsername(ctx, request.Username)

	if errors.Is(err, models.ErrNotFound) {
		s.log.Warn("user not found with username", zap.Any("username", request.Username))
		return models.AuthenticationResponse{}, customValidator.NewValidationError("username", "incorrect username or password")
	} else if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("error getting user", zap.Any("request", request), zap.Error(err))
		return models.AuthenticationResponse{}, errors.Wrap(err, "could not get user")
	}

	valid, err := s.security.CompareHash(request.Password, user.PasswordHash)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("error validating password", zap.Any("request", request), zap.Error(err))
		return models.AuthenticationResponse{}, errors.Wrap(err, "could not verify password")
	}

	if !valid {
		s.log.Warn("password for user is not valid", zap.Any("request", request))
		return models.AuthenticationResponse{}, customValidator.NewValidationError("username", "incorrect username or password")
	}

	accessToken, refreshToken, err := s.security.GenerateToken(user)

	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not generate token", zap.Any("request", request), zap.Error(err))
		return models.AuthenticationResponse{}, errors.Wrap(err, "could not generate token")
	}

	permissions, err := s.permissionRepository.GetByUser(ctx, user.ID)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get user permissions", zap.Error(err), zap.String("userID", user.ID))
	}

	return models.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Permissions:  permissions,
	}, nil
}

func (s *service) Refresh(ctx context.Context, request models.RefreshTokenRequest) (models.AuthenticationResponse, error) {
	tokenUser, err := s.security.VerifyToken(request.Token, true)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return models.AuthenticationResponse{}, customValidator.NewValidationError("token", "invalid refresh token")
		}

		if errors.Is(err, jwt.ErrExpiredToken) {
			return models.AuthenticationResponse{}, customValidator.NewValidationError("token", "expired refresh token")
		}

		s.sentry.HandleError(err)
		s.log.Error("could not verify token", zap.Error(err))
		return models.AuthenticationResponse{}, errors.Wrap(err, "could not verify token")
	}

	user, err := s.userRepository.GetByUsername(ctx, tokenUser.Username)

	if errors.Is(err, models.ErrNotFound) {
		return models.AuthenticationResponse{}, customValidator.NewValidationError("token", "user not exists")
	} else if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get user", zap.Error(err), zap.Any("request", request))
		return models.AuthenticationResponse{}, errors.Wrap(err, "could not get user")
	}

	accessToken, refreshToken, err := s.security.GenerateToken(user)

	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not generate token", zap.Error(err), zap.Any("request", request))
		return models.AuthenticationResponse{}, errors.Wrap(err, "could not generate token")
	}

	permissions, err := s.permissionRepository.GetByUser(ctx, user.ID)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get user permissions", zap.Error(err), zap.String("userID", user.ID))
	}

	return models.AuthenticationResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Permissions:  permissions,
	}, nil
}

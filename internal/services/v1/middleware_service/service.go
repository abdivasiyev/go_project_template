package middleware_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/helpers"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/security"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var Module = fx.Provide(New)

type service struct {
	environment          string
	log                  logger.Logger
	sentry               sentry.Handler
	security             security.Handler
	permissionRepository repository.Permission
	roleRepository       repository.Role
	cache                storage.Cacher
}

type Params struct {
	fx.In
	Config               config.Config
	Logger               logger.Logger
	Sentry               sentry.Handler
	PermissionRepository repository.Permission
	RoleRepository       repository.Role
	Security             security.Handler
	Cache                storage.Cacher
}

func New(params Params) v1.MiddlewareServiceV1 {
	return &service{
		environment:          params.Config.GetString(config.EnvironmentKey),
		log:                  params.Logger,
		sentry:               params.Sentry,
		security:             params.Security,
		permissionRepository: params.PermissionRepository,
		roleRepository:       params.RoleRepository,
		cache:                params.Cache,
	}
}

func (s *service) RecoverPanic(_ context.Context, request *http.Request, err error, stack bool) {
	if err == nil {
		return
	}
	// Check for a broken connection, as it is not really a
	// condition that warrants a panic stack trace.
	var brokenPipe bool
	if ne, isOpErr := err.(*net.OpError); isOpErr {
		if se, isSysCallErr := ne.Err.(*os.SyscallError); isSysCallErr {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
				strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}

	httpRequest, dumpErr := httputil.DumpRequest(request, true)
	if dumpErr != nil {
		s.sentry.HandleError(dumpErr)
		s.log.Error("could not dump request", zap.Error(dumpErr))
	}
	if brokenPipe {
		s.log.Error("unexpected client connection closed",
			zap.String("requestPath", request.URL.Path),
			zap.Any("error", err),
			zap.String("request", string(httpRequest)),
		)
		return
	}

	if stack {
		s.log.Error("[Recovery from panic]",
			zap.Time("time", time.Now()),
			zap.Any("error", err),
			zap.String("request", string(httpRequest)),
			zap.String("stack", string(debug.Stack())),
		)
	} else {
		s.log.Error("[Recovery from panic]",
			zap.Time("time", time.Now()),
			zap.Any("error", err),
			zap.String("request", string(httpRequest)),
		)
	}

	// send error to sentry
	s.sentry.HandleError(err)
}

func (s *service) Log(
	_ context.Context,
	statusCode int,
	request *http.Request,
	clientIP string,
	startTime, endTime time.Time,
	errors []error,
	timeFormat string,
) {
	latency := endTime.Sub(startTime)

	fields := []zapcore.Field{
		zap.Int("status", statusCode),
		zap.String("method", request.Method),
		zap.String("path", request.URL.RawPath),
		zap.String("query", request.URL.RawQuery),
		zap.String("ip", clientIP),
		zap.String("user-agent", request.UserAgent()),
		zap.String("latency", latency.String()),
		zap.Errors("errors", errors),
	}
	if timeFormat != "" {
		fields = append(
			fields,
			zap.String("startTime", startTime.Format(timeFormat)),
			zap.String("endTime", endTime.Format(timeFormat)),
		)
	}

	if len(errors) > 0 {
		for _, err := range errors {
			s.sentry.HandleError(err)
		}
	}

	s.logByStatusCode(statusCode, "[REQUEST]", fields...)
}

func (s *service) logByStatusCode(statusCode int, message string, fields ...zap.Field) {
	if statusCode >= http.StatusOK && statusCode < http.StatusBadRequest {
		s.log.Info(message, fields...)
		return
	}

	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		s.log.Warn(message, fields...)
		return
	}

	s.log.Error(message, fields...)
}

func (s *service) HasAccess(ctx context.Context, userID, path, method string, fn func(queryParam string) string) error {
	var permission models.GetPermissionResponse

	key := fmt.Sprintf("permission:%s:%s:%s", url.QueryEscape(path), method, userID)

	err := s.cache.GetObj(ctx, key, &permission)
	if err == nil {
		if helpers.IsEmpty(permission.QueryParam) {
			return nil
		}
		queryParamValue := fn(permission.QueryParam)

		if queryParamValue != permission.QueryParamValue {
			return models.ErrForbidden
		}
	}

	permission, err = s.permissionRepository.GetPermissionByUserAndPathAndMethod(ctx, userID, path, method)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not check access", zap.Error(err), zap.String("userId", userID), zap.String("path", path), zap.String("method", method))
			return err
		}

		if err = s.isAdmin(ctx, userID); err != nil {
			return err
		}
		if err = s.cache.SetObj(ctx, key, permission, 1*time.Hour); err != nil {
			s.sentry.HandleError(err)
			s.log.Error("could not save to cache", zap.Error(err))
		}
		return nil
	}

	if !helpers.IsEmpty(permission.QueryParam) {
		queryParamValue := fn(permission.QueryParam)
		if queryParamValue != permission.QueryParamValue {
			return models.ErrForbidden
		}
	}

	if err = s.cache.SetObj(ctx, key, permission, 30*time.Second); err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not save to cache", zap.Error(err))
	}

	return nil
}

func (s *service) isAdmin(ctx context.Context, userID string) error {
	isAdmin, err := s.roleRepository.IsAdmin(ctx, userID)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not check access", zap.Error(err), zap.String("userId", userID))
		return models.ErrForbidden
	}

	if isAdmin {
		return nil
	}

	return models.ErrForbidden
}

func (s *service) CheckAuth(ctx context.Context, token string) (models.GetUserResponse, error) {
	token = strings.TrimSpace(token)

	if token == "" {
		return models.GetUserResponse{}, models.ErrUnauthorized
	}

	tokens := strings.Split(token, " ")
	if len(tokens) != 2 {
		return models.GetUserResponse{}, models.ErrUnauthorized
	}

	token = tokens[1]

	user, err := s.security.VerifyToken(token, false)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not verify access token: %v", zap.Error(err))
		return models.GetUserResponse{}, err
	}

	return user, nil
}

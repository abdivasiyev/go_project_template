package middleware

import (
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Provide(New)

type Handler interface {
	HasAccess() gin.HandlerFunc
	BearerAuth() gin.HandlerFunc
	Log(timeFormat string, utc bool) gin.HandlerFunc
	LogWithConfig(conf *Config) gin.HandlerFunc
	RecoverWithLog(stack bool) gin.HandlerFunc
	WrapHttpHandler(f http.Handler) gin.HandlerFunc
}

type middleware struct {
	log     logger.Logger
	service v1.MiddlewareServiceV1
}

type Params struct {
	fx.In
	Log     logger.Logger
	Service v1.MiddlewareServiceV1
}

func New(params Params) Handler {
	return &middleware{
		log:     params.Log,
		service: params.Service,
	}
}

package pprof

import (
	"net/http/pprof"

	"github.com/abdivasiyev/project_template/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In
	Config config.Config
}

type Handler struct {
	environment string
}

func New(params Params) *Handler {
	return &Handler{environment: params.Config.GetString(config.EnvironmentKey)}
}

func (h *Handler) Cmdline() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.environment != config.Development {
			return
		}
		pprof.Cmdline(c.Writer, c.Request)
	}
}

func (h *Handler) Profile() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.environment != config.Development {
			return
		}
		pprof.Profile(c.Writer, c.Request)
	}
}

func (h *Handler) Symbol() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.environment != config.Development {
			return
		}
		pprof.Symbol(c.Writer, c.Request)
	}
}

func (h *Handler) Heap() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.environment != config.Development {
			return
		}
		heapHandler := pprof.Handler("heap")
		heapHandler.ServeHTTP(c.Writer, c.Request)
	}
}

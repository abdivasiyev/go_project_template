package doc

import (
	"net/http"

	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Handler struct {
	environment string
	log         logger.Logger
	title       string
	swaggerURL  string
	description string
}

type Params struct {
	fx.In
	Config config.Config
	Log    logger.Logger
}

func New(params Params) *Handler {
	return &Handler{
		environment: params.Config.GetString(config.EnvironmentKey),
		log:         params.Log,
		title:       params.Config.GetString(config.SpecTitle),
		swaggerURL:  params.Config.GetString(config.SpecUrl),
		description: params.Config.GetString(config.SpecDescription),
	}
}

func (h *Handler) Render() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":       h.title,
			"url":         h.swaggerURL,
			"description": h.description,
		})
	}
}

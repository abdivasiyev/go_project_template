package app

import (
	"github.com/abdivasiyev/project_template/config"
	serviceV1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
)

var Module = fx.Provide(NewHandler)

type Handler struct {
	environment string
	log         logger.Logger
	service     serviceV1.AppServiceV1
}

type Params struct {
	fx.In
	Config  config.Config
	Log     logger.Logger
	Service serviceV1.AppServiceV1
}

func NewHandler(params Params) *Handler {
	return &Handler{
		environment: params.Config.GetString(config.EnvironmentKey),
		log:         params.Log,
		service:     params.Service,
	}
}

// GetVersion godoc
// @Summary Returns application version
// @Description Returns application version
// @Accept  json
// @Produce  json
// @Success 200 {object} models.GetAppVersionResponse
// @Failure default {object} models.ErrorResponse
// @Tags app
// @Router /v1/app/version [get]
func (h *Handler) GetVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := h.service.GetVersion(c)
		if err != nil {
			h.log.Errorf("could not get app version: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get app version",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    resp,
			StatusCode: http.StatusOK,
		})
	}
}

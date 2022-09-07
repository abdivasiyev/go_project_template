package auth

import (
	"github.com/abdivasiyev/project_template/config"
	"go.uber.org/fx"
	"net/http"

	serviceV1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/response"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/gin-gonic/gin"
)

var Module = fx.Provide(NewHandler)

type Handler struct {
	environment string
	log         logger.Logger
	service     serviceV1.AuthServiceV1
}

type Params struct {
	fx.In
	Config  config.Config
	Log     logger.Logger
	Service serviceV1.AuthServiceV1
}

func NewHandler(params Params) *Handler {
	return &Handler{
		environment: params.Config.GetString(config.EnvironmentKey),
		log:         params.Log,
		service:     params.Service,
	}
}

// ResetPassword godoc
// @Summary ResetPassword resets user password
// @Description Returns ok if success
// @Accept  json
// @Produce  json
// @Param refreshTokenRequest body models.ResetPasswordRequest true "reset password request"
// @Success 200 {object} models.SuccessResponse
// @Failure default {object} models.ErrorResponse
// @Tags auth
// @Router /v1/auth/reset-password [post]
func (h *Handler) ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.ResetPasswordRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			h.log.Errorf("could not unmarshal json request: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not unmarshal json body",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		resp, err := h.service.ResetPassword(c, request)
		if err != nil {
			h.log.Errorf("could not reset user password: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not reset password",
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

// Refresh godoc
// @Summary Refresh access token to new one
// @Description Returns new access token response
// @Accept  json
// @Produce  json
// @Param refreshTokenRequest body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.AuthenticationResponse
// @Failure default {object} models.ErrorResponse
// @Tags auth
// @Router /v1/auth/refresh [post]
func (h *Handler) Refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.RefreshTokenRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			h.log.Errorf("could not unmarshal json request: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not unmarshal json body",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		resp, err := h.service.Refresh(c, request)
		if err != nil {
			h.log.Errorf("could not login user: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not generate token",
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

// Login godoc
// @Summary Authorize by username and password
// @Description Returns access token response
// @Accept  json
// @Produce  json
// @Param loginRequest body models.LoginRequest true "Auth credentials"
// @Success 200 {object} models.AuthenticationResponse
// @Failure default {object} models.ErrorResponse
// @Tags auth
// @Router /v1/auth/login [post]
func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			h.log.Errorf("could not unmarshal json request: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not unmarshal json body",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		resp, err := h.service.Login(c, request)
		if err != nil {
			h.log.Errorf("could not login user: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not generate token",
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

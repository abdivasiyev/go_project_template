package user

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/abdivasiyev/project_template/config"
	serviceV1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/response"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var Module = fx.Provide(NewHandler)

type Handler struct {
	environment string
	log         logger.Logger
	service     serviceV1.UserServiceV1
}

type Params struct {
	fx.In
	Config  config.Config
	Log     logger.Logger
	Service serviceV1.UserServiceV1
}

func NewHandler(params Params) *Handler {
	return &Handler{
		environment: params.Config.GetString(config.EnvironmentKey),
		log:         params.Log,
		service:     params.Service,
	}
}

// Update godoc
// @Security ApiKeyAuth
// @Summary Updates user
// @Description Returns updated user
// @Accept  json
// @Produce  json
// @Param id path string true "User id"
// @Param updateUser body models.UpdateUserRequest true "Update user request"
// @Success 200 {object} models.GetUserResponse
// @Failure default {object} models.ErrorResponse
// @Tags user
// @Router /v1/user/{id} [put]
func (h *Handler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.UpdateUserRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			h.log.Errorf("could not bind json: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not parse request body",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		request.ID = c.Param("id")

		resp, err := h.service.Update(c, request)
		if err != nil {
			h.log.Errorf("could not update user: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not update user",
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

// UpdateProfile godoc
// @Security ApiKeyAuth
// @Summary Updates user profile
// @Description Returns updated user
// @Accept  json
// @Produce  json
// @Param updateUser body models.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} models.GetUserResponse
// @Failure default {object} models.ErrorResponse
// @Tags user
// @Router /v1/user/profile [put]
func (h *Handler) UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.UpdateProfileRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			h.log.Errorf("could not bind json: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not parse request body",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		user, _ := c.Get("user")

		request.ID = (user.(models.GetUserResponse)).ID

		resp, err := h.service.UpdateProfile(c, request)
		if err != nil {
			h.log.Errorf("could not update user: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not update user",
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

// GetAll godoc
// @Security ApiKeyAuth
// @Summary Returns all users
// @Description Returns all users
// @Accept  json
// @Produce  json
// @Param filter query models.GetAllUsersRequest true "Filter params"
// @Success 200 {object} models.GetAllUsersResponse
// @Failure default {object} models.ErrorResponse
// @Tags user
// @Router /v1/user [get]
func (h *Handler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.GetAllUsersRequest

		if err := c.ShouldBindQuery(&request); err != nil {
			h.log.Errorf("could not bind query: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not bind query",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		users, err := h.service.GetAll(c, request)
		if err != nil {
			h.log.Errorf("could not get users: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get users",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    users,
			StatusCode: http.StatusOK,
		})
	}
}

// Get godoc
// @Security ApiKeyAuth
// @Summary Gets user
// @Description Returns user
// @Accept  json
// @Produce  json
// @Param id path string true "User id"
// @Success 200 {object} models.GetUserResponse
// @Failure default {object} models.ErrorResponse
// @Tags user
// @Router /v1/user/{id} [get]
func (h *Handler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, err := h.service.Get(c, id)
		if err != nil {
			h.log.Errorf("could not get user: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get user",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    user,
			StatusCode: http.StatusOK,
		})
	}
}

// GetProfile godoc
// @Security ApiKeyAuth
// @Summary Gets current user profile
// @Description Returns current user
// @Accept  json
// @Produce  json
// @Success 200 {object} models.GetUserResponse
// @Failure default {object} models.ErrorResponse
// @Tags user
// @Router /v1/user/profile [get]
func (h *Handler) GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, ok := c.Get("user")
		if !ok {
			response.JSON(c, response.Params{
				Err:        errors.New("user not authorized"),
				Message:    "user not authorized",
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		user, err := h.service.Get(c, currentUser.(models.GetUserResponse).ID)
		if err != nil {
			h.log.Errorf("could not get user: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get user",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    user,
			StatusCode: http.StatusOK,
		})
	}
}

// Delete godoc
// @Security ApiKeyAuth
// @Summary Deletes user
// @Description Deletes requested user
// @Accept  json
// @Produce  json
// @Param id path string true "User id"
// @Success 204 {object} models.SuccessResponse
// @Failure default {object} models.ErrorResponse
// @Tags user
// @Router /v1/user/{id} [delete]
func (h *Handler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := h.service.Delete(c, id)
		if err != nil {
			h.log.Errorf("could not delete user: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not delete user",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj: models.SuccessResponse{
				Ok: true,
			},
			StatusCode: http.StatusNoContent,
		})
	}
}

// Create godoc
// @Security ApiKeyAuth
// @Summary Creates new user
// @Description Returns created user
// @Accept  json
// @Produce  json
// @Param createUser body models.CreateUserRequest true "Create user request"
// @Success 201 {object} models.GetUserResponse
// @Failure default {object} models.ErrorResponse
// @Tags user
// @Router /v1/user [post]
func (h *Handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateUserRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			h.log.Errorf("could not bind json: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not parse request body",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		resp, err := h.service.Create(c, request)
		if err != nil {
			h.log.Errorf("could not create user: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not create user",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    resp,
			StatusCode: http.StatusCreated,
		})
	}
}

package role

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/abdivasiyev/project_template/config"
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
	service     serviceV1.RoleServiceV1
}

type Params struct {
	fx.In
	Config  config.Config
	Log     logger.Logger
	Service serviceV1.RoleServiceV1
}

func NewHandler(params Params) *Handler {
	return &Handler{
		environment: params.Config.GetString(config.EnvironmentKey),
		log:         params.Log,
		service:     params.Service,
	}
}

// Create godoc
// @Security ApiKeyAuth
// @Summary Creates new role
// @Description Returns created role
// @Accept  json
// @Produce  json
// @Param createForm body models.CreateRoleRequest true "Role"
// @Success 201 {object} models.GetRoleResponse
// @Failure default {object} models.ErrorResponse
// @Tags role
// @Router /v1/role [post]
func (h *Handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateRoleRequest

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
			h.log.Errorf("could not create role: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not create role",
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

// GetAll godoc
// @Security ApiKeyAuth
// @Summary Returns all roles
// @Description Returns all roles
// @Accept  json
// @Produce  json
// @Param filter query models.GetAllRoleRequest false "Filter"
// @Success 200 {object} models.GetAllRoleResponse
// @Failure default {object} models.ErrorResponse
// @Tags role
// @Router /v1/role [get]
func (h *Handler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.GetAllRoleRequest

		if err := c.ShouldBindQuery(&request); err != nil {
			h.log.Errorf("could not bind query: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not bind query",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		roles, err := h.service.GetAll(c, request)
		if err != nil {
			h.log.Errorf("could not get roles: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get roles",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    roles,
			StatusCode: http.StatusOK,
		})
	}
}

// Get godoc
// @Security ApiKeyAuth
// @Summary Returns role
// @Description Returns all role modules
// @Accept  json
// @Produce  json
// @Param id path string true "Role id"
// @Success 200 {object} models.GetRoleResponse
// @Failure default {object} models.ErrorResponse
// @Tags role
// @Router /v1/role/{id} [get]
func (h *Handler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		role, err := h.service.Get(c, id)
		if err != nil {
			h.log.Errorf("could not get role: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get role",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    role,
			StatusCode: http.StatusOK,
		})
	}
}

// GetModules godoc
// @Security ApiKeyAuth
// @Summary Returns all role modules
// @Description Returns all role modules
// @Accept  json
// @Produce  json
// @Success 200 {object} models.GetModulesResponse
// @Failure default {object} models.ErrorResponse
// @Tags role
// @Router /v1/role/modules [get]
func (h *Handler) GetModules() gin.HandlerFunc {
	return func(c *gin.Context) {
		modules, err := h.service.GetModules(c)
		if err != nil {
			h.log.Errorf("could not get role modules: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get role modules",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    modules,
			StatusCode: http.StatusOK,
		})
	}
}

// Update godoc
// @Security ApiKeyAuth
// @Summary Updates new role
// @Description Returns created role
// @Accept  json
// @Produce  json
// @Param id path string true "Role id"
// @Param createForm body models.CreateRoleRequest true "Role"
// @Success 200 {object} models.GetRoleResponse
// @Failure default {object} models.ErrorResponse
// @Tags role
// @Router /v1/role/{id} [put]
func (h *Handler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.CreateRoleRequest

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
			h.log.Errorf("could not create role: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not create role",
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

// Delete godoc
// @Security ApiKeyAuth
// @Summary Deletes role
// @Description Returns created role
// @Accept  json
// @Produce  json
// @Param id path string true "Role id"
// @Success 204
// @Failure default {object} models.ErrorResponse
// @Tags role
// @Router /v1/role/{id} [delete]
func (h *Handler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h.service.Delete(c, c.Param("id"))
		if err != nil {
			h.log.Errorf("could not delete role: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not create role",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    nil,
			StatusCode: http.StatusNoContent,
		})
	}
}

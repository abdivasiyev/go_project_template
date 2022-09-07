// Package v1 handler
// @title Logistics application API
// @version 1.0
// @description API Gateway.
// @query.collection.format multi
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package v1

import (
	"context"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/handler/v1/app"
	"github.com/abdivasiyev/project_template/internal/handler/v1/auth"
	"github.com/abdivasiyev/project_template/internal/handler/v1/doc"
	"github.com/abdivasiyev/project_template/internal/handler/v1/file"
	"github.com/abdivasiyev/project_template/internal/handler/v1/pprof"
	"github.com/abdivasiyev/project_template/internal/handler/v1/role"
	"github.com/abdivasiyev/project_template/internal/handler/v1/user"
	"github.com/abdivasiyev/project_template/internal/middleware"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Invoke(New),
)

type Params struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Config     config.Config
	Router     *gin.Engine
	Logger     logger.Logger
	Security   security.Handler
	Middleware middleware.Handler
	Auth       *auth.Handler
	File       *file.Handler
	Role       *role.Handler
	User       *user.Handler
	Doc        *doc.Handler
	App        *app.Handler
	Pprof      *pprof.Handler
}

type Handler struct {
	logger            logger.Logger
	security          security.Handler
	middleware        middleware.Handler
	auth              *auth.Handler
	file              *file.Handler
	role              *role.Handler
	user              *user.Handler
	doc               *doc.Handler
	pprof             *pprof.Handler
	app               *app.Handler
	basicAuthUser     string
	basicAuthPassword string
	swaggerPath       string
}

func New(params Params) {
	handler := &Handler{
		auth:              params.Auth,
		file:              params.File,
		role:              params.Role,
		user:              params.User,
		pprof:             params.Pprof,
		basicAuthUser:     params.Config.GetString(config.BasicAuthUserKey),
		basicAuthPassword: params.Config.GetString(config.BasicAuthPasswordKey),
		logger:            params.Logger,
		security:          params.Security,
		middleware:        params.Middleware,
		doc:               params.Doc,
		swaggerPath:       params.Config.GetString(config.SpecPath),
		app:               params.App,
	}

	params.Lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				params.Logger.Info("Handlers registering")
				handler.registerRoutes(params.Router)
				return nil
			},
		},
	)
}

func (h *Handler) registerRoutes(router gin.IRouter) {
	apiV1 := router.Group("/v1")

	basicAuth := apiV1.Group("/", gin.BasicAuth(gin.Accounts{
		h.basicAuthUser: h.basicAuthPassword,
	}))

	authRequired := apiV1.Group("/")
	authRequired.Use(
		h.middleware.BearerAuth(),
		h.middleware.HasAccess(),
	)

	// basic auth
	h.registerDoc(basicAuth)
	// no auth
	h.registerAuth(apiV1)
	h.registerApp(apiV1)
	// auth required
	h.registerUser(authRequired)
	h.registerRole(authRequired)
	h.registerFile(authRequired)
	h.registerPprof(apiV1)
}

func (h *Handler) registerApp(group gin.IRouter) {
	routerGroup := group.Group("/app")
	{
		routerGroup.GET("/version", h.app.GetVersion())
	}
}

func (h *Handler) registerUser(group gin.IRouter) {
	routerGroup := group.Group("/user")
	{
		routerGroup.POST("/", h.user.Create())
		routerGroup.PUT("/:id", h.user.Update())
		routerGroup.DELETE("/:id", h.user.Delete())
		routerGroup.GET("/", h.user.GetAll())
		routerGroup.GET("/:id", h.user.Get())
		routerGroup.GET("/profile", h.user.GetProfile())
		routerGroup.PUT("/profile", h.user.UpdateProfile())
	}
}

func (h *Handler) registerRole(group gin.IRouter) {
	routerGroup := group.Group("/role")
	{
		routerGroup.POST("/", h.role.Create())
		routerGroup.PUT("/:id", h.role.Update())
		routerGroup.DELETE("/:id", h.role.Delete())
		routerGroup.GET("/:id", h.role.Get())
		routerGroup.GET("/", h.role.GetAll())
		routerGroup.GET("/modules", h.role.GetModules())
	}
}

func (h *Handler) registerFile(group gin.IRouter) {
	routerGroup := group.Group("/file")
	{
		routerGroup.POST("/", h.file.Upload())
		routerGroup.GET("/:id", h.file.Get())
	}
}

func (h *Handler) registerDoc(group gin.IRouter) {
	routerGroup := group.Group("/docs")
	{
		routerGroup.GET("/", h.doc.Render())
		routerGroup.StaticFile("swagger.json", h.swaggerPath)
	}
}

// Auth handlers
func (h *Handler) registerAuth(group gin.IRouter) {
	routerGroup := group.Group("/auth")
	{
		routerGroup.POST("/login", h.auth.Login())
		routerGroup.POST("/refresh", h.auth.Refresh())
		routerGroup.POST("/reset-password", h.auth.ResetPassword())
	}
}

func (h *Handler) registerPprof(group gin.IRouter) {
	routerGroup := group.Group("/debug")
	{
		routerGroup.Any("/cmdline", h.pprof.Cmdline())
		routerGroup.Any("/profile", h.pprof.Profile())
		routerGroup.Any("/heap", h.pprof.Heap())
		routerGroup.Any("/symbol", h.pprof.Symbol())
	}
}

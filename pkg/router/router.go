package router

import (
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/middleware"
	"github.com/abdivasiyev/project_template/pkg/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In
	Config     config.Config
	Middleware middleware.Handler
}

// New constructor of router
func New(params Params) *gin.Engine {
	router := gin.New()

	if params.Config.GetString(config.EnvironmentKey) == config.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup middlewares
	router.Use(params.Middleware.Log(config.DateTimeFormat, true))
	router.Use(params.Middleware.RecoverWithLog(true))
	router.Use(cors.Default())

	router.LoadHTMLGlob(config.RootDir() + "/templates/*.html")
	// register validator
	binding.Validator = new(validator.TranslatableValidator)

	return router
}

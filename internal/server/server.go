package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

var Module = fx.Options(fx.Invoke(New))

type Params struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    config.Config
	Handler   *gin.Engine
	Logger    logger.Logger
}

func New(params Params) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", params.Config.GetInt(config.HttpPortKey)),
		Handler: params.Handler,
	}

	params.Lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				params.Logger.Info("Application started")
				go func() {
					if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						params.Logger.Error("Could not start server", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				params.Logger.Info("Application stopped")
				return server.Shutdown(ctx)
			},
		},
	)
}

package sentry

import (
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Handler interface {
	HandleError(err error)
	HandleMessage(format string, args ...any)
}

type Params struct {
	fx.In
	Config config.Config
}

type sentryHandler struct {
	hub         *sentry.Hub
	environment string
}

func New(params Params) (Handler, error) {
	debug := false

	if params.Config.GetString(config.EnvironmentKey) == config.Development {
		debug = true
	}

	sentryClient, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              params.Config.GetString(config.SentryDSNKey),
		Environment:      params.Config.GetString(config.EnvironmentKey),
		Debug:            debug,
		AttachStacktrace: true,
		TracesSampleRate: 1,
		ServerName:       params.Config.GetString(config.NamespaceKey),
	})

	if err != nil {
		return &sentryHandler{}, errors.Wrap(err, "could not connect to client")
	}

	hub := sentry.CurrentHub()

	hub.BindClient(sentryClient)

	return &sentryHandler{
		hub:         hub,
		environment: params.Config.GetString(config.EnvironmentKey),
	}, nil
}

func (s *sentryHandler) isProd() bool {
	return s.environment == config.Production
}

func (s *sentryHandler) HandleMessage(format string, args ...any) {
	if !s.isProd() {
		return
	}
	message := fmt.Sprintf(format, args...)
	s.hub.CaptureMessage(message)
}

func (s *sentryHandler) HandleError(err error) {
	if !s.isProd() {
		return
	}
	s.hub.CaptureException(err)
}

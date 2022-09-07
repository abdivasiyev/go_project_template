package middleware

import (
	"net/http"
	"time"

	"github.com/abdivasiyev/project_template/internal/models"

	"github.com/gin-gonic/gin"
)

// Config is config setting for Ginzap
type Config struct {
	TimeFormat string
	UTC        bool
	SkipPaths  []string
}

// Log returns a gin.HandlerFunc (middleware) that logs requests using custom logger
func (m *middleware) Log(timeFormat string, utc bool) gin.HandlerFunc {
	return m.LogWithConfig(&Config{TimeFormat: timeFormat, UTC: utc})
}

// LogWithConfig returns a gin.HandlerFunc using configs
func (m *middleware) LogWithConfig(conf *Config) gin.HandlerFunc {
	skipPaths := make(map[string]struct{}, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = struct{}{}
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		request := c.Request.Clone(c)

		c.Next()

		if _, ok := skipPaths[request.URL.Path]; !ok {
			end := time.Now()
			if conf.UTC {
				end = end.UTC()
			}

			errs := make([]error, len(c.Errors))

			for _, err := range c.Errors {
				errs = append(errs, err.Err)
			}

			m.service.Log(c, c.Writer.Status(), request, c.ClientIP(), start, end, errs, conf.TimeFormat)
		}
	}
}

// RecoverWithLog returns a gin.HandlerFunc (middleware)
func (m *middleware) RecoverWithLog(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				m.service.RecoverPanic(c, c.Request, err.(error), stack)

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					models.ErrorResponse{
						ErrorCode:    http.StatusInternalServerError,
						ErrorMessage: http.StatusText(http.StatusInternalServerError),
					},
				)
			}
		}()
		c.Next()
	}
}

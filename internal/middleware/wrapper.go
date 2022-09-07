package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *middleware) WrapHttpHandler(f http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		f.ServeHTTP(c.Writer, c.Request)
	}
}

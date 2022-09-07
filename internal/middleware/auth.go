package middleware

import (
	"errors"
	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/pkg/helpers"
	"github.com/abdivasiyev/project_template/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *middleware) HasAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := c.Get("user")

		if !ok {
			response.JSON(c, response.Params{
				Err:        models.ErrUnauthorized,
				StatusCode: http.StatusUnauthorized,
				Message:    "unauthorized",
			})
			return
		}

		err := m.service.HasAccess(c, user.(models.GetUserResponse).ID, c.FullPath(), c.Request.Method, func(queryParam string) string {
			result := c.Param(queryParam)

			if helpers.IsEmpty(result) {
				result = c.Query(queryParam)
			}

			return result
		})

		if err != nil {
			if errors.Is(err, models.ErrForbidden) {
				response.JSON(c, response.Params{
					Err:        models.ErrForbidden,
					StatusCode: http.StatusForbidden,
					Message:    "forbidden operation",
				})
				return
			}
			response.JSON(c, response.Params{
				Err:        err,
				StatusCode: http.StatusInternalServerError,
				Message:    "internal error",
			})
			return
		}
	}
}

func (m *middleware) BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")

		user, err := m.service.CheckAuth(c, accessToken)
		if err != nil {
			if errors.Is(err, models.ErrUnauthorized) {
				response.JSON(c, response.Params{
					Err:        err,
					StatusCode: http.StatusUnauthorized,
					Message:    "unauthorized",
				})
				return
			}
			response.JSON(c, response.Params{
				Err:        err,
				StatusCode: http.StatusInternalServerError,
				Message:    "internal error",
			})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

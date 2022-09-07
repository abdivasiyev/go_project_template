package response

import (
	"errors"
	"github.com/abdivasiyev/project_template/pkg/helpers"
	"github.com/abdivasiyev/project_template/pkg/security/jwt"
	"io"
	"net/http"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/gin-gonic/gin"
)

// Params arguments for returning response
type Params struct {
	JsonObj    interface{}
	Err        error
	StatusCode int
	Message    string
}

// JSON returns json response to gin context
func JSON(c *gin.Context, params Params) {
	if jsonWithError(c, params) {
		return
	}

	// handle another errors
	if params.StatusCode >= http.StatusBadRequest {
		c.AbortWithStatusJSON(params.StatusCode, models.ErrorResponse{
			ErrorCode:    params.StatusCode,
			ErrorMessage: params.Message,
		})
		return
	}

	c.AbortWithStatusJSON(params.StatusCode, params.JsonObj)
}

func jsonWithError(c *gin.Context, params Params) bool {
	// handle error and check validation errors
	if params.Err == nil {
		return false
	}

	switchedErr := false
	resp := models.ErrorResponse{}

	switch {
	case errors.Is(params.Err, models.ErrNotFound):
		resp.ErrorCode = http.StatusNotFound
		resp.ErrorMessage = http.StatusText(http.StatusNotFound)
		switchedErr = !switchedErr
	case errors.Is(params.Err, jwt.ErrInvalidToken):
		resp.ErrorCode = http.StatusBadRequest
		resp.ErrorMessage = jwt.ErrInvalidToken.Error()
		switchedErr = !switchedErr
	case errors.Is(params.Err, jwt.ErrExpiredToken):
		resp.ErrorCode = http.StatusBadRequest
		resp.ErrorMessage = jwt.ErrExpiredToken.Error()
		switchedErr = !switchedErr
	case errors.Is(params.Err, io.EOF) || errors.Is(params.Err, io.ErrUnexpectedEOF):
		resp.ErrorCode = http.StatusBadRequest
		resp.ErrorMessage = http.StatusText(http.StatusBadRequest)
		switchedErr = !switchedErr
	case errors.Is(params.Err, models.ErrForbidden):
		resp.ErrorCode = http.StatusForbidden
		resp.ErrorMessage = "you are not allowed to perform this action"
		switchedErr = !switchedErr
	}

	if switchedErr {
		c.AbortWithStatusJSON(resp.ErrorCode, resp)
		return true
	}

	resp = helpers.ConvertErrorToErrorResponse(params.StatusCode, params.Err)
	resp.ErrorMessage = params.Message

	c.AbortWithStatusJSON(resp.ErrorCode, resp)
	return true
}

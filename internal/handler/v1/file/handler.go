package file

import (
	"crypto/sha512"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/fx"

	"github.com/abdivasiyev/project_template/config"
	serviceV1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/response"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewHandler)

type Handler struct {
	environment string
	uploadPath  string
	log         logger.Logger
	service     serviceV1.FileServiceV1
}

type Params struct {
	fx.In
	Config  config.Config
	Log     logger.Logger
	Service serviceV1.FileServiceV1
}

func NewHandler(params Params) *Handler {
	return &Handler{
		environment: params.Config.GetString(config.EnvironmentKey),
		uploadPath:  params.Config.GetString(config.UploadPathKey),
		log:         params.Log,
		service:     params.Service,
	}
}

// Upload godoc
// @Summary Uploads file to minio service
// @Description Returns file url and file id
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "File body"
// @Success 200 {object} models.GetFileResponse
// @Failure default {object} models.ErrorResponse
// @Tags file
// @Router /v1/file [post]
func (h *Handler) Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			h.log.Errorf("could not get file: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get file",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		h.log.Debug("multipart file", zap.Any("header", fileHeader))

		uploadedResponse, err := h.service.UploadFile(c, fileHeader)
		if err != nil {
			h.log.Errorf("could not save file: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not save file",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		response.JSON(c, response.Params{
			JsonObj:    uploadedResponse,
			StatusCode: http.StatusOK,
		})
	}
}

// Get godoc
// @Summary Gets uploaded file by id
// @Description Returns file url and file id
// @Accept  json
// @Produce  json
// @Param id path string true "File name"
// @Success 200 {string} string "file in bytes"
// @Failure default {object} models.ErrorResponse
// @Tags file
// @Router /v1/file/{id} [get]
func (h *Handler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.GetFileRequest

		if err := c.ShouldBindUri(&request); err != nil {
			h.log.Errorf("could not bind uri params: %v", err)

			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not bind uri params",
				StatusCode: http.StatusBadRequest,
			})
			return
		}

		cacheData := []byte(fmt.Sprint(time.Now().UTC().UnixNano()))
		etag := fmt.Sprintf("%x", sha512.Sum512(cacheData))

		if match := c.GetHeader("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				c.Status(http.StatusNotModified)
				return
			}
		}

		uploadResponse, filePath, err := h.service.GetFile(c, request.ID)

		if err != nil {
			h.log.Errorf("could not get file: %v", err)
			response.JSON(c, response.Params{
				Err:        err,
				Message:    "could not get file",
				StatusCode: http.StatusInternalServerError,
			})
			return
		}

		c.Header("Cache-Control", "public, max-age=3600")
		c.Header("ETag", etag)
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", uploadResponse.FileName))

		c.File(filePath)
	}
}

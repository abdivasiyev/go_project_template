package file_service

import (
	"context"
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"github.com/abdivasiyev/project_template/pkg/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/internal/repository"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var Module = fx.Provide(NewService)

type service struct {
	environment    string
	uploadPath     string
	cdnURL         string
	log            logger.Logger
	sentry         sentry.Handler
	fileRepository repository.File
	cache          storage.Cacher
}

type Params struct {
	fx.In
	Config         config.Config
	Log            logger.Logger
	Sentry         sentry.Handler
	FileRepository repository.File
	Cache          storage.Cacher
}

func NewService(params Params) v1.FileServiceV1 {
	return &service{
		environment:    params.Config.GetString(config.EnvironmentKey),
		log:            params.Log,
		sentry:         params.Sentry,
		fileRepository: params.FileRepository,
		uploadPath:     params.Config.GetString(config.UploadPathKey),
		cdnURL:         params.Config.GetString(config.CdnURLKey),
		cache:          params.Cache,
	}
}

func (s *service) UploadFile(ctx context.Context, multipartFileHeader *multipart.FileHeader) (models.GetFileResponse, error) {
	var (
		fileID   = uuid.New().String()
		fileName = fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(multipartFileHeader.Filename))
		filePath = fmt.Sprintf("%s/%s", s.uploadPath, fileName)
		fileURL  = fmt.Sprintf("%s/%s", s.cdnURL, fileName)
	)

	if err := s.upload(multipartFileHeader, filePath); err != nil {
		s.sentry.HandleError(err)
		s.log.Errorf("could not save multipartFileHeader: %v", err)
		return models.GetFileResponse{}, errors.Wrap(err, "could not save file")
	}

	fileResponse := models.GetFileResponse{
		FileID:   fileID,
		FileName: fileName,
		FileURL:  fileURL,
	}

	err := s.fileRepository.Create(ctx, fileResponse)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not create file", zap.Error(err), zap.Any("fileResponse", fileResponse))
	}

	return fileResponse, errors.Wrap(err, "could not save file")
}

func (s *service) upload(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not open file", zap.Error(err))
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not create file", zap.Error(err))
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not copy file", zap.Error(err))
	}
	return err
}

func (s *service) GetFile(ctx context.Context, id string) (models.GetFileResponse, string, error) {
	var (
		resp models.GetFileResponse
		key  = "file:storage:" + id
	)

	err := s.cache.GetObj(ctx, key, &resp)
	if err == nil {
		return s.getFilePath(resp)
	}

	resp, err = s.fileRepository.Get(ctx, id)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			s.sentry.HandleError(err)
			s.log.Error("could not get file", zap.Error(err), zap.Any("fileID", id))
		}
		return models.GetFileResponse{}, "", errors.Wrap(err, "could not get uploadResponse")
	}

	err = s.cache.SetObj(ctx, key, resp, 12*time.Hour)
	if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not cache file", zap.Error(err), zap.Any("fileID", id))
	}

	return s.getFilePath(resp)
}

func (s *service) getFilePath(resp models.GetFileResponse) (models.GetFileResponse, string, error) {
	filePath := fmt.Sprintf("%s/%s", s.uploadPath, resp.FileName)

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		s.sentry.HandleError(err)
		s.log.Warn("osStat file not exist")
		return resp, "", models.ErrNotFound
	} else if err != nil {
		s.sentry.HandleError(err)
		s.log.Error("could not get stat file", zap.Error(err))
		return resp, "", errors.Wrap(err, "could not get file stat")
	}

	if fileInfo.IsDir() {
		s.log.Warn("path is a directory")
		return resp, "", models.ErrNotFound
	}

	return resp, filePath, nil
}

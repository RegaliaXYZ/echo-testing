package v1

import (
	"bp-echo-test/internal/database"
	"bp-echo-test/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type V1Handler struct {
	database database.Service
	logger   *zap.Logger
	gcs      utils.GoogleServiceInterface
}

type V1Service interface {
	Health(c echo.Context) error
	Create(c echo.Context) error
	UploadContent(c echo.Context) error

	DeleteModels(c echo.Context) error
	ListModels(c echo.Context) error
}

func NewV1Handler(database database.Service, logger *zap.Logger, gcs utils.GoogleServiceInterface) V1Service {
	return &V1Handler{
		database: database,
		logger:   logger,
		gcs:      gcs,
	}
}

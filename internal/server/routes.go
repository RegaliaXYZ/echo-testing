package server

import (
	v1 "bp-echo-test/internal/server/v1"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func (s *Server) RegisterRoutes(auth_token string) http.Handler {
	e := echo.New()
	//e.Use(middleware.Logger())
	logger, _ := zap.NewDevelopment()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:   true,
		LogMethod:    true,
		LogStatus:    true,
		LogRoutePath: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.Int64("latency", v.Latency.Milliseconds()),
				zap.String("method", v.Method),
				zap.Int("status", v.Status),
				zap.String("path", v.RoutePath),
			}
			if v.Status == http.StatusOK {
				logger.Info("", fields...)
			} else {
				logger.Error("", fields...)
			}
			return nil
		},
	}))

	e.Use(middleware.Recover())

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "Basic",
		Validator:  func(auth string, c echo.Context) (bool, error) { return auth == auth_token, nil },
	}))

	e.GET("/", s.HelloWorldHandler)
	e.GET("/health", s.healthHandler)

	v1_handler := v1.NewV1Handler(s.db, logger, s.gcp_client)
	v1 := e.Group("/api/v1/ml/tenants/:tenant/type/:model_type")
	{

		v1.Use(validateModelType)
		v1.DELETE("/models", v1_handler.DeleteModels)
		v1.POST("/models", v1_handler.Create)
		v1.GET("/models", v1_handler.ListModels)
		models := v1.Group("/models/:model")
		{
			models.POST("/populate", v1_handler.UploadContent)
		}
	}

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func validateModelType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		modelType := strings.ToLower(c.Param("model_type"))

		// Check if the model_type is either "nlu" or "ner"
		if modelType != "nlu" && modelType != "ner" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid model type. It must be either 'nlu' or 'ner'.")
		}

		// Continue to the next middleware or route handler
		return next(c)
	}
}

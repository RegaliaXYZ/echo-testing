package v1

import (
	"bp-echo-test/internal/database"
	"bp-echo-test/internal/models"
	"bp-echo-test/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (v1 *V1Handler) Create(c echo.Context) error {
	tenant := strings.ToLower(c.Param("tenant"))
	model_type := strings.ToLower(c.Param("model_type"))
	input := new(models.CreateModelInput)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot bind input")
	}
	// v1.logger.Info("input", zap.Any("input", input), zap.String("tenant", tenant), zap.String("model_type", model_type))
	model, err := v1.database.GetByName(tenant, model_type, input.Name)
	if err != nil && err != database.ErrNotFound {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot get model")
	}
	if model != (database.Model{}) {
		return echo.NewHTTPError(http.StatusBadRequest, "model already exists")
	}

	id, err := v1.database.Create(tenant, database.Model{
		ModelType:     model_type,
		Name:          input.Name,
		Service:       input.Service,
		SubService:    input.SubService,
		Language:      input.Language,
		Status:        "initiated",
		NumUtterances: 0,
		CreatedAt:     time.Now(),
	})
	//id, err := database.PGSQLDatabase.InsertModel(tenant, task, userInput.Model, userInput)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "cannot create model")
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "model created successfully",
		"id":     id.String(),
	})
}

func (v1 *V1Handler) UploadContent(c echo.Context) error {
	tenant := strings.ToLower(c.Param("tenant"))
	model_type := strings.ToLower(c.Param("model_type"))
	model_name := strings.ToLower(c.Param("model"))

	decoded, err := utils.GetEncodedData(c.Request().Body, c.Request().Header.Get("Content-Encoding"))
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot decode body")
	}
	var input interface{}

	status := "populated"
	if model_type == "nlu" {
		input = new(models.UploadContentNLUInput)
		if err := json.Unmarshal(decoded, &input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cannot bind input")
		}
		if strings.ToLower(input.(*models.UploadContentNLUInput).Complete) != "yes" {
			status = "populating"
		}
	} else if model_type == "ner" {
		input = new(models.UploadContentNERInput)
		if err := json.Unmarshal(decoded, &input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cannot bind input")
		}
		if strings.ToLower(input.(*models.UploadContentNERInput).Complete) != "yes" {
			status = "populating"
		}
	}
	fmt.Println(input)
	// check if input has complete field

	v1.logger.Info(status)
	model, err := v1.database.GetByName(tenant, model_type, model_name)
	if err != nil && err != database.ErrNotFound {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot get model")
	}
	if model == (database.Model{}) {
		return echo.NewHTTPError(http.StatusNotFound, "model does not exist")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":     "ok",
		"tenant":     tenant,
		"type":       model_type,
		"model_name": model_name,
		"model":      model,
		"input":      input,
	})
}

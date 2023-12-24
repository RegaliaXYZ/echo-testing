package v1

import (
	"bp-echo-test/internal/database"
	"bp-echo-test/internal/models"
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
	model, err := v1.database.GetByName(tenant, model_type, input.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get model")
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
	model := strings.ToLower(c.Param("model"))

	// get request body

	requestBody := c.Request().Body

	// Your handling logic here

	// Example: Reading the request body
	// You may want to use a specific data structure or parse the body as needed.
	// Here, we're just printing the body to the console.
	buf := make([]byte, 0)
	_, err := requestBody.Read(buf)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read request body"})
	}
	fmt.Println(string(buf))

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"tenant": tenant,
		"type":   model_type,
		"model":  model,
	})
}

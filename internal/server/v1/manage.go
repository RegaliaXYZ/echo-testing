package v1

import (
	"bp-echo-test/internal/database"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (v1 *V1Handler) ListModels(c echo.Context) error {
	tenant := c.Param("tenant")
	model_type := strings.ToLower(c.Param("model_type"))
	//models, err := database.GetModelsFromDB(tenant, task, db)
	models, err := v1.database.GetByModelType(tenant, model_type)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot get models")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"value":  "OK",
		"models": models,
	})
}

func (v1 *V1Handler) DeleteModels(c echo.Context) error {
	var userInput *struct {
		Models []string `json:"models" binding:"required"`
	}
	tenant := c.Param("tenant")
	task := strings.ToLower(c.Param("model_type"))

	if err := c.Bind(userInput); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot bind input")
	}

	var deletedModels = make([]string, 0)
	var notDeletedModels = make([]string, 0)
	v1.logger.Debug("DELETE_MODELS", zap.Any("models_to_delete", userInput.Models))
	for _, model_name := range userInput.Models {
		v1.logger.Debug("model_name", zap.String("model_name", model_name))
		model, err := v1.database.GetByName(tenant, task, model_name)
		if err != nil {
			notDeletedModels = append(notDeletedModels, model_name)
			v1.logger.Debug("error fetching the model", zap.Error(err))
			continue
		}
		v1.logger.Debug("model", zap.Any("model", model))
		if model == (database.Model{}) {
			notDeletedModels = append(notDeletedModels, model_name)
			continue
		}
		if model.Status == "deleted" {
			deletedModels = append(deletedModels, model_name)
			continue
		}

		if model.Status == "training" {
			deletedModels = append(deletedModels, model_name)
			continue
		}
		//deployments, err := database.GetDeploymentsFromDB(model.ID, "", db)

		updatedModel := database.Model{
			ID:            model.ID,
			ModelType:     model.ModelType,
			Name:          model.Name,
			Service:       model.Service,
			SubService:    model.SubService,
			Language:      model.Language,
			Status:        "deleted",
			JobID:         model.JobID,
			RunID:         model.RunID,
			NumUtterances: model.NumUtterances,
			CreatedAt:     model.CreatedAt,
		}
		// WE CAN UNDEPLOY OTHERWISE
		//v1.gke.undeployModels(c, a.config.GCP.Namespace, tenant, model.Name)

		// deleting model from GCP Bucket
		path := "tenants/" + tenant + "/" + model.Name + "/"
		err = v1.gcs.DeleteFolder(path)
		if err != nil {
			notDeletedModels = append(notDeletedModels, model_name)
			continue
		}
		model.Status = "deleted"
		err = v1.database.Update(tenant, updatedModel)
		if err != nil {
			notDeletedModels = append(notDeletedModels, model_name)
			continue
		}
		deletedModels = append(deletedModels, model_name)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"deletedModels":    deletedModels,
		"notDeletedModels": notDeletedModels,
	})
}

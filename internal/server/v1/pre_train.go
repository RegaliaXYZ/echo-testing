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
	"go.uber.org/zap"
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
	input := new(models.GlobalUploadContentInput)

	status := "populated"
	if err := json.Unmarshal(decoded, &input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot bind input")
	}
	if strings.ToLower(input.Complete) != "yes" {
		status = "populating"
	}
	/*
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
		}*/
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

	model_folder := "tenants/" + tenant + "/" + model_type + "/"
	err = v1.gcs.WriteFolder(model_folder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not write folder")
	}

	json_full_filename := model_folder + model_name + "/json_corpus/" + model_name + "_full.json"
	if status == "populated" {
		defer func() error {
			// delete the json file
			err := v1.gcs.DeleteFile(json_full_filename)
			if err != nil {
				v1.logger.Error("could not delete file", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "could not delete json full file")
			}
			return nil
		}()
	}

	var content []models.Payload
	numUtterances := 0
	fileExists, err := v1.gcs.FileExists(json_full_filename)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not check if file exists")
	}
	if fileExists {
		// read the file
		err = v1.gcs.ReadFile(json_full_filename, &content)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not read & parse json full filename")
		}
		content = append(content, input.Payload...)
		numUtterances = len(content)
	} else {
		numUtterances = len(input.Payload)
		content = input.Payload
	}
	err = v1.gcs.WriteFile(json_full_filename, content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not write json full filename")
	}

	corpus_file := model_folder + model_name + "/corpus/" + model_name + ".csv"
	w := v1.gcs.GetWriter(corpus_file)
	err = utils.CorpusToCSV(model_type, content, "train", w)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not write convert or write csv train")
	}
	w.Close()
	corpus_test_file := model_folder + model_name + "/test_corpus/" + model_name + ".csv"
	w = v1.gcs.GetWriter(corpus_test_file)
	err = utils.CorpusToCSV(model_type, content, "test", w)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not write convert or write csv test")
	}
	w.Close()
	updatedModel := database.Model{
		ID:            model.ID,
		ModelType:     model_type,
		Name:          model_name,
		Service:       model.Service,
		SubService:    model.SubService,
		Language:      model.Language,
		Status:        status,
		JobID:         model.JobID,
		RunID:         model.RunID,
		NumUtterances: numUtterances,
		CreatedAt:     model.CreatedAt,
	}

	err = v1.database.Update(tenant, updatedModel)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not update model with id: "+model.ID.String())
	}
	return c.JSON(http.StatusOK, map[string]any{
		"value":     "OK",
		"modelInfo": updatedModel,
	})
	/*
		return c.JSON(http.StatusOK, map[string]any{
			"status":         "ok",
			"tenant":         tenant,
			"type":           model_type,
			"model_name":     model_name,
			"model":          model,
			"input":          input,
			"num_utterances": numUtterances,
		})
	*/
}

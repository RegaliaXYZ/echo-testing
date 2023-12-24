package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (v1 *V1Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}

package setting_view

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func LevelOfPrivacy() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "Not implemented yet."})
	}
}

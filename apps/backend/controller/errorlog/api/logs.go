package errorlog_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/database/sqlc"
	"resuming/service"
)

func GetErrorAuditLogs() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key("error_log_data").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to retrieve error logs."})
		}

		var data []sqlc.ErrorLog
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to parse error logs."})
		}

		c.Set("response_data", data)
		return nil
	}
}

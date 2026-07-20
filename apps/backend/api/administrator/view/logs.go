package administrator_view

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/database/sqlc"
	"resuming/tool"
)

func GetClientAuditLogs() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_audit_log_data").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to retrieve clients audit logs."})
		}

		var data []sqlc.ClientAuditLog
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to parse client audit logs."})
		}

		c.Set("response_data", data)
		return nil
	}
}

func GetAdminAuditLogs() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_audit_log_data").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to retrieve admin audit logs."})
		}

		var data []sqlc.AdminAuditLog
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to parse admin audit logs."})
		}

		c.Set("response_data", data)
		return nil
	}
}

func GetErrorAuditLogs() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("error_log_data").Build()).ToString()
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

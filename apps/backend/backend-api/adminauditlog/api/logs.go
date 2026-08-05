package adminauditlog_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/database/sqlc"
	"resuming/tool"
)

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

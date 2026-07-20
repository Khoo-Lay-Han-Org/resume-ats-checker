package administrator_view

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/tool"
)

func GetClients() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_configs").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to retrieve clients config data."})
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to parse clients config data."})
		}

		c.Set("response_data", data)
		return nil
	}
}

func GetAdmins() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_configs").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to retrieve admins config data."})
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to parse admins config data."})
		}

		c.Set("response_data", data)
		return nil
	}
}

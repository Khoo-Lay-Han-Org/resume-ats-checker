package session_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	dto "resuming/controller/session/dto"
	validator "resuming/controller/session/validator"
	"resuming/service"
)

func RemoveIndividualUserSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request dto.SessionControlRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve request."})
		}

		polished_request, err := validator.ValidateSessionControlRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		public_user_id := polished_request.PublicUserId

		admin_session_id := c.Get("public_user_id")
		if admin_session_id != nil && public_user_id == admin_session_id.(string) {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Cannot remove your own session."})
		}

		ctx := c.Request().Context()
		exists, err := service.Valkey.Do(
			ctx,
			service.Valkey.B().Exists().Key(public_user_id+":session_data").Build(),
		).AsInt64()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to verify session."})
		}
		if exists == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Session not found."})
		}

		err = service.Valkey.Do(
			ctx,
			service.Valkey.B().Del().
				Key(public_user_id+":session_data",
					public_user_id+":jwt_data",
					public_user_id+":user_data").
				Build(),
		).Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete session"})
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "Successfully deleted session"})
	}
}

func RemoveAllClientSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key("client_configs").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to find client configs."})
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse configs."})
		}

		for _, item := range data {
			public_user_id, ok := item["public_user_id"].(string)
			if !ok {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Invalid client config format."})
			}

			ctx := c.Request().Context()
			_, err := service.Valkey.Do(
				ctx,
				service.Valkey.B().Del().
					Key(public_user_id+":session_data").
					Build(),
			).ToString()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete session"})
			}
		}

		admin_id := c.Get("public_user_id")
		if admin_id != nil {
			admin_key := admin_id.(string) + ":session_data"
			ttl_seconds, err := service.Valkey.Do(ctx, service.Valkey.B().Ttl().Key(admin_key).Build()).AsInt64()
			if err == nil && ttl_seconds > 300 {
				service.Valkey.Do(ctx, service.Valkey.B().Expire().Key(admin_key).Seconds(300).Build())
			}
		}

		return nil
	}
}

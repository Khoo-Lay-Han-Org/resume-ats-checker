package client_support_view

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	typing "resuming/api/client-support/typing"
	validator "resuming/api/client-support/validator"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func ClientCommunicateToAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get session data."})
		}
		public_user_id := retrieved_public_user_id.(string)

		var request typing.ClientCommunicateRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to get request."})
		}

		polished_request, err := validator.ValidateClientCommunicateRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		type_of_message := polished_request.Type
		client_message := polished_request.Message

		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_support_messages").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get support message data."})
		}

		var data []map[string]any
		if err := json.Unmarshal([]byte(retrieved_data), &data); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse data."})
		}

		new_client_comm := map[string]any{
			"public_id":                 uuid.New().String(),
			"user_id":                   public_user_id,
			"type":                      type_of_message,
			"message":                   client_message,
			"sender_type":               "client",
			"client_comm_log_public_id": "",
			"created_at":                time.Now(),
		}

		data = append(data, new_client_comm)

		serialised_data, err := json.Marshal(data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to serialise support message data."})
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("client_support_messages").Value(string(serialised_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "failed to store support message data."})
		}

		return nil
	}
}

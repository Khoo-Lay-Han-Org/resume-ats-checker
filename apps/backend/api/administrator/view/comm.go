package administrator_view

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	valkey "github.com/valkey-io/valkey-go"
	typing "resuming/api/administrator/typing"
	validator "resuming/api/administrator/validator"
	"resuming/database"
	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func GetSupportMessages() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_support_messages").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve support messages."})
		}

		var data []map[string]any
		if err := json.Unmarshal([]byte(retrieved_data), &data); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse support messages."})
		}

		var polished_data_arr []map[string]any
		for _, item := range data {
			public_user_id, _ := item["user_id"].(string)
			message_type, _ := item["type"].(string)
			message, _ := item["message"].(string)
			sender_type, _ := item["sender_type"].(string)

			user_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
			if err != nil {
				if valkey.IsValkeyNil(err) {
					user, dbErr := database.FindUserByPublicId(public_user_id)
					if dbErr != nil {
						continue
					}
					if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
						continue
					}
					user_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
					if err != nil {
						continue
					}
				} else {
					continue
				}
			}

			var user sqlc.User
			if err := json.Unmarshal([]byte(user_data), &user); err != nil {
				continue
			}

			polished_data := map[string]any{
				"message_id":      item["public_id"],
				"username":        user.Username,
				"displayname":     user.Displayname,
				"message_type":    message_type,
				"message_content": message,
				"sender_type":     sender_type,
			}

			polished_data_arr = append(polished_data_arr, polished_data)
		}

		c.Set("response_data", polished_data_arr)
		return nil
	}
}

func ClientCommunicationReply() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Session not found"})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ClientCommunicationReplyRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to get request."})
		}

		polished_request, err := validator.ValidateClientCommunicationReplyRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		original_public_id := polished_request.PublicId
		message := polished_request.Message

		ctx := c.Request().Context()
		retrieved_comms_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_support_messages").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve support messages."})
		}

		var all_messages []map[string]any
		if err := json.Unmarshal([]byte(retrieved_comms_data), &all_messages); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse support messages."})
		}

		var original_type string
		for _, m := range all_messages {
			if pid, _ := m["public_id"].(string); pid == original_public_id {
				original_type, _ = m["type"].(string)
				break
			}
		}

		retrieved_user_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve admin data"})
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve admin data"})
				}
				retrieved_user_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve admin data"})
				}
			} else {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve admin data"})
			}
		}

		var user_data sqlc.User
		if err := json.Unmarshal([]byte(retrieved_user_data), &user_data); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse admin data"})
		}

		new_admin_comm := map[string]any{
			"public_id":                 uuid.New().String(),
			"user_id":                   user_data.PublicID.String(),
			"type":                      original_type,
			"message":                   message,
			"sender_type":               "admin",
			"client_comm_log_public_id": original_public_id,
			"created_at":                time.Now(),
		}

		all_messages = append(all_messages, new_admin_comm)

		serialised_data, err := json.Marshal(all_messages)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to serialise support message data"})
		}

		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("client_support_messages").Value(string(serialised_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to store support message data"})
		}

		return nil
	}
}

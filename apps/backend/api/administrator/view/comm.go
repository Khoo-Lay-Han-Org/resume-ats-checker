package administrator_view

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	valkey "github.com/valkey-io/valkey-go"
	typing "resuming/api/administrator/typing"
	validator "resuming/api/administrator/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func GetSupportMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_support_messages").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve support messages."})
			return
		}

		var data []map[string]any
		if err := json.Unmarshal([]byte(retrieved_data), &data); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse support messages."})
			return
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

			var user database.User
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
	}
}

func ClientCommunicationReply() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session not found"})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ClientCommunicationReplyRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to get request."})
			return
		}

		polished_request, err := validator.ValidateClientCommunicationReplyRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		original_public_id := polished_request.PublicId
		message := polished_request.Message

		ctx := c.Request.Context()
		retrieved_comms_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_support_messages").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve support messages."})
			return
		}

		var all_messages []map[string]any
		if err := json.Unmarshal([]byte(retrieved_comms_data), &all_messages); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse support messages."})
			return
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
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve admin data"})
					return
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve admin data"})
					return
				}
				retrieved_user_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve admin data"})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve admin data"})
				return
			}
		}

		var user_data database.User
		if err := json.Unmarshal([]byte(retrieved_user_data), &user_data); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse admin data"})
			return
		}

		new_admin_comm := map[string]any{
			"public_id":                 uuid.New().String(),
			"user_id":                   user_data.PublicId.String(),
			"type":                      original_type,
			"message":                   message,
			"sender_type":               "admin",
			"client_comm_log_public_id": original_public_id,
			"created_at":                time.Now(),
		}

		all_messages = append(all_messages, new_admin_comm)

		serialised_data, err := json.Marshal(all_messages)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to serialise support message data"})
			return
		}

		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("client_support_messages").Value(string(serialised_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store support message data"})
			return
		}

	}
}

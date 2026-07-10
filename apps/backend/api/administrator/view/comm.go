package administrator_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	valkey "github.com/valkey-io/valkey-go"
	typing "resuming/api/administrator/typing"
	validator "resuming/api/administrator/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func GetClientCommunicationLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_comms").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve client communication logs."})
			return
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse client communication logs."})
			return
		}

		var polished_data_arr []map[string]any
		for _, item := range data {
			public_user_id := item["public_user_id"].(uuid.UUID)
			message_type := item["type"].(string)
			message := item["message"].(string)

			user_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id.String()+":user_data").Build()).ToString()
			if err != nil {
				if valkey.IsValkeyNil(err) {
					user, dbErr := database.FindUserByPublicId(public_user_id.String())
					if dbErr != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user."})
						return
					}
					if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id.String(), user); syncErr != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user."})
						return
					}
					user_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id.String()+":user_data").Build()).ToString()
					if err != nil {
						c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user."})
						return
					}
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user."})
					return
				}
			}

			var user database.User
			err = json.Unmarshal([]byte(user_data), &user)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user."})
				return
			}

			polished_data := map[string]any{
				"username":        user.Username,
				"displayname":     user.Displayname,
				"message_type":    message_type,
				"message_content": message,
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

		public_id := polished_request.PublicId
		message := polished_request.Message

		ctx := c.Request.Context()
		retrieved_comms_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_comms").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve admin communication logs."})
			return
		}

		var comms_data []map[string]any
		err = json.Unmarshal([]byte(retrieved_comms_data), &comms_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse admin communication logs."})
			return
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
		err = json.Unmarshal([]byte(retrieved_user_data), &user_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse admin data"})
			return
		}

		new_admin_comm := map[string]any{
			"admin_comm_public_id":      uuid.New().String(),
			"admin_user_public_id":      user_data.PublicId.String(),
			"client_comm_log_public_id": public_id,
			"message":                   message,
		}

		comms_data = append(comms_data, new_admin_comm)

		serialised_admin_comms, err := json.Marshal(comms_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to serialise communication data"})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("admin_comms").Value(string(serialised_admin_comms)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store communication data"})
			return
		}

	}
}

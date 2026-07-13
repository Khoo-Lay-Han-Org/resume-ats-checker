package client_support_view

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	typing "resuming/api/client-support/typing"
	validator "resuming/api/client-support/validator"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func ClientCommunicateToAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to get session data."})
			return
		}
		public_user_id := retrieved_public_user_id.(string)

		var request typing.ClientCommunicateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to get request."})
			return
		}

		polished_request, err := validator.ValidateClientCommunicateRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		type_of_message := polished_request.Type
		client_message := polished_request.Message

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_support_messages").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get support message data."})
			return
		}

		var data []map[string]any
		if err := json.Unmarshal([]byte(retrieved_data), &data); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse data."})
			return
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to serialise support message data."})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("client_support_messages").Value(string(serialised_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to store support message data."})
			return
		}

	}
}

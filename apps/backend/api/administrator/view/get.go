package administrator_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"resuming/tool"
)

func GetClients() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_configs").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve clients config data."})
			return
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to parse clients config data."})
			return
		}

		c.Set("response_data", data)
	}
}

func GetAdmins() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_configs").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve admins config data."})
			return
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to parse admins config data."})
			return
		}

		c.Set("response_data", data)
	}
}

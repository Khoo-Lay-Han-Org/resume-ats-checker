package administrator_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"resuming/database"
	"resuming/tool"
)

func GetClientAuditLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_audit_log_data").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve clients audit logs."})
			return
		}

		var data []database.ClientAuditLog
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to parse client audit logs."})
			return
		}

		c.Set("response_data", data)
	}
}

func GetAdminAuditLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_audit_log_data").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve admin audit logs."})
			return
		}

		var data []database.AdminAuditLog
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to parse admin audit logs."})
			return
		}

		c.Set("response_data", data)
	}
}

func GetErrorAuditLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("error_log_data").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to retrieve error logs."})
			return
		}

		var data []database.ErrorLog
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to parse error logs."})
			return
		}

		c.Set("response_data", data)
	}
}

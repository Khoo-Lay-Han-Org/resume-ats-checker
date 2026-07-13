package client_support_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	typing "resuming/api/client-support/typing"
	validator "resuming/api/client-support/validator"
	"resuming/tool"
)

func ClientReportOtherClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_reporting_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve session id."})
			return
		}

		reporting_public_user_id := retrieved_reporting_public_user_id.(string)

		var request typing.ClientReportRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to retrieve request."})
			return
		}

		polished_request, err := validator.ValidateClientReportRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		report_type := polished_request.ReportType
		target_public_user_id := polished_request.TargetClientPublicUserId

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_report_logs").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user."})
			return
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user data."})
			return
		}

		new_report_log := map[string]any{
			"public_id":                uuid.New(),
			"reporting_public_user_id": reporting_public_user_id,
			"target_public_user_id":    target_public_user_id,
			"type":                     report_type,
		}

		data = append(data, new_report_log)

		serialised_data, err := json.Marshal(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to serialise new report log"})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("client_report_logs").
				Value(string(serialised_data)).
				Build()).
			Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store client report logs"})
			return
		}

	}
}

func AIReport() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

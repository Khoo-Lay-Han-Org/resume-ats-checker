package showcaserecord_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	valkey "github.com/valkey-io/valkey-go"
	util "resuming/api/showcaserecord/util"
	"resuming/database"
	"resuming/tool"
)

func RetrieveShowCaseRecordData() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve showcase record data."})
					return
				}
				if syncErr := database.SyncIndividualShowCaseRecordDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve showcase record data."})
					return
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve showcase record data."})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve showcase record data."})
				return
			}
		}

		var data map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process showcase record data."})
			return
		}

		response_data := gin.H{
			"name":           util.ToStringSlice(data["name"]),
			"email":          util.ToStringSlice(data["email"]),
			"phone_number":   util.ToStringSlice(data["phone_number"]),
			"address":        util.ToStringSlice(data["address"]),
			"social_media":   util.ToStringSlice(data["social_media"]),
			"job_experience": util.ToJSON(data["job_experience"]),
			"education":      util.ToJSON(data["education"]),
			"skill":          util.ToStringSlice(data["skill"]),
			"certificate":    util.ToJSON(data["certificate"]),
			"language":       util.ToStringSlice(data["language"]),
			"project":        util.ToJSON(data["project"]),
		}
		c.Set("response_data", response_data)
	}
}

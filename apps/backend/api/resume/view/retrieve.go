package resume_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	valkey "github.com/valkey-io/valkey-go"
	"gorm.io/datatypes"
	"resuming/database"
	"resuming/tool"
)

func RetrieveResumeData() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve session data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":resume_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve resume data."})
					return
				}
				if syncErr := database.SyncIndividualResumeDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve resume data."})
					return
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":resume_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve resume data."})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve resume data."})
				return
			}
		}

		var data map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process resume data."})
			return
		}

		data_detail := data["detail"].(map[string]any)

		template_id := data["template_id"].(int)
		name := data_detail["name"].([]string)
		email := data_detail["email"].([]string)
		phone_number := data_detail["phone_number"].([]string)
		address := data_detail["address"].([]string)
		social_media := data_detail["social_media"].([]string)
		job_experience := data_detail["job_experience"].(datatypes.JSON)
		education := data_detail["education"].(datatypes.JSON)
		skill := data_detail["skill"].([]string)
		certificate := data_detail["certificate"].(datatypes.JSON)
		language := data_detail["language"].([]string)
		project := data_detail["project"].(datatypes.JSON)

		response_data := gin.H{
			"template_id":    template_id,
			"name":           name,
			"email":          email,
			"phone_number":   phone_number,
			"address":        address,
			"social_media":   social_media,
			"job_experience": job_experience,
			"education":      education,
			"skill":          skill,
			"certificate":    certificate,
			"language":       language,
			"project":        project,
		}
		c.Set("response_data", response_data)
	}
}

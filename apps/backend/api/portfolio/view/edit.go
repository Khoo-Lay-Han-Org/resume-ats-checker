package portfolio_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	valkey "github.com/valkey-io/valkey-go"
	typing "resuming/api/portfolio/typing"
	validator "resuming/api/portfolio/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func ChooseTemplate() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve user data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChooseTemplateRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		template_id, err := validator.ValidateTemplateID(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":portfolio_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve portfolio data."})
					return
				}
				if syncErr := database.SyncIndividualPortfolioDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve portfolio data."})
					return
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":portfolio_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve portfolio data."})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve portfolio data."})
				return
			}
		}

		var data map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process portfolio data."})
			return
		}

		data["template_id"] = template_id

		serialised_data, err := json.Marshal(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process portfolio data."})
			return
		}

		err = tool.Valkey.Do(
			c.Request.Context(),
			tool.Valkey.B().Set().
				Key(public_user_id+":portfolio_data").Value(string(serialised_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to save portfolio data."})
			return
		}

	}
}

package administrator_view

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	typing "resuming/api/administrator/typing"
	validator "resuming/api/administrator/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func BanClient() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request typing.UserControlRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve request."})
			return
		}

		polished_request, err := validator.ValidateUserControlRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		public_user_id := polished_request.PublicUserId

		ctx := c.Request.Context()
		group_data, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("user_data").Build(),
		).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve cached data."})
			return
		}

		var all_users []database.User
		if err := json.Unmarshal([]byte(group_data), &all_users); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse cached data."})
			return
		}

		var target_user *database.User
		now := time.Now()
		for i, user := range all_users {
			if user.PublicId.String() == public_user_id {
				all_users[i].BannedAt = &now
				target_user = &all_users[i]
				break
			}
		}

		if target_user == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}

		serialised_group, err := json.Marshal(all_users)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process updated data."})
			return
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("user_data").Value(string(serialised_group)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user data."})
			return
		}

		individual_data, err := json.Marshal(target_user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to serialise user data"})
			return
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(public_user_id+":user_data").
				Value(string(individual_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to ban user."})
			return
		}

		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id)
	}
}

func RemoveIndividualUserSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request typing.SessionControlRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve request."})
			return
		}

		polished_request, err := validator.ValidateSessionControlRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		public_user_id := polished_request.PublicUserId

		admin_session_id, _ := c.Get("public_user_id")
		if public_user_id == admin_session_id {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Cannot remove your own session."})
			return
		}

		ctx := c.Request.Context()
		exists, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Exists().Key(public_user_id+":session_data").Build(),
		).AsInt64()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify session."})
			return
		}
		if exists == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Session not found."})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Del().
				Key(public_user_id+":session_data",
					public_user_id+":jwt_data",
					public_user_id+":user_data").
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted session"})
	}
}

func RemoveAllClientSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_configs").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to find client configs."})
			return
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse configs."})
			return
		}

		for _, item := range data {
			public_user_id, ok := item["public_user_id"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Invalid client config format."})
				return
			}

			ctx := c.Request.Context()
			_, err := tool.Valkey.Do(
				ctx,
				tool.Valkey.B().Del().
					Key(public_user_id+":session_data").
					Build(),
			).ToString()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete session"})
				return
			}
		}

		admin_id, exists := c.Get("public_user_id")
		if exists {
			admin_key := admin_id.(string) + ":session_data"
			ttl_seconds, err := tool.Valkey.Do(ctx, tool.Valkey.B().Ttl().Key(admin_key).Build()).AsInt64()
			if err == nil && ttl_seconds > 300 {
				tool.Valkey.Do(ctx, tool.Valkey.B().Expire().Key(admin_key).Seconds(300).Build())
			}
		}
	}
}

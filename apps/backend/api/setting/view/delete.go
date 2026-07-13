package setting_view

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	valkey "github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
	typing "resuming/api/setting/typing"
	util "resuming/api/setting/util"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func PrepareDeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to get user data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
					return
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
					return
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
				return
			}
		}

		var user database.User
		err = json.Unmarshal([]byte(retrieved_data), &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user data."})
			return
		}

		err = util.SendOTP(user.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP."})
			return
		}

	}
}

func DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to get user data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
					return
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
					return
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user data."})
				return
			}
		}

		var user database.User
		err = json.Unmarshal([]byte(retrieved_data), &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user data."})
			return
		}

		var request typing.OTPRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		err = util.CheckOTP(user.Email, request.OTP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP."})
			return
		}

		expires_at := time.Now().Add(5 * 365 * 24 * time.Hour)
		new_user_struct := database.User{
			PublicId:    user.PublicId,
			Username:    user.Username,
			Email:       user.Email,
			Displayname: user.Displayname,
			UserType:    user.UserType,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			BannedAt:    user.BannedAt,
			DeletedAt:   gorm.DeletedAt{Time: time.Now(), Valid: true},
			ExpiresAt:   &expires_at,
		}

		serialised_new_user_struct, err := json.Marshal(new_user_struct)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store user data."})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(public_user_id+":user_data").Value(string(serialised_new_user_struct)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store user data."})
			return
		}

	}
}

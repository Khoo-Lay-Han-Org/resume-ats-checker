package setting_view

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	valkey "github.com/valkey-io/valkey-go"
	"golang.org/x/crypto/bcrypt"
	typing "resuming/api/setting/typing"
	util "resuming/api/setting/util"
	validator "resuming/api/setting/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func ChangeUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve session data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangeUsernameRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to receive request."})
			return
		}

		validated_request, err := validator.ValidateUsernameRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		new_data := validated_request.Username

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

		new_user_struct := database.User{
			PublicId:    user.PublicId,
			Username:    new_data,
			Email:       user.Email,
			Displayname: user.Displayname,
			UserType:    user.UserType,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			BannedAt:    user.BannedAt,
			DeletedAt:   user.DeletedAt,
			ExpiresAt:   user.ExpiresAt,
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

func ChangeDisplayname() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to get user data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangeDisplaynameRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to receive request."})
			return
		}

		validated_request, err := validator.ValidateDisplaynameRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		new_data := validated_request.Displayname

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

		new_user_struct := database.User{
			PublicId:    user.PublicId,
			Username:    user.Username,
			Email:       user.Email,
			Displayname: new_data,
			UserType:    user.UserType,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			BannedAt:    user.BannedAt,
			DeletedAt:   user.DeletedAt,
			ExpiresAt:   user.ExpiresAt,
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

func PrepareChangeEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to get user data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangeEmailRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to receive request."})
			return
		}

		validated_request, err := validator.ValidateEmailRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		new_data := validated_request.Email

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

		if new_data == user.Email {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Email is already in use."})
			return
		}

		err = tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().
			Set().
			Key(user.Email+":change-email").
			Value(new_data).
			Build()).
			Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Connection to in-memory data stores failed."})
			return
		}

		err = tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().
			Expire().
			Key(user.Email+":change-email").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Connection to in-memory data stores failed."})
			return
		}

		err = util.SendOTP(user.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP."})
			return
		}

	}
}

func ChangeEmail() gin.HandlerFunc {
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

		new_data, err := tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().Get().Key(user.Email+":change-email").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Email change request expired or not found."})
			return
		}

		new_user_struct := database.User{
			PublicId:    user.PublicId,
			Username:    user.Username,
			Email:       new_data,
			Displayname: user.Displayname,
			UserType:    user.UserType,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			BannedAt:    user.BannedAt,
			DeletedAt:   user.DeletedAt,
			ExpiresAt:   user.ExpiresAt,
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

func PrepareChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_public_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to get user data."})
			return
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangePasswordRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to receive request."})
			return
		}

		validated_request, err := validator.ValidatePasswordRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		new_data := validated_request.Password

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

		var db_user database.User
		result := database.DB.Where("public_id = ?", user.PublicId).First(&db_user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Unable to find user."})
			return
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(new_data), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process request."})
			return
		}

		err = tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().
			Set().
			Key(db_user.Email+":change-password").
			Value(string(hashed_password)).
			Build()).
			Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Connection to in-memory data stores failed."})
			return
		}

		err = tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().
			Expire().
			Key(db_user.Email+":change-password").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Connection to in-memory data stores failed."})
			return
		}

		err = util.SendOTP(db_user.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP."})
			return
		}

	}
}

func ChangePassword() gin.HandlerFunc {
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

		var db_user database.User
		result := database.DB.Where("public_id = ?", user.PublicId).First(&db_user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Unable to find user."})
			return
		}

		var request typing.OTPRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		err = util.CheckOTP(db_user.Email, request.OTP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP."})
			return
		}

		new_password, err := tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().Get().Key(db_user.Email+":change-password").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve new password."})
			return
		}

		result = database.DB.Model(&db_user).Update("password", []byte(new_password))
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update."})
			return
		}

	}
}

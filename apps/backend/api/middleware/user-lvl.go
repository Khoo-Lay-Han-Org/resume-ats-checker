package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
	valkey "github.com/valkey-io/valkey-go"
	"resuming/database"
	"resuming/tool"
)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session not found."})
			return
		}

		public_user_id := retrieved_user_id.(string)

		enforcer, err := casbin.NewEnforcer("api/middleware/config/user-lvl.conf", "api/middleware/config/user-lvl.csv")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve access control configurations."})
			return
		}

		ctx := context.Background()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
					return
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
					return
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
				return
			}
		}

		var user database.User
		err = json.Unmarshal([]byte(retrieved_data), &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user data."})
			return
		}

		ok, err := enforcer.Enforce(string(user.UserType))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to evaluate user accessibility."})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User is not authorised for this process."})
			return
		}

		c.Next()
	}
}

func OnlySuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_user_id, exists := c.Get("public_user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Session not found."})
			return
		}

		public_user_id := retrieved_user_id.(string)

		enforcer, err := casbin.NewEnforcer("api/middleware/config/super-admin-lvl.conf", "api/middleware/config/super-admin-lvl.csv")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve access control configurations."})
			return
		}

		ctx := context.Background()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
					return
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
					return
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
					return
				}
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user data."})
				return
			}
		}

		var user database.User
		err = json.Unmarshal([]byte(retrieved_data), &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user data."})
			return
		}

		ok, err := enforcer.Enforce(string(user.UserType))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to evaluate user accessibility."})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User is not authorised for this process."})
			return
		}

		c.Next()
	}
}

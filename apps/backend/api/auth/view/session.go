package auth_view

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
	util "resuming/api/auth/util"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func SetSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_data, exists := c.Get("private_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve user data."})
			return
		}

		private_id := retrieved_data.(int)

		user_pointer, err := util.FindUser(private_id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}

		user := *user_pointer

		session_key := ksuid.New().String()
		public_user_id := user.PublicId.String()

		claim := jwt.MapClaims{
			"user_public_id": public_user_id,
			"session_key":    session_key,
			"exp":            time.Now().Add(systemconfig.SessionExpiryDuration).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		token_string, err := token.SignedString([]byte(session_key))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate JWT."})
			return
		}

		session_data := map[string]any{
			"public_id":   public_user_id,
			"session_key": session_key,
			"token":       token_string,
		}

		session_json, err := json.Marshal(session_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate session."})
			return
		}

		c.Set("public_user_id", public_user_id)
		c.Set("session_key", session_key)
		c.Set("user", &user)
		c.SetCookie("session", string(session_json), int(systemconfig.SessionExpiryDuration.Seconds()), "/", "", false, true)

		psid := public_user_id
		ctx := context.Background()
		session_store_data, _ := json.Marshal(map[string]string{"session_key": session_key})
		if err := tool.Valkey.Do(ctx, tool.Valkey.B().Set().
			Key(psid+":session_data").
			Value(string(session_store_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
		).Error(); err != nil {
			log.Printf("Failed to write session to Valkey: %v", err)
		}
	}
}

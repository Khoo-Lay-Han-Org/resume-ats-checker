package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	util "resuming/api/middleware/util"
)

func SessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session_cookie, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Failed to retrieve cookie."})
			return
		}

		user_public_id, err := util.CheckSession(session_cookie)
		if err != nil {
			c.SetCookie("session", "deleting cookie", -1, "/", "", false, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired session."})
			return
		}

		c.Set("public_user_id", user_public_id)
		c.Next()
	}
}

package auth_view

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
	util "resuming/api/auth/util"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func SetSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_data := c.Get("private_id")
		if retrieved_data == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve user data."})
		}

		private_id := retrieved_data.(int32)

		user_pointer, err := util.FindUser(private_id)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found."})
		}

		user := *user_pointer

		session_key := ksuid.New().String()
		signing_key := ksuid.New().String()
		public_user_id := user.PublicID.String()

		claim := jwt.MapClaims{
			"user_public_id": public_user_id,
			"exp":            time.Now().Add(systemconfig.SessionExpiryDuration).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		token_string, err := token.SignedString([]byte(signing_key))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to generate JWT."})
		}

		session_data := map[string]any{
			"public_id":   public_user_id,
			"session_key": session_key,
			"token":       token_string,
		}

		session_json, err := json.Marshal(session_data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to generate session."})
		}

		psid := public_user_id
		ctx := c.Request().Context()
		session_store_data, _ := json.Marshal(map[string]string{"session_key": session_key})
		if err := tool.Valkey.Do(ctx, tool.Valkey.B().Set().
			Key(psid+":session_data").
			Value(string(session_store_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
		).Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create session."})
		}

		c.Set("public_user_id", public_user_id)
		c.Set("session_key", session_key)
		c.Set("signing_key", signing_key)
		c.Set("user", &user)
		c.SetCookie(&http.Cookie{
			Name:     "session",
			Value:    string(session_json),
			MaxAge:   int(systemconfig.SessionExpiryDuration.Seconds()),
			Path:     "/",
			Domain:   "",
			Secure:   systemconfig.ApplicationHosted,
			HttpOnly: true,
		})

		return nil
	}
}

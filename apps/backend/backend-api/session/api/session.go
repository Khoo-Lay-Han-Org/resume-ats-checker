package session_api

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
	auth_find "resuming/backend-api/user/find"
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

		user_pointer, err := auth_find.FindUser(private_id)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found."})
		}

		user := *user_pointer

		session_key := ksuid.New().String()
		signing_key := ksuid.New().String()
		public_user_id := user.PublicID.String()

		payload, _ := json.Marshal(map[string]any{
			"user_public_id": public_user_id,
			"exp":            time.Now().Add(systemconfig.SessionExpiryDuration).Unix(),
		})

		key := sha256.Sum256([]byte(signing_key))

		encrypter, err := jose.NewEncrypter(
			jose.A128GCM,
			jose.Recipient{Algorithm: jose.DIRECT, Key: key[:16]},
			nil,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create encrypter."})
		}

		object, err := encrypter.Encrypt(payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to encrypt session."})
		}

		token_string, err := object.CompactSerialize()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to serialize session."})
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

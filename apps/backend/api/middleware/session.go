package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	util "resuming/api/middleware/util"
	systemconfig "resuming/system-config"
)

func SessionCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session_cookie, err := c.Cookie("session")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve cookie."})
			}

			user_public_id, err := util.CheckSession(session_cookie.Value)
			if err != nil {
				c.SetCookie(&http.Cookie{
					Name:   "session",
					Value:  "deleting cookie",
					MaxAge: -1,
					Path:   "/",
					Domain: "",
					Secure: systemconfig.ApplicationHosted,
				})
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid or expired session."})
			}

			c.Set("public_user_id", user_public_id)
			return next(c)
		}
	}
}

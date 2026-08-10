package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
	find "resuming/shared/find"
)

func OnlyAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			retrieved_user_id := c.Get("public_user_id")
			if retrieved_user_id == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Session not found."})
			}
			public_user_id := retrieved_user_id.(string)

			enforcer, err := casbin.NewEnforcer("controller/middleware/config/user-lvl.conf", "controller/middleware/config/user-lvl.csv")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve access control configurations."})
			}

			user, err := find.GetUser(public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
			}

			ok, err := enforcer.Enforce(string(user.UserType))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to evaluate user accessibility."})
			}

			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User is not authorised for this process."})
			}

			return next(c)
		}
	}
}

func OnlySuperAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			retrieved_user_id := c.Get("public_user_id")
			if retrieved_user_id == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Session not found."})
			}
			public_user_id := retrieved_user_id.(string)

			enforcer, err := casbin.NewEnforcer("controller/middleware/config/super-admin-lvl.conf", "controller/middleware/config/super-admin-lvl.csv")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve access control configurations."})
			}

			user, err := find.GetUser(public_user_id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
			}

			ok, err := enforcer.Enforce(string(user.UserType))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to evaluate user accessibility."})
			}

			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "User is not authorised for this process."})
			}

			return next(c)
		}
	}
}

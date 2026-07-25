package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
	valkey "github.com/valkey-io/valkey-go"
	"resuming/database"
	"resuming/database/sqlc"
	"resuming/tool"
)

func OnlyAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			retrieved_user_id := c.Get("public_user_id")
			if retrieved_user_id == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Session not found."})
			}
			public_user_id := retrieved_user_id.(string)

			enforcer, err := casbin.NewEnforcer("api/middleware/config/user-lvl.conf", "api/middleware/config/user-lvl.csv")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve access control configurations."})
			}

			ctx := c.Request().Context()
			retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
			if err != nil {
				if valkey.IsValkeyNil(err) {
					user, dbErr := database.FindUserByPublicId(public_user_id)
					if dbErr != nil {
						return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
					}
					if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
						return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
					}
					retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
					if err != nil {
						return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
					}
				} else {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
				}
			}

			var user sqlc.User
			err = json.Unmarshal([]byte(retrieved_data), &user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user data."})
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

			enforcer, err := casbin.NewEnforcer("api/middleware/config/super-admin-lvl.conf", "api/middleware/config/super-admin-lvl.csv")
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve access control configurations."})
			}

			ctx := c.Request().Context()
			retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
			if err != nil {
				if valkey.IsValkeyNil(err) {
					user, dbErr := database.FindUserByPublicId(public_user_id)
					if dbErr != nil {
						return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
					}
					if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
						return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
					}
					retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
					if err != nil {
						return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
					}
				} else {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user data."})
				}
			}

			var user sqlc.User
			err = json.Unmarshal([]byte(retrieved_data), &user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user data."})
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

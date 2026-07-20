package setting_view

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	valkey "github.com/valkey-io/valkey-go"
	typing "resuming/api/setting/typing"
	util "resuming/api/setting/util"
	"resuming/database"
	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func PrepareDeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
				}
			} else {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
			}
		}

		var user sqlc.User
		err = json.Unmarshal([]byte(retrieved_data), &user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user data."})
		}

		err = util.SendOTP(user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}

		return nil
	}
}

func DeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
		if err != nil {
			if valkey.IsValkeyNil(err) {
				user, dbErr := database.FindUserByPublicId(public_user_id)
				if dbErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
				}
				if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
				}
				retrieved_data, err = tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
				if err != nil {
					return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
				}
			} else {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
			}
		}

		var user sqlc.User
		err = json.Unmarshal([]byte(retrieved_data), &user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse user data."})
		}

		var request typing.OTPRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		err = util.CheckOTP(user.Email, request.OTP)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid OTP."})
		}

		expires_at := time.Now().Add(5 * 365 * 24 * time.Hour)
		new_user_struct := sqlc.User{
			PublicID:    user.PublicID,
			Username:    user.Username,
			Email:       user.Email,
			Displayname: user.Displayname,
			UserType:    user.UserType,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			BannedAt:    user.BannedAt,
			DeletedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			ExpiresAt:   pgtype.Timestamptz{Time: expires_at, Valid: true},
		}

		serialised_new_user_struct, err := json.Marshal(new_user_struct)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to store user data."})
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(public_user_id+":user_data").Value(string(serialised_new_user_struct)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to store user data."})
		}

		return nil
	}
}

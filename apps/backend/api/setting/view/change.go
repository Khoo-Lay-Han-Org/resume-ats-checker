package setting_view

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	valkey "github.com/valkey-io/valkey-go"
	"golang.org/x/crypto/bcrypt"
	typing "resuming/api/setting/typing"
	util "resuming/api/setting/util"
	validator "resuming/api/setting/validator"
	"resuming/database"
	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func ChangeUsername() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve session data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangeUsernameRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidateUsernameRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Username

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

		new_user_struct := sqlc.User{
			PublicID:    user.PublicID,
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

func ChangeDisplayname() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangeDisplaynameRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidateDisplaynameRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Displayname

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

		new_user_struct := sqlc.User{
			PublicID:    user.PublicID,
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

func PrepareChangeEmail() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangeEmailRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidateEmailRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Email

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

		if new_data == user.Email {
			return c.JSON(http.StatusConflict, echo.Map{"message": "Email is already in use."})
		}

		err = tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().
			Set().
			Key(user.Email+":change-email").
			Value(new_data).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().
			Expire().
			Key(user.Email+":change-email").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = util.SendOTP(user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}

		return nil
	}
}

func ChangeEmail() echo.HandlerFunc {
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

		new_data, err := tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().Get().Key(user.Email+":change-email").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Email change request expired or not found."})
		}

		new_user_struct := sqlc.User{
			PublicID:    user.PublicID,
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

func PrepareChangePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request typing.ChangePasswordRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidatePasswordRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Password

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

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(new_data), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process request."})
		}

		err = tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().
			Set().
			Key(user.Email+":change-password").
			Value(string(hashed_password)).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().
			Expire().
			Key(user.Email+":change-password").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = util.SendOTP(user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}

		return nil
	}
}

func ChangePassword() echo.HandlerFunc {
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

		new_password, err := tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().Get().Key(user.Email+":change-password").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve new password."})
		}

		err = database.Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Displayname: user.Displayname,
			Password:    []byte(new_password),
			UserType:    user.UserType,
			BannedAt:    user.BannedAt,
			DeletedAt:   user.DeletedAt,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update."})
		}

		return nil
	}
}

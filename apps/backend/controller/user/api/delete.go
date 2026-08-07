package user_api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	dto "resuming/controller/user/dto"
	email "resuming/controller/user/email"
	otp "resuming/controller/user/otp"
	sms "resuming/controller/user/sms"
	"resuming/database/sqlc"
	"resuming/service"
	shared_find "resuming/shared/find"
	"resuming/systemconfig"
)

func PrepareDeleteAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		factor, ok := supportedTwoFactorType(c)
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Unsupported 2FA type."})
		}

		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}

		if factor == two_factor_sms {
			if user.PhoneNumber == "" {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "No phone number on file."})
			}
			err = sms.SendSMSOTP(user.PhoneNumber)
		} else {
			err = email.SendEmailOTP(user.Email)
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}
		setTwoFactorCookie(c, factor)

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

		factor, err := twoFactorFromCookie(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve 2FA type."})
		}

		ctx := c.Request().Context()
		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}

		var request dto.OTPRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		if factor == two_factor_sms {
			err = otp.CheckSMSOTP(user.PhoneNumber, request.OTP)
		} else {
			err = otp.CheckEmailOTP(user.Email, request.OTP)
		}
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

		err = service.Valkey.Do(
			ctx,
			service.Valkey.B().Set().
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

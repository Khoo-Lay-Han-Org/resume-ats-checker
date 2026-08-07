package user_api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	dto "resuming/controller/user/dto"
	email "resuming/controller/user/email"
	otp "resuming/controller/user/otp"
	sms "resuming/controller/user/sms"
	validator "resuming/controller/user/validator"
	"resuming/database"
	"resuming/database/sqlc"
	"resuming/service"
	shared_find "resuming/shared/find"
	"resuming/systemconfig"
)

func ChangeUsername() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to retrieve session data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request dto.ChangeUsernameRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidateUsernameRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Username

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}
		ctx := c.Request().Context()

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

func ChangeDisplayname() echo.HandlerFunc {
	return func(c echo.Context) error {
		retrieved_public_user_id := c.Get("public_user_id")
		if retrieved_public_user_id == nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Failed to get user data."})
		}

		public_user_id := retrieved_public_user_id.(string)

		var request dto.ChangeDisplaynameRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidateDisplaynameRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Displayname

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}
		ctx := c.Request().Context()

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

func PrepareChangeEmail() echo.HandlerFunc {
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

		var request dto.ChangeEmailRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidateEmailRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Email

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}

		if new_data == user.Email {
			return c.JSON(http.StatusConflict, echo.Map{"message": "Email is already in use."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Set().
			Key(user.Email+":change-email").
			Value(new_data).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Expire().
			Key(user.Email+":change-email").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
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

func ChangeEmail() echo.HandlerFunc {
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

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}
		ctx := c.Request().Context()

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

		new_data, err := service.Valkey.Do(c.Request().Context(), service.Valkey.B().Get().Key(user.Email+":change-email").Build()).ToString()
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

func PrepareChangePhoneNumber() echo.HandlerFunc {
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

		var request dto.ChangePhoneNumberRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidatePhoneNumberRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.PhoneNumber

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}

		if new_data == user.PhoneNumber {
			return c.JSON(http.StatusConflict, echo.Map{"message": "Phone number is already in use."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Set().
			Key(user.Email+":change-phone-number").
			Value(new_data).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Expire().
			Key(user.Email+":change-phone-number").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
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

func ChangePhoneNumber() echo.HandlerFunc {
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

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}
		ctx := c.Request().Context()

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

		new_data, err := service.Valkey.Do(c.Request().Context(), service.Valkey.B().Get().Key(user.Email+":change-phone-number").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Phone number change request expired or not found."})
		}

		new_user_struct := sqlc.User{
			PublicID:    user.PublicID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: new_data,
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

func PrepareChangePassword() echo.HandlerFunc {
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

		var request dto.ChangePasswordRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to receive request."})
		}

		validated_request, err := validator.ValidatePasswordRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		new_data := validated_request.Password

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(new_data), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process request."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Set().
			Key(user.Email+":change-password").
			Value(string(hashed_password)).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Expire().
			Key(user.Email+":change-password").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
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

func ChangePassword() echo.HandlerFunc {
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

		user, err := shared_find.GetUser(public_user_id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to get user data."})
		}
		ctx := c.Request().Context()

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

		new_password, err := service.Valkey.Do(c.Request().Context(), service.Valkey.B().Get().Key(user.Email+":change-password").Build()).ToString()
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

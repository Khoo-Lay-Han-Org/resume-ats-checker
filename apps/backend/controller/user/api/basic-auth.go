package user_api

import (
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
	"resuming/systemconfig"
)

func PrepareRegistration() echo.HandlerFunc {
	return func(c echo.Context) error {
		factor, ok := otp.SupportedTwoFactorType(c)
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Unsupported 2FA type."})
		}
		if factor != otp.Two_factor_email {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Registration only supports email 2FA."})
		}

		var request dto.Register
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		validated_request, err := validator.ValidateRegistration(request)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": err.Error()})
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process request."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Hset().
			Key(request.Email+":session").
			FieldValue().
			FieldValue("username", request.Username).
			FieldValue("password", string(hashed_password)).
			FieldValue("displayname", request.Displayname).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = service.Valkey.Do(c.Request().Context(), service.Valkey.B().
			Expire().
			Key(request.Email+":session").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = email.SendEmailOTP(validated_request.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}

		otp.SetTwoFactorCookie(c, otp.Two_factor_email)
		c.SetCookie(&http.Cookie{
			Name:     "email_for_otp",
			Value:    validated_request.Email,
			MaxAge:   int(systemconfig.OtpExpiryDuration.Seconds()),
			Path:     "/",
			Domain:   "",
			Secure:   systemconfig.ApplicationHosted,
			HttpOnly: true,
		})

		return nil
	}
}

func RegisterUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("email_for_otp")
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve cookie."})
		}
		email := cookie.Value

		var request dto.OTP
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		err = otp.CheckEmailOTP(email, request.OTP)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid OTP."})
		}

		user_details, err := service.Valkey.Do(c.Request().Context(),
			service.Valkey.B().Hgetall().Key(email+":session").Build()).
			ToMap()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user detail."})
		}

		username_msg := user_details["username"]
		displayname_msg := user_details["displayname"]
		password_msg := user_details["password"]

		username, err := (&username_msg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse username."})
		}

		displayname, err := (&displayname_msg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse displayname."})
		}

		hashed_password, err := (&password_msg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse password."})
		}

		_, err = database.Queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
			Username:    username,
			Email:       email,
			Password:    []byte(hashed_password),
			Displayname: displayname,
			UserType:    sqlc.UserTypeClient,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create user."})
		}

		return nil
	}
}

func RegisterAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("email_for_otp")
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve cookie."})
		}
		email := cookie.Value

		var request dto.OTP
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		err = otp.CheckEmailOTP(email, request.OTP)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid OTP."})
		}

		if email != systemconfig.Email {
			return c.JSON(http.StatusForbidden, echo.Map{"message": "Only the system admin can register as admin."})
		}

		count, err := database.Queries.CountUsersByType(c.Request().Context(), sqlc.UserTypeAdmin)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to check admin count."})
		}
		super_admin_count, err := database.Queries.CountUsersByType(c.Request().Context(), sqlc.UserTypeSuperAdmin)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to check admin count."})
		}
		if count+super_admin_count > 0 {
			return c.JSON(http.StatusForbidden, echo.Map{"message": "Admin already registered."})
		}

		user_details, err := service.Valkey.Do(c.Request().Context(),
			service.Valkey.B().Hgetall().Key(email+":session").Build()).
			ToMap()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user detail."})
		}

		username_msg := user_details["username"]
		displayname_msg := user_details["displayname"]
		password_msg := user_details["password"]

		username, err := (&username_msg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse username."})
		}

		displayname, err := (&displayname_msg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse displayname."})
		}

		hashed_password, err := (&password_msg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse password."})
		}

		_, err = database.Queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
			Username:    username,
			Email:       email,
			Password:    []byte(hashed_password),
			Displayname: displayname,
			UserType:    sqlc.UserTypeSuperAdmin,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create user."})
		}

		return nil
	}
}

func PrepareLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		factor, ok := otp.SupportedTwoFactorType(c)
		if !ok {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Unsupported 2FA type."})
		}

		var request dto.Login
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		validated_request, err := validator.ValidateLogin(request)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": err.Error()})
		}

		user, err := database.Queries.FindUserByEmail(c.Request().Context(), validated_request.Email)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User does not exist."})
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid password."})
		}

		if factor == otp.Two_factor_sms {
			if user.PhoneNumber == "" {
				return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "No phone number on file."})
			}
			err = sms.SendSMSOTP(user.PhoneNumber)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
			}
			otp.SetTwoFactorCookie(c, otp.Two_factor_sms)
			c.SetCookie(&http.Cookie{
				Name:     "email_for_otp",
				Value:    validated_request.Email,
				MaxAge:   int(systemconfig.OtpExpiryDuration.Seconds()),
				Path:     "/",
				Domain:   "",
				Secure:   systemconfig.ApplicationHosted,
				HttpOnly: true,
			})
			return nil
		}

		err = email.SendEmailOTP(validated_request.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}

		otp.SetTwoFactorCookie(c, otp.Two_factor_email)
		c.SetCookie(&http.Cookie{
			Name:     "email_for_otp",
			Value:    validated_request.Email,
			MaxAge:   int(systemconfig.OtpExpiryDuration.Seconds()),
			Path:     "/",
			Domain:   "",
			Secure:   systemconfig.ApplicationHosted,
			HttpOnly: true,
		})

		return nil
	}
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("email_for_otp")
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve cookie."})
		}
		email := cookie.Value

		factor, err := otp.TwoFactorFromCookie(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve 2FA type."})
		}

		var request dto.OTP
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		ctx := c.Request().Context()
		user, err := database.Queries.FindUserByEmail(ctx, email)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve user."})
		}

		if factor == otp.Two_factor_sms {
			err = otp.CheckSMSOTP(user.PhoneNumber, request.OTP)
		} else {
			err = otp.CheckEmailOTP(email, request.OTP)
		}
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid OTP"})
		}

		c.Set("private_id", user.ID)
		return nil
	}
}

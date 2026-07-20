package auth_view

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	typing "resuming/api/auth/typing"
	util "resuming/api/auth/util"
	validator "resuming/api/auth/validator"
	"resuming/database"
	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func PrepareRegistration() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request typing.Register
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		validated_request, err := validator.ValidateRegistration(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		hashed_password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process request."})
		}

		err = tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().
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

		err = tool.Valkey.Do(c.Request().Context(), tool.Valkey.B().
			Expire().
			Key(request.Email+":session").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Connection to in-memory data stores failed."})
		}

		err = util.SendOTP(validated_request.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}

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

func Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		type_of_user := c.Param("type-of-user")

		cookie, err := c.Cookie("email_for_otp")
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve cookie."})
		}
		email := cookie.Value

		var request typing.OTP
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		err = util.CheckOTP(email, request.OTP)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid OTP."})
		}

		user_details, err := tool.Valkey.Do(c.Request().Context(),
			tool.Valkey.B().Hgetall().Key(email+":session").Build()).
			ToMap()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve user detail."})
		}

		usernameMsg := user_details["username"]
		displaynameMsg := user_details["displayname"]
		passwordMsg := user_details["password"]

		username, err := (&usernameMsg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse username."})
		}

		displayname, err := (&displaynameMsg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse displayname."})
		}

		hashed_password, err := (&passwordMsg).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse password."})
		}

		var user_type sqlc.UserType
		switch type_of_user {
		case "client":
			user_type = sqlc.UserTypeClient
		case "admin":
			if email != systemconfig.Email {
				return c.JSON(http.StatusForbidden, echo.Map{"message": "Only the system admin can register as admin."})
			}
			count, err := database.Queries.CountUsersByType(c.Request().Context(), sqlc.UserTypeAdmin)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to check admin count."})
			}
			superAdminCount, err := database.Queries.CountUsersByType(c.Request().Context(), sqlc.UserTypeSuperAdmin)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to check admin count."})
			}
			if count+superAdminCount > 0 {
				return c.JSON(http.StatusForbidden, echo.Map{"message": "Admin already registered."})
			}
			user_type = sqlc.UserTypeSuperAdmin
		default:
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid account type."})
		}

		_, err = database.Queries.CreateUser(c.Request().Context(), sqlc.CreateUserParams{
			Username:    username,
			Email:       email,
			Password:    []byte(hashed_password),
			Displayname: displayname,
			UserType:    user_type,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create user."})
		}

		return nil
	}
}

func PrepareLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request typing.Login
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		validated_request, err := validator.ValidateLogin(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		user, err := database.Queries.FindUserByEmail(c.Request().Context(), validated_request.Email)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User does not exist."})
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid password."})
		}

		err = util.SendOTP(validated_request.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send OTP."})
		}

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

		var request typing.OTP
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"message": "Failed to process request."})
		}

		ctx := c.Request().Context()
		value, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(email+":otp").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process OTP."})
		}
		err = bcrypt.CompareHashAndPassword([]byte(value), []byte(request.OTP))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid OTP"})
		}

		err = tool.Valkey.Do(ctx, tool.Valkey.B().Del().Key(email+":otp").Build()).Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process OTP"})
		}

		user, err := database.Queries.FindUserByEmail(ctx, email)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve user."})
		}

		c.Set("private_id", user.ID)
		return nil
	}
}

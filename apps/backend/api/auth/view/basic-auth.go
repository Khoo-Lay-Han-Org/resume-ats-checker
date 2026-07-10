package auth_view

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	typing "resuming/api/auth/typing"
	util "resuming/api/auth/util"
	validator "resuming/api/auth/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func PrepareRegistration() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get request
		var request typing.Register
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		// validate data
		validated_request, err := validator.ValidateRegistration(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// store the registeration data for after validating OTP
		hashed_password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process request."})
			return
		}

		err = tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().
			Hset().
			Key(request.Email+":session").
			FieldValue().
			FieldValue("username", request.Username).
			FieldValue("password", string(hashed_password)).
			FieldValue("displayname", request.Displayname).
			Build()).
			Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Connection to in-memory data stores failed."})
			return
		}

		err = tool.Valkey.Do(c.Request.Context(), tool.Valkey.B().
			Expire().
			Key(request.Email+":session").
			Seconds(int64(systemconfig.OtpExpiryDuration.Seconds())).
			Build()).
			Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Connection to in-memory data stores failed."})
			return
		}

		// send OTP
		err = util.SendOTP(validated_request.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP."})
			return
		}

		// send cookie
		c.SetCookie("email_for_otp", validated_request.Email, int(systemconfig.OtpExpiryDuration.Seconds()), "/", "", false, true)
	}
}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		type_of_user := c.Param("type-of-user")

		// get cookie
		email, err := c.Cookie("email_for_otp")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cookie."})
			return
		}

		// get request data (OTP)
		var request typing.OTP
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		err = util.CheckOTP(email, request.OTP)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid OTP."})
			return
		}

		// if OTP valid, create new user
		user_details, err := tool.Valkey.Do(c.Request.Context(),
			tool.Valkey.B().Hgetall().Key(email+":session").Build()).
			ToMap()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve user detail."})
			return
		}

		usernameMsg := user_details["username"]
		displaynameMsg := user_details["displayname"]
		passwordMsg := user_details["password"]

		username, err := (&usernameMsg).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse username."})
			return
		}

		displayname, err := (&displaynameMsg).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse displayname."})
			return
		}

		hashed_password, err := (&passwordMsg).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse password."})
			return
		}

		var user_type database.UserType
		switch type_of_user {
		case "client":
			user_type = database.Client
		case "admin":
			if email != systemconfig.Email {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Only the system admin can register as admin."})
				return
			}
			var count int64
			database.DB.Model(&database.User{}).
				Where("user_type IN ?", []database.UserType{database.Admin, database.SuperAdmin}).
				Count(&count)
			if count > 0 {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Admin already registered."})
				return
			}
			user_type = database.SuperAdmin
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid account type."})
			return
		}

		new_user := database.User{
			Displayname: displayname,
			Username:    username,
			Email:       email,
			Password:    []byte(hashed_password),
			UserType:    user_type,
		}

		result := database.DB.Create(&new_user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
			return
		}

	}
}

func PrepareLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get request
		var request typing.Login
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		// validate input
		validated_request, err := validator.ValidateLogin(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// check password
		var user database.User
		result := database.DB.Where("email = ?", validated_request.Email).First(&user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User does not exist."})
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(request.Password))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid password."})
			return
		}

		// send OTP
		err = util.SendOTP(validated_request.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to send OTP."})
			return
		}

		// set cookie
		c.SetCookie("email_for_otp", validated_request.Email, int(systemconfig.OtpExpiryDuration.Seconds()), "/", "", false, true)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get cookie
		email, err := c.Cookie("email_for_otp")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve cookie."})
			return
		}

		// get request
		var request typing.OTP
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to process request."})
			return
		}

		// check OTP matching
		ctx := c.Request.Context()
		value, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(email+":otp").Build()).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process OTP."})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(value), []byte(request.OTP))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
			return
		}

		err = tool.Valkey.Do(ctx, tool.Valkey.B().Del().Key(email+":otp").Build()).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process OTP"})
			return
		}

		// find user
		var user database.User
		result := database.DB.Where("email = ?", email).First(&user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve user."})
			return
		}

		c.Set("private_id", user.Id)
		c.Next()
	}
}

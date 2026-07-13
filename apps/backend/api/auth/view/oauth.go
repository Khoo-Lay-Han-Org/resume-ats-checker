package auth_view

/*
import (
	"context"
	"crypto/rand"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	typing "resuming/api/auth/typing"
	util "resuming/api/auth/util"
	"resuming/database"
	"resuming/env"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func InitiateOAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		provider := c.Param("provider")

		var config *oauth2.Config
		switch provider {
		case "google":
			config = &oauth2.Config{
				ClientID:     systemconfig.Google_client_id,
				ClientSecret: systemconfig.Google_client_secret,
				RedirectURL:  systemconfig.Google_redirect_uri,
				Endpoint:     google.Endpoint,
				Scopes:       []string{"email", "profile"},
			}
		case "facebook":
			config = &oauth2.Config{
				ClientID:     systemconfig.Facebook_client_id,
				ClientSecret: systemconfig.Facebook_client_secret,
				RedirectURL:  systemconfig.Facebook_redirect_uri,
				Endpoint:     facebook.Endpoint,
				Scopes:       []string{"email"},
			}
		case "github":
			config = &oauth2.Config{
				ClientID:     systemconfig.Github_client_id,
				ClientSecret: systemconfig.Github_client_secret,
				RedirectURL:  systemconfig.Github_redirect_uri,
				Endpoint:     github.Endpoint,
				Scopes:       []string{"user:email"},
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Unsupported provider."})
			return
		}

		// create oauth session key
		state_key := make([]byte, 20)

		_, err := rand.Read(state_key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create state."})
		}

		state := hex.EncodeToString(state_key)

		c.Set("oauth-config", config)
		c.SetCookie("oauth_status", state, 300, "/", "", false, true)

		// trigger oauth
		url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func OAuthCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		retrieved_config, exists := c.Get("oauth-config")
		if !exists {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve OAuth detail."})
			return
		}

		config, ok := retrieved_config.(*oauth2.Config)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to run OAuth."})
			return
		}

		state, err := c.Cookie("oauth_status")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch cookie."})
			return
		}

		// check session key matching for security
		stored_state := c.Query("status")
		if stored_state != state {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Oauth session expired."})
			return
		}

		// query code that represent user consent
		code := c.Query("code")
		if code == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "No authorisation code found."})
			return
		}

		// exchange code to token that represent the permission to process
		ctx := context.Background()
		token, err := config.Exchange(ctx, code)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange token."})
			return
		}

		// call for the permission to use the oauth provider's client
		client := config.Client(ctx, token)
		// retrieve response
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user info."})
			return
		}
		defer resp.Body.Close()

		// parse response
		var user_info typing.OAuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&user_info); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user info."})
			return
		}

		c.Set("oauth_user_detail", user_info)
		c.Next()
	}
}

func OAuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_info, exists := c.Get("oauth_user_detail")
		if !exists {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to fetch user detail."})
			return
		}

		user_detail, ok := user_info.(typing.OAuthResponse)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user detail."})
			return
		}

		var user database.User
		result := database.DB.Where("email = ?", user_detail.Email).First(user)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to find user."})
			return
		}

		// Set session expiry time (3 days)
		session_expiry_time := time.Hour * 24 * 3

		// Generate secret for JWT key
		secret := make([]byte, 20)
		_, err := crypto_rand.Read(secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}

		// Generate ID for session store data
		session_store_id := ksuid.New().String()

		// Prepare data to be stored in session store
		user_data := map[string]interface{}{
			"user_data": map[string]string{
				"username":    user.Username,
				"displayname": user.Displayname,
				"email":       user.Email,
			},
		}
		portfolio_data := map[string]interface{}{
			"portfolio_data": user.Portfolio,
		}
		resume_data := map[string]interface{}{
			"resume_data": user.Resume,
		}
		ats_data := map[string]interface{}{
			"ats_data": user.ATS,
		}

		// Serialise details
		serialised_user_data, err := json.Marshal(user_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}
		serialised_portfolio_data, err := json.Marshal(portfolio_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}
		serialised_resume_data, err := json.Marshal(resume_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}
		serialised_ats_data, err := json.Marshal(ats_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}

		ctx := context.Background()
		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(session_store_id+":user_data").Value(string(serialised_user_data)).
				Ex(session_expiry_time).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(session_store_id+":portfolio_data").Value(string(serialised_portfolio_data)).
				Ex(session_expiry_time).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(session_store_id+":resume_data").Value(string(serialised_resume_data)).
				Ex(session_expiry_time).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(session_store_id+":ats_data").Value(string(serialised_ats_data)).
				Ex(session_expiry_time).
				Build(),
		).Error()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}

		// Store JWT key into database
		new_jwt_key := database.JwtKey{
			Key: session_store_id,
		}

		result = database.DB.Create(&new_jwt_key)
		if result.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			return
		}

		// Create JWT
		claim := jwt.MapClaims{
			"session_store_id": session_store_id,
			"exp":              time.Now().Add(session_expiry_time).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		token_string, err := token.SignedString(secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate JWT."})
			return
		}

		// create session
		session_data := map[string]interface{}{
			"public_id": user.PublicId,
			"token":     token_string,
		}

		session_json, err := json.Marshal(session_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate session."})
			return
		}

		c.SetCookie("session", string(session_json), int(systemconfig.SessionExpiryDuration.Seconds()), "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in."})
		return
	}

}

func OAuthRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_info, exists := c.Get("oauth_user_detail")
		if !exists {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to fetch user detail."})
			return
		}

		user_detail, ok := user_info.(typing.OAuthResponse)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse user detail."})
			return
		}

		var user database.User
		result := database.DB.Where("email = ?", user_detail.Email).First(user)
		if result.Error != nil {
			err, password := util.GenerateRandomPassword()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
				return
			}

			new_user := database.User{
				Username: strings.ToLower(user_detail.Name),
				Email:    user_detail.Email,
				Password: password,
			}

			result = database.DB.Create(&new_user)
			if result.Error != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "User already exists."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully created user."})
		return
	}
}
*/

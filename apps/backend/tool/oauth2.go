package tool

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"resuming/env"
)

func GetGoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     env.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: env.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  env.GetEnv("BACKEND_URI") + "oauth/google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

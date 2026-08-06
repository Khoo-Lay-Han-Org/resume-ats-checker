package systemconfig

import "resuming/env"

var TWILIOAccountSID string
var TWILIOAuthToken string
var TWILIOPhoneNumber string

func init() {
	TWILIOAccountSID = env.GetEnv("TWILIO_ACCOUNT_SID")
	TWILIOAuthToken = env.GetEnv("TWILIO_AUTH_TOKEN")
	TWILIOPhoneNumber = env.GetEnv("TWILIO_PHONE_NUMBER")
}

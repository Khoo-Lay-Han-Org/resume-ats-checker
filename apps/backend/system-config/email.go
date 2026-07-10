package systemconfig

import "resuming/env"

var Email string
var MailgunDomain string
var MailgunAPIKey string

func init() {
	Email = env.GetEnv("EMAIL")
	MailgunDomain = env.GetEnv("MAILGUN_DOMAIN")
	MailgunAPIKey = env.GetEnv("MAILGUN_API_KEY")
}

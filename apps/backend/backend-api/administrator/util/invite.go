package administrator_util

import (
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func EmailInvitationToBecomeAdmin(email, token string) error {
	invite_link := systemconfig.BackendUri + "/accept-become-admin/" + token
	email_message := `
		<!DOCTYPE html>
		<html>
		<body>
			<a href="` + invite_link + `">ACCEPT</a>
		</body>
		</html>
	`

	return tool.SendEmail(email, "Your Invitation to Become an Admin for Resuming", email_message, true)
}

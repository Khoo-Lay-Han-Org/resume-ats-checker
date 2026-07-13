package systemconfig

import (
	"strings"

	"resuming/env"
)

var FrontendDomain string
var FrontendPort string
var FrontendUri string

func init() {
	FrontendDomain = env.GetEnv("FRONTEND_DOMAIN")
	if !strings.HasSuffix(FrontendDomain, ":") {
		new_string := strings.TrimSpace(FrontendDomain) + ":"
		FrontendDomain = new_string
	}
	FrontendPort = env.GetEnv("FRONTEND_PORT")

	if ApplicationHosted {
		FrontendUri = strings.TrimSuffix(FrontendDomain, "/")
	} else if !ApplicationHosted {
		FrontendUri = strings.TrimSuffix(FrontendDomain+FrontendPort, "/")
	}
}

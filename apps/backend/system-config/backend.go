package systemconfig

import (
	"strings"

	"resuming/env"
)

var BackendDomain string
var BackendPort string
var BackendUri string

func init() {
	BackendDomain = env.GetEnv("BACKEND_DOMAIN")
	if !strings.HasSuffix(BackendDomain, ":") {
		new_string := strings.TrimSpace(BackendDomain) + ":"
		BackendDomain = new_string
	}

	// only to satisfy PaaS (Render) health check (start)
	port := env.GetEnv("PORT")
	if port == "" {
		port = env.GetEnv("BACKEND_PORT")
		if port == "" {
			port = "5321"
		}
	}
	// only to satisfy PaaS health check (end)

	BackendPort = port

	if ApplicationHosted {
		BackendUri = strings.TrimSuffix(BackendDomain, "/")
	} else if !ApplicationHosted {
		BackendUri = strings.TrimSuffix(BackendDomain+BackendPort, "/")
	}
}

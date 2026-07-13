package systemconfig

import (
	"strings"

	"resuming/env"
)

var ValkeyDomain string
var ValkeyPort string
var ValkeyPassword string
var ValkeyTLS bool
var ValkeyUri string

func init() {
	ValkeyDomain = env.GetEnv("VALKEY_DOMAIN")
	if !strings.HasSuffix(ValkeyDomain, ":") {
		new_string := strings.TrimSpace(ValkeyDomain) + ":"
		ValkeyDomain = new_string
	}
	ValkeyPort = env.GetEnv("VALKEY_PORT")
	ValkeyPassword = env.GetEnv("VALKEY_PASSWORD")
	retrieved_valkey_tls := env.GetEnv("VALKEY_TLS")
	if retrieved_valkey_tls == "true" || retrieved_valkey_tls == "True" || retrieved_valkey_tls == "TRUE" {
		ValkeyTLS = true
	} else {
		ValkeyTLS = false
	}

	ValkeyUri = strings.TrimSuffix(ValkeyDomain+ValkeyPort, "/")
}

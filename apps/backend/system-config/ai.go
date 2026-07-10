package systemconfig

import (
	"strings"

	"resuming/env"
)

var AiModelsDomain string
var AiModelsPort string
var AiModelsUri string

func init() {
	AiModelsDomain = env.GetEnv("AI_MODELS_DOMAIN")
	if !strings.HasSuffix(AiModelsDomain, ":") {
		new_string := strings.TrimSpace(AiModelsDomain) + ":"
		AiModelsDomain = new_string
	}
	AiModelsPort = env.GetEnv("AI_MODELS_PORT")
	if !strings.HasSuffix(AiModelsPort, "/") {
		new_string := strings.TrimSpace(AiModelsPort) + "/"
		AiModelsPort = new_string
	}

	if ApplicationHosted {
		AiModelsUri = AiModelsDomain
	} else if !ApplicationHosted {
		AiModelsUri = AiModelsDomain + AiModelsPort
	}
}

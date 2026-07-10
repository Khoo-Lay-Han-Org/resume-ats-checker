package client_support_validator

import (
	"github.com/bobch27/valtra-go"
	typing "resuming/api/client-support/typing"
)

func ValidateClientCommunicateRequest(request typing.ClientCommunicateRequest) (typing.ClientCommunicateRequest, error) {
	v := valtra.NewCollector()

	client_communicate_request := typing.ClientCommunicateRequest{
		Type: valtra.Val(request.Type, "Message type").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Message type is required."),
			).
			Collect(v),

		Message: valtra.Val(request.Message, "Message").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Message is required."),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return client_communicate_request, nil
}

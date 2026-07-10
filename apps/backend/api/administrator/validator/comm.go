package administrator_validator

import (
	"github.com/bobch27/valtra-go"
	typing "resuming/api/administrator/typing"
)

func ValidateClientCommunicationReplyRequest(request typing.ClientCommunicationReplyRequest) (typing.ClientCommunicationReplyRequest, error) {
	v := valtra.NewCollector()

	client_communication_reply_request := typing.ClientCommunicationReplyRequest{
		PublicId: valtra.Val(request.PublicId, "Public ID").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Public ID is required."),
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

	return client_communication_reply_request, nil
}

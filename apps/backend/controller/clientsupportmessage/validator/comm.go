package clientsupportmessage_validator

import (
	"github.com/bobch27/valtra-go"
	dto "resuming/controller/clientsupportmessage/dto"
)

func ValidateClientCommunicationReplyRequest(request dto.ClientCommunicationReplyRequest) (dto.ClientCommunicationReplyRequest, error) {
	v := valtra.NewCollector()

	client_communication_reply_request := dto.ClientCommunicationReplyRequest{
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

func ValidateClientCommunicateRequest(request dto.ClientCommunicateRequest) (dto.ClientCommunicateRequest, error) {
	v := valtra.NewCollector()

	client_communicate_request := dto.ClientCommunicateRequest{
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

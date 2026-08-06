package clientsupportmessage_validator

import (
	dto "resuming/controller/clientsupportmessage/dto"
	"resuming/controller/validation"
)

func ValidateClientCommunicationReplyRequest(request dto.ClientCommunicationReplyRequest) (dto.ClientCommunicationReplyRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ClientCommunicationReplyRequest), nil
}

func ValidateClientCommunicateRequest(request dto.ClientCommunicateRequest) (dto.ClientCommunicateRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ClientCommunicateRequest), nil
}

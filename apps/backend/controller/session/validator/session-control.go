package session_validator

import (
	dto "resuming/controller/session/dto"
	"resuming/controller/validation"
)

func ValidateSessionControlRequest(request dto.SessionControlRequest) (dto.SessionControlRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.SessionControlRequest), nil
}

package user_validator

import (
	dto "resuming/controller/user/dto"
	"resuming/shared/validation"
)

func ValidateUserControlRequest(request dto.UserControlRequest) (dto.UserControlRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.UserControlRequest), nil
}

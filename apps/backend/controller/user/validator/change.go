package user_validator

import (
	dto "resuming/controller/user/dto"
	"resuming/shared/validation"
)

func ValidateUsernameRequest(request dto.ChangeUsernameRequest) (dto.ChangeUsernameRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ChangeUsernameRequest), nil
}

func ValidateDisplaynameRequest(request dto.ChangeDisplaynameRequest) (dto.ChangeDisplaynameRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ChangeDisplaynameRequest), nil
}

func ValidateEmailRequest(request dto.ChangeEmailRequest) (dto.ChangeEmailRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ChangeEmailRequest), nil
}

func ValidatePhoneNumberRequest(request dto.ChangePhoneNumberRequest) (dto.ChangePhoneNumberRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ChangePhoneNumberRequest), nil
}

func ValidatePasswordRequest(request dto.ChangePasswordRequest) (dto.ChangePasswordRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ChangePasswordRequest), nil
}

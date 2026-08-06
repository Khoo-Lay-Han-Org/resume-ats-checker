package user_validator

import (
	dto "resuming/controller/user/dto"
	"resuming/controller/validation"
)

func ValidateRegistration(request dto.Register) (dto.Register, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.Register{}, err
	}
	return cleaned.(dto.Register), nil
}

func ValidateLogin(request dto.Login) (dto.Login, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.Login{}, err
	}
	return cleaned.(dto.Login), nil
}

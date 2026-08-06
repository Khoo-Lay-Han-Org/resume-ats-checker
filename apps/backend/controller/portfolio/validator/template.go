package portfolio_validator

import (
	dto "resuming/controller/portfolio/dto"
	"resuming/controller/validation"
)

func ValidateTemplateID(request dto.ChooseTemplateRequest) (dto.ChooseTemplateRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.ChooseTemplateRequest{}, err
	}
	return cleaned.(dto.ChooseTemplateRequest), nil
}

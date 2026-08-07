package resume_validator

import (
	dto "resuming/controller/resume/dto"
	"resuming/shared/validation"
)

func ValidateTemplateID(request dto.ChooseTemplateRequest) (dto.ChooseTemplateRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return dto.ChooseTemplateRequest{}, err
	}
	return cleaned.(dto.ChooseTemplateRequest), nil
}

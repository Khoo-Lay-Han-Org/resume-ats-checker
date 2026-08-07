package showcase_validator

import (
	dto "resuming/controller/showcase/dto"
	"resuming/shared/validation"
)

func ValidateSpecificPortfolioDataRequest(request dto.SpecificPortoflioDataRequest) (dto.SpecificPortoflioDataRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.SpecificPortoflioDataRequest), nil
}

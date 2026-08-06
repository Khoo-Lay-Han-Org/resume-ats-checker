package portfolio_validator

import (
	"github.com/bobch27/valtra-go"
	dto "resuming/controller/portfolio/dto"
)

func ValidateTemplateID(request dto.ChooseTemplateRequest) (dto.ChooseTemplateRequest, error) {
	v := valtra.NewCollector()

	template_id := dto.ChooseTemplateRequest{
		TemplateId: valtra.Val(request.TemplateId, "Template ID").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Template ID is required."),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return dto.ChooseTemplateRequest{}, v.Errors()[0]
	}

	return template_id, nil
}

package portfolio_validator

import (
	"github.com/bobch27/valtra-go"
	typing "resuming/api/portfolio/typing"
)

func ValidateTemplateID(request typing.ChooseTemplateRequest) (typing.ChooseTemplateRequest, error) {
	v := valtra.NewCollector()

	template_id := typing.ChooseTemplateRequest{
		TemplateId: valtra.Val(request.TemplateId, "Template ID").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Template ID is required."),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return typing.ChooseTemplateRequest{}, v.Errors()[0]
	}

	return template_id, nil
}

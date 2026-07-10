package showcaserecord_validator

import (
	"github.com/bobch27/valtra-go"
	typing "resuming/api/showcaserecord/typing"
)

func ValidateSpecificPortfolioDataRequest(request typing.SpecificPortoflioDataRequest) (typing.SpecificPortoflioDataRequest, error) {
	v := valtra.NewCollector()

	specific_portfolio_request := typing.SpecificPortoflioDataRequest{
		SectionTitle: valtra.Val(request.SectionTitle, "Section title").
			Transform(valtra.TrimSpace(), valtra.Lowercase()).
			Validate(
				valtra.Required[string]("Section title is required."),
				valtra.MinLengthString(3, "Section title must be at least 3 characters"),
				valtra.MaxLengthString(300, "Section title must be at most 300 characters"),
			).
			Collect(v),

		Index: request.Index,
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return specific_portfolio_request, nil
}

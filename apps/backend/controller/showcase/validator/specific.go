package showcase_validator

import (
	"github.com/bobch27/valtra-go"
	dto "resuming/controller/showcase/dto"
)

func ValidateSpecificPortfolioDataRequest(request dto.SpecificPortoflioDataRequest) (dto.SpecificPortoflioDataRequest, error) {
	v := valtra.NewCollector()

	specific_portfolio_request := dto.SpecificPortoflioDataRequest{
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

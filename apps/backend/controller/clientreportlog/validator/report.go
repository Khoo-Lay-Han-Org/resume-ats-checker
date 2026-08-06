package clientreportlog_validator

import (
	"github.com/bobch27/valtra-go"
	dto "resuming/controller/clientreportlog/dto"
)

func ValidateClientReportRequest(request dto.ClientReportRequest) (dto.ClientReportRequest, error) {
	v := valtra.NewCollector()

	client_report_request := dto.ClientReportRequest{
		TargetClientPublicUserId: valtra.Val(request.TargetClientPublicUserId, "Target client").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Target client public session ID is required."),
			).
			Collect(v),

		ReportType: valtra.Val(request.ReportType, "Report type").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Report type is required."),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return client_report_request, nil
}

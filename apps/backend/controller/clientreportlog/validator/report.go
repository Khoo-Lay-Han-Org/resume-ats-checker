package clientreportlog_validator

import (
	dto "resuming/controller/clientreportlog/dto"
	"resuming/controller/validation"
)

func ValidateClientReportRequest(request dto.ClientReportRequest) (dto.ClientReportRequest, error) {
	cleaned, err := validation.TransformAndValidate(request)
	if err != nil {
		return request, err
	}
	return cleaned.(dto.ClientReportRequest), nil
}

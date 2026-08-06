package session_validator

import (
	"github.com/bobch27/valtra-go"
	dto "resuming/controller/session/dto"
)

func ValidateSessionControlRequest(request dto.SessionControlRequest) (dto.SessionControlRequest, error) {
	v := valtra.NewCollector()

	session_control_request := dto.SessionControlRequest{
		PublicUserId: valtra.Val(request.PublicUserId, "Public session ID").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Public session ID is required."),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return session_control_request, nil
}

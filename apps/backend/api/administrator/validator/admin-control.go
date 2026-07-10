package administrator_validator

import (
	"github.com/bobch27/valtra-go"
	typing "resuming/api/administrator/typing"
)

func ValidateUserControlRequest(request typing.UserControlRequest) (typing.UserControlRequest, error) {
	v := valtra.NewCollector()

	user_control_request := typing.UserControlRequest{
		PublicUserId: valtra.Val(request.PublicUserId, "Public user ID").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Public user ID is required."),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return user_control_request, nil
}

func ValidateSessionControlRequest(request typing.SessionControlRequest) (typing.SessionControlRequest, error) {
	v := valtra.NewCollector()

	session_control_request := typing.SessionControlRequest{
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

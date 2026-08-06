package user_validator

import (
	"github.com/bobch27/valtra-go"
	dto "resuming/controller/user/dto"
)

func ValidateUserControlRequest(request dto.UserControlRequest) (dto.UserControlRequest, error) {
	v := valtra.NewCollector()

	user_control_request := dto.UserControlRequest{
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

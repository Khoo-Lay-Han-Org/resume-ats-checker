package setting_validator

import (
	"fmt"

	"github.com/bobch27/valtra-go"
	setting_util "resuming/api/setting/util"
	typing "resuming/api/setting/typing"
)

func ValidateUsernameRequest(request typing.ChangeUsernameRequest) (typing.ChangeUsernameRequest, error) {
	v := valtra.NewCollector()

	username := typing.ChangeUsernameRequest{
		Username: valtra.Val(request.Username, "Username").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Username is required."),
				valtra.MinLengthString(4, "Username must be at least 4 characters"),
				valtra.MaxLengthString(255, "Username must be at most 255 characters")).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return username, nil
}

func ValidateDisplaynameRequest(request typing.ChangeDisplaynameRequest) (typing.ChangeDisplaynameRequest, error) {
	v := valtra.NewCollector()

	displayname := typing.ChangeDisplaynameRequest{
		Displayname: valtra.Val(request.Displayname, "Displayname").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Displayname is required."),
				valtra.MinLengthString(4, "Displayname must be at least 4 characters"),
				valtra.MaxLengthString(255, "Displayname must be at most 255 characters")).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return displayname, nil
}

func ValidateEmailRequest(request typing.ChangeEmailRequest) (typing.ChangeEmailRequest, error) {
	v := valtra.NewCollector()

	email := typing.ChangeEmailRequest{
		Email: valtra.Val(request.Email, "Email").
			Transform(valtra.TrimSpace(), valtra.Lowercase()).
			Validate(
				valtra.Required[string]("Email is required."),
				valtra.Email("Email must be in correct email format"),
			func(v valtra.Value[string]) error {
				if !setting_util.ValidateEmailMX(v.Value()) {
					return fmt.Errorf("Email domain must have valid MX or A records")
				}
				return nil
			},
			).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return email, nil
}

func ValidatePasswordRequest(request typing.ChangePasswordRequest) (typing.ChangePasswordRequest, error) {
	v := valtra.NewCollector()

	password := typing.ChangePasswordRequest{
		Password: valtra.Val(request.Password, "Password").
			Validate(
				valtra.Required[string]("Password is required"),
				valtra.MinLengthString(8, "Password must be at least 8 characters"),
				valtra.MaxLengthString(20, "Password must be at most 20 characters"),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return password, nil
}

package user_validator

import (
	"fmt"

	"github.com/bobch27/valtra-go"
	dto "resuming/controller/user/dto"
)

func ValidateUsernameRequest(request dto.ChangeUsernameRequest) (dto.ChangeUsernameRequest, error) {
	v := valtra.NewCollector()

	username := dto.ChangeUsernameRequest{
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

func ValidateDisplaynameRequest(request dto.ChangeDisplaynameRequest) (dto.ChangeDisplaynameRequest, error) {
	v := valtra.NewCollector()

	displayname := dto.ChangeDisplaynameRequest{
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

func ValidateEmailRequest(request dto.ChangeEmailRequest) (dto.ChangeEmailRequest, error) {
	v := valtra.NewCollector()

	email := dto.ChangeEmailRequest{
		Email: valtra.Val(request.Email, "Email").
			Transform(valtra.TrimSpace(), valtra.Lowercase()).
			Validate(
				valtra.Required[string]("Email is required."),
				valtra.Email("Email must be in correct email format"),
				func(v valtra.Value[string]) error {
					if !ValidateEmailMX(v.Value()) {
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

func ValidatePhoneNumberRequest(request dto.ChangePhoneNumberRequest) (dto.ChangePhoneNumberRequest, error) {
	v := valtra.NewCollector()

	phone_number := dto.ChangePhoneNumberRequest{
		PhoneNumber: valtra.Val(request.PhoneNumber, "Phone Number").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Phone number is required."),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return request, v.Errors()[0]
	}

	return phone_number, nil
}

func ValidatePasswordRequest(request dto.ChangePasswordRequest) (dto.ChangePasswordRequest, error) {
	v := valtra.NewCollector()

	password := dto.ChangePasswordRequest{
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

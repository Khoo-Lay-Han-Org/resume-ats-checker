package auth_validator

import (
	"fmt"

	"github.com/bobch27/valtra-go"
	auth_util "resuming/api/auth/util"
	typing "resuming/api/auth/typing"
)

func ValidateRegistration(request typing.Register) (typing.Register, error) {
	v := valtra.NewCollector()

	registration := typing.Register{
		Username: valtra.Val(request.Username, "Username").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Username is required."),
				valtra.MinLengthString(4, "Username must be at least 4 characters"),
				valtra.MaxLengthString(255, "Username must be at most 255 characters")).
			Collect(v),

		Displayname: valtra.Val(request.Displayname, "Displayname").
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Displayname is required."),
				valtra.MinLengthString(4, "Displayname must be at least 4 characters"),
				valtra.MaxLengthString(255, "Displayname must be at most 255 characters")).
			Collect(v),

		Email: valtra.Val(request.Email, "Email").
			Transform(valtra.TrimSpace(), valtra.Lowercase()).
			Validate(
				valtra.Required[string]("Email is required."),
				valtra.Email("Email must be in correct email format"),
			func(v valtra.Value[string]) error {
				if !auth_util.ValidateEmailMX(v.Value()) {
					return fmt.Errorf("Email domain must have valid MX or A records")
				}
				return nil
			},
			).
			Collect(v),

		Password: valtra.Val(request.Password, "Password").
			Validate(
				valtra.Required[string]("Password is required"),
				valtra.MinLengthString(8, "Password must be at least 8 characters"),
				valtra.MaxLengthString(20, "Password must be at most 20 characters"),
			).
			Collect(v),
	}

	if !v.IsValid() {
		return typing.Register{}, v.Errors()[0]
	}

	return registration, nil
}

func ValidateLogin(request typing.Login) (typing.Login, error) {
	v := valtra.NewCollector()

	login := typing.Login{
		Email: valtra.Val(request.Email).
			Transform(valtra.Lowercase(), valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Email is required."),
			).Collect(v),

		Password: valtra.Val(request.Password).
			Transform(valtra.TrimSpace()).
			Validate(
				valtra.Required[string]("Password is required."),
			).Collect(v),
	}

	if !v.IsValid() {
		return typing.Login{}, v.Errors()[0]
	}

	return login, nil
}

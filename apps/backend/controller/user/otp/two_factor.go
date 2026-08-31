package user_otp

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/systemconfig"
)

const (
	Two_factor_email = "email"
	Two_factor_sms   = "sms"
)

func SupportedTwoFactorType(c echo.Context) (string, bool) {
	factor := c.Param("2fa-type")
	switch factor {
	case Two_factor_email, Two_factor_sms:
		return factor, true
	default:
		return "", false
	}
}

func TwoFactorFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("two_factor_type")
	if err != nil {
		return "", errors.New("failed to retrieve 2FA type")
	}
	switch cookie.Value {
	case Two_factor_email, Two_factor_sms:
		return cookie.Value, nil
	default:
		return "", errors.New("unsupported 2FA type")
	}
}

func SetTwoFactorCookie(c echo.Context, factor string) {
	c.SetCookie(&http.Cookie{
		Name:     "two_factor_type",
		Value:    factor,
		MaxAge:   int(systemconfig.OtpExpiryDuration.Seconds()),
		Path:     "/",
		Domain:   "",
		Secure:   systemconfig.ApplicationHosted,
		HttpOnly: true,
	})
}

package user_api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"resuming/systemconfig"
)

const (
	two_factor_email = "email"
	two_factor_sms   = "sms"
)

func supportedTwoFactorType(c echo.Context) (string, bool) {
	factor := c.Param("2fa-type")
	switch factor {
	case two_factor_email, two_factor_sms:
		return factor, true
	default:
		return "", false
	}
}

func twoFactorFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("two_factor_type")
	if err != nil {
		return "", errors.New("failed to retrieve 2FA type")
	}
	switch cookie.Value {
	case two_factor_email, two_factor_sms:
		return cookie.Value, nil
	default:
		return "", errors.New("unsupported 2FA type")
	}
}

func setTwoFactorCookie(c echo.Context, factor string) {
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

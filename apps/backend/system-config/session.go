package systemconfig

import (
	"time"
)

var SessionExpiryDuration time.Duration
var OtpExpiryDuration time.Duration

func init() {
	SessionExpiryDuration = time.Hour * 24 * 3
	OtpExpiryDuration = time.Minute * 5
}

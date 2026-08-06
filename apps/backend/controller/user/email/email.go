package user_email

import (
	"errors"
	"log"

	"resuming/service"
)

func SendOTPEmail(email, otp string) error {
	if err := service.SendEmail(email, "Your OTP for Resuming", "Here is your OTP: "+otp, false); err != nil {
		log.Printf("Failed to send OTP: %v", err)
		return errors.New("failed to send OTP")
	}

	return nil
}

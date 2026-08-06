package user_email

import (
	"context"
	"errors"
	"log"

	"resuming/controller/user/otp"
	"resuming/service"
	systemconfig "resuming/system-config"
)

func SendEmailOTP(recipient_email string) error {
	otp_string, hashed_otp, err := user_otp.GenerateOTP()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = service.Valkey.Do(ctx, service.Valkey.B().Set().Key(recipient_email+"otp-email").Value(hashed_otp).Ex(systemconfig.OtpExpiryDuration).Build()).Error()
	if err != nil {
		log.Printf("Failed to store OTP in Valkey: %v", err)
		return errors.New("failed to process OTP")
	}

	if err := service.SendEmail(recipient_email, "Your OTP for Resuming", "Here is your OTP: "+otp_string, false); err != nil {
		log.Printf("Failed to send OTP: %v", err)
		return errors.New("failed to send OTP")
	}

	return nil
}

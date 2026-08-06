package user_sms

import (
	"context"
	"errors"
	"log"

	otp "resuming/controller/user/otp"
	"resuming/service"
	"resuming/systemconfig"
)

func SendSMSOTP(receipient_phone_number string) error {
	otp_string, hashed_otp, err := otp.GenerateOTP()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = service.Valkey.Do(ctx, service.Valkey.B().Set().Key(receipient_phone_number+"otp-sms").Value(hashed_otp).Ex(systemconfig.OtpExpiryDuration).Build()).Error()
	if err != nil {
		log.Printf("Failed to store OTP in Valkey: %v", err)
		return errors.New("failed to process OTP")
	}

	if err := service.SendSMS(receipient_phone_number, "Your OTP for Resuming", "Here is your OTP: "+otp_string); err != nil {
		log.Printf("Failed to send OTP SMS: %v", err)
		return errors.New("failed to send OTP")
	}

	return nil
}

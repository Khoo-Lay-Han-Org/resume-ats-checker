package user_otp

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"resuming/controller/user/email"
	"resuming/service"
	systemconfig "resuming/system-config"
)

func GenerateOTP() (string, string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		return "", "", errors.New("failed to process OTP")
	}

	otp := n.Int64() + 100000
	otp_string := strconv.FormatInt(otp, 10)

	hashed_otp, err := bcrypt.GenerateFromPassword([]byte(otp_string), bcrypt.MinCost)
	if err != nil {
		log.Printf("Failed to hash OTP: %v", err)
		return "", "", errors.New("failed to process OTP")
	}

	return otp_string, string(hashed_otp), nil
}

func SendOTP(recipient_email string) error {
	otp_string, hashed_otp, err := GenerateOTP()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = service.Valkey.Do(ctx, service.Valkey.B().Set().Key(recipient_email+":otp").Value(hashed_otp).Ex(systemconfig.OtpExpiryDuration).Build()).Error()
	if err != nil {
		log.Printf("Failed to store OTP in Valkey: %v", err)
		return errors.New("failed to process OTP")
	}

	if err := user_email.SendOTPEmail(recipient_email, otp_string); err != nil {
		return err
	}

	return nil
}

func CheckOTP(email, otp string) error {
	ctx := context.Background()
	value, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key(email+":otp").Build()).ToString()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(value), []byte(otp))
	if err != nil {
		return err
	}

	err = service.Valkey.Do(ctx, service.Valkey.B().Del().Key(email+":otp").Build()).Error()
	if err != nil {
		return err
	}

	return nil
}

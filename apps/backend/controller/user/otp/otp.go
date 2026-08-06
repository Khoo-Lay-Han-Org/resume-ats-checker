package user_otp

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"resuming/service"
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

func CheckEmailOTP(email, otp string) error {
	return checkStoredOTP(email+"otp-email", otp)
}

func CheckSMSOTP(phone_number, otp string) error {
	return checkStoredOTP(phone_number+"otp-sms", otp)
}

func checkStoredOTP(key, otp string) error {
	ctx := context.Background()
	value, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key(key).Build()).ToString()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(value), []byte(otp))
	if err != nil {
		return err
	}

	err = service.Valkey.Do(ctx, service.Valkey.B().Del().Key(key).Build()).Error()
	if err != nil {
		return err
	}

	return nil
}

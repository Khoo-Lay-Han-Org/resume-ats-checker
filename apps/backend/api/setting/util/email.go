package setting_util

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func SendOTP(email string) error {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		return errors.New("failed to process OTP")
	}
	otp := n.Int64() + 100000

	otp_string := strconv.FormatInt(otp, 10)

	hashed_otp, err := bcrypt.GenerateFromPassword([]byte(otp_string), bcrypt.MinCost)
	if err != nil {
		log.Printf("Failed to hash OTP: %v", err)
		return errors.New("failed to process OTP")
	}

	ctx := context.Background()
	err = tool.Valkey.Do(ctx, tool.Valkey.B().Set().Key(email+":otp").Value(string(hashed_otp)).Ex(systemconfig.OtpExpiryDuration).Build()).Error()

	if err != nil {
		log.Printf("Failed to store OTP in Valkey: %v", err)
		return errors.New("failed to process OTP")
	}

	if err := tool.SendEmail(email, "Your OTP for Resuming", "Here is your OTP: "+otp_string, false); err != nil {
		log.Printf("Failed to send OTP: %v", err)
		return errors.New("failed to send OTP")
	}

	return nil
}

func CheckOTP(email, otp string) error {
	ctx := context.Background()
	value, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(email+":otp").Build()).ToString()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(value), []byte(otp))
	if err != nil {
		return err
	}

	err = tool.Valkey.Do(ctx, tool.Valkey.B().Del().Key(email+":otp").Build()).Error()
	if err != nil {
		return err
	}

	return nil
}

func ValidateEmailMX(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	resolver := net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mx, err := resolver.LookupMX(ctx, domain)
	if err == nil && len(mx) > 0 {
		return true
	}

	ips, err := resolver.LookupIPAddr(ctx, domain)
	if err == nil && len(ips) > 0 {
		return true
	}

	return false
}

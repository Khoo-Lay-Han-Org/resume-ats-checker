package user_oauth

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandomPassword() ([]byte, error) {
	password := make([]byte, 20)

	_, err := rand.Read(password)
	if err != nil {
		return nil, err
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashed_password, nil
}

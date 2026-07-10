package auth_util

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

// the passwrod field is required, but the user login using oauth, therefore, randomly generate one for them
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

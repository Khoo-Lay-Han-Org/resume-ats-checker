package auth_util

import (
	"errors"

	"resuming/database"
)

func FindUser(private_id int) (*database.User, error) {
	var user database.User
	result := database.DB.Where("id = ?", private_id).First(&user)
	if result.Error != nil {
		return nil, errors.New("failed to find user")
	}
	return &user, nil
}

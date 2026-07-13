package database

import (
	"context"
	"encoding/json"
	"fmt"

	valkey "github.com/valkey-io/valkey-go"
	"resuming/tool"
)

func FindUser(private_id int) (*User, error) {
	var user User
	result := DB.Preload("ShowcaseRecord").
		Preload("Session").
		Preload("Resume").
		Preload("Portfolio").
		Preload("ATS").
		Preload("JWTKey").
		Scopes(NonExpiredUser).
		Where("id = ?", private_id).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}
	return &user, nil
}

func FindUserByPublicId(publicId string) (*User, error) {
	var user User
	result := DB.Preload("ShowcaseRecord").
		Preload("Session").
		Preload("Resume").
		Preload("Portfolio").
		Preload("ATS").
		Scopes(NonExpiredUser).
		Where("public_id = ?", publicId).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user by public_id: %w", result.Error)
	}
	return &user, nil
}

func resolveUserId(public_user_id string) (int, error) {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read user data for ID resolution: %w", err)
	}
	var user User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return 0, fmt.Errorf("failed to parse user data for ID resolution: %w", err)
	}
	var dbUser User
	if err := DB.Where("public_id = ?", user.PublicId).First(&dbUser).Error; err != nil {
		return 0, fmt.Errorf("failed to find user by public_id: %w", err)
	}
	return dbUser.Id, nil
}

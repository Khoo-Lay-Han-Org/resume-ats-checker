package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	valkey "github.com/valkey-io/valkey-go"
	"resuming/database/sqlc"
	"resuming/tool"
)

func FindUser(private_id int32) (*sqlc.User, error) {
	user, err := Queries.FindUserById(context.Background(), private_id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func FindUserByPublicId(publicId string) (*sqlc.User, error) {
	uid := pgtype.UUID{}
	if err := uid.Scan(publicId); err != nil {
		return nil, fmt.Errorf("invalid public_id: %w", err)
	}
	user, err := Queries.FindUserByPublicId(context.Background(), uid)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by public_id: %w", err)
	}
	return &user, nil
}

func resolveUserId(public_user_id string) (int32, error) {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read user data for ID resolution: %w", err)
	}
	var user sqlc.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return 0, fmt.Errorf("failed to parse user data for ID resolution: %w", err)
	}
	uid := pgtype.UUID{}
	if err := uid.Scan(user.PublicID); err != nil {
		return 0, fmt.Errorf("invalid public_id in cache: %w", err)
	}
	dbUser, err := Queries.FindUserByPublicId(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("failed to find user by public_id: %w", err)
	}
	return dbUser.ID, nil
}

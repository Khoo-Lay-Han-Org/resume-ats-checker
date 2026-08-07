package shared_find

import (
	"context"
	"encoding/json"

	valkey "github.com/valkey-io/valkey-go"
	"resuming/database"
	"resuming/database/sqlc"
	"resuming/service"
)

func FindUser(private_id int32) (*sqlc.User, error) {
	return database.FindUser(private_id)
}

func GetUser(public_user_id string) (*sqlc.User, error) {
	ctx := context.Background()
	retrieved_data, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
	if err == nil {
		var user sqlc.User
		if err := json.Unmarshal([]byte(retrieved_data), &user); err != nil {
			return nil, err
		}
		return &user, nil
	}

	if !valkey.IsValkeyNil(err) {
		return nil, err
	}

	user, dbErr := database.FindUserByPublicId(public_user_id)
	if dbErr != nil {
		return nil, dbErr
	}

	if syncErr := database.SyncIndividualUserDataSessionStore(public_user_id, user); syncErr != nil {
		return nil, syncErr
	}

	return user, nil
}

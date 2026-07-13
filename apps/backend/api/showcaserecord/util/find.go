package showcaserecord_util

import (
	"context"
	"encoding/json"

	"resuming/database"
	"resuming/tool"
)

func FindUser(public_user_id string) (*database.User, error) {
	ctx := context.Background()
	retrieved_data, err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Get().
			Key(public_user_id+":user_data").
			Build(),
	).ToString()
	if err != nil {
		return nil, err
	}

	var user database.User
	err = json.Unmarshal([]byte(retrieved_data), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

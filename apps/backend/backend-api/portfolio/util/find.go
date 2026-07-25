package portfolio_util

import (
	"context"
	"encoding/json"

	"resuming/database/sqlc"
	"resuming/tool"
)

func FindUser(public_user_id string) (*sqlc.User, error) {
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

	var user sqlc.User
	err = json.Unmarshal([]byte(retrieved_data), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

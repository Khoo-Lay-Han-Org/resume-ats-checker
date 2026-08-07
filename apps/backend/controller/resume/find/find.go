package resume_find

import (
	"context"

	valkey "github.com/valkey-io/valkey-go"
	"resuming/database"
	"resuming/service"
)

func GetResumeData(public_user_id string) (string, error) {
	ctx := context.Background()
	retrieved_data, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":resume_data").Build()).ToString()
	if err == nil {
		return retrieved_data, nil
	}

	if !valkey.IsValkeyNil(err) {
		return "", err
	}

	user, dbErr := database.FindUserByPublicId(public_user_id)
	if dbErr != nil {
		return "", dbErr
	}

	resume, dbErr := database.Queries.FindResumeByUserId(ctx, user.ID)
	if dbErr != nil {
		return "", dbErr
	}

	if syncErr := database.SyncIndividualResumeDataSessionStore(public_user_id, &resume); syncErr != nil {
		return "", syncErr
	}

	retrieved_data, err = service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":resume_data").Build()).ToString()
	if err != nil {
		return "", err
	}

	return retrieved_data, nil
}

package showcase_find

import (
	"context"

	valkey "github.com/valkey-io/valkey-go"
	"resuming/database"
	"resuming/service"
)

func GetShowcaseRecordData(public_user_id string) (string, error) {
	ctx := context.Background()
	retrieved_data, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
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

	showcase, dbErr := database.Queries.FindShowcaseRecordByUserId(ctx, user.ID)
	if dbErr != nil {
		return "", dbErr
	}

	if syncErr := database.SyncIndividualShowCaseRecordDataSessionStore(public_user_id, &showcase); syncErr != nil {
		return "", syncErr
	}

	retrieved_data, err = service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
	if err != nil {
		return "", err
	}

	return retrieved_data, nil
}

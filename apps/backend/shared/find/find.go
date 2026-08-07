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

func GetPortfolioData(public_user_id string) (string, error) {
	ctx := context.Background()
	retrieved_data, err := service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":portfolio_data").Build()).ToString()
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

	portfolio, dbErr := database.Queries.FindPortfolioByUserId(ctx, user.ID)
	if dbErr != nil {
		return "", dbErr
	}

	if syncErr := database.SyncIndividualPortfolioDataSessionStore(public_user_id, &portfolio); syncErr != nil {
		return "", syncErr
	}

	retrieved_data, err = service.Valkey.Do(ctx, service.Valkey.B().Get().Key(public_user_id+":portfolio_data").Build()).ToString()
	if err != nil {
		return "", err
	}

	return retrieved_data, nil
}

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

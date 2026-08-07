package portfolio_find

import (
	"context"

	valkey "github.com/valkey-io/valkey-go"
	"resuming/database"
	"resuming/service"
)

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

package user_api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	typing "resuming/controller/user/dto"
	validator "resuming/controller/user/validator"
	"resuming/database"
	"resuming/database/sqlc"
	"resuming/service"
	systemconfig "resuming/system-config"
)

func BanClient() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request typing.UserControlRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve request."})
		}

		polished_request, err := validator.ValidateUserControlRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		public_user_id := polished_request.PublicUserId

		ctx := c.Request().Context()
		group_data, err := service.Valkey.Do(
			ctx,
			service.Valkey.B().Get().Key("user_data").Build(),
		).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve cached data."})
		}

		var all_users []sqlc.User
		if err := json.Unmarshal([]byte(group_data), &all_users); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse cached data."})
		}

		var target_user *sqlc.User
		now := time.Now()
		for i, user := range all_users {
			if user.PublicID.String() == public_user_id {
				all_users[i].BannedAt = pgtype.Timestamptz{Time: now, Valid: true}
				target_user = &all_users[i]
				break
			}
		}

		if target_user == nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found."})
		}

		serialised_group, err := json.Marshal(all_users)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process updated data."})
		}
		if err := service.Valkey.Do(
			ctx,
			service.Valkey.B().Set().
				Key("user_data").Value(string(serialised_group)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update user data."})
		}

		individual_data, err := json.Marshal(target_user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to serialise user data"})
		}
		if err := service.Valkey.Do(
			ctx,
			service.Valkey.B().Set().
				Key(public_user_id+":user_data").
				Value(string(individual_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to ban user."})
		}

		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id)

		return nil
	}
}

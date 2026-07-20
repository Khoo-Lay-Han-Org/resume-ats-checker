package administrator_view

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	typing "resuming/api/administrator/typing"
	validator "resuming/api/administrator/validator"
	"resuming/database"
	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
	"resuming/tool"
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
		group_data, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("user_data").Build(),
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
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
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
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
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

func RemoveIndividualUserSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		var request typing.SessionControlRequest
		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Failed to retrieve request."})
		}

		polished_request, err := validator.ValidateSessionControlRequest(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		}

		public_user_id := polished_request.PublicUserId

		admin_session_id := c.Get("public_user_id")
		if admin_session_id != nil && public_user_id == admin_session_id.(string) {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Cannot remove your own session."})
		}

		ctx := c.Request().Context()
		exists, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Exists().Key(public_user_id+":session_data").Build(),
		).AsInt64()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to verify session."})
		}
		if exists == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Session not found."})
		}

		err = tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Del().
				Key(public_user_id+":session_data",
					public_user_id+":jwt_data",
					public_user_id+":user_data").
				Build(),
		).Error()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete session"})
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "Successfully deleted session"})
	}
}

func RemoveAllClientSession() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_configs").Build()).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to find client configs."})
		}

		var data []map[string]any
		err = json.Unmarshal([]byte(retrieved_data), &data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse configs."})
		}

		for _, item := range data {
			public_user_id, ok := item["public_user_id"].(string)
			if !ok {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Invalid client config format."})
			}

			ctx := c.Request().Context()
			_, err := tool.Valkey.Do(
				ctx,
				tool.Valkey.B().Del().
					Key(public_user_id+":session_data").
					Build(),
			).ToString()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to delete session"})
			}
		}

		admin_id := c.Get("public_user_id")
		if admin_id != nil {
			admin_key := admin_id.(string) + ":session_data"
			ttl_seconds, err := tool.Valkey.Do(ctx, tool.Valkey.B().Ttl().Key(admin_key).Build()).AsInt64()
			if err == nil && ttl_seconds > 300 {
				tool.Valkey.Do(ctx, tool.Valkey.B().Expire().Key(admin_key).Seconds(300).Build())
			}
		}

		return nil
	}
}

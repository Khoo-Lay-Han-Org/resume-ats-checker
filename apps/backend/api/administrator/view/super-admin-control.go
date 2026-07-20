package administrator_view

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	typing "resuming/api/administrator/typing"
	util "resuming/api/administrator/util"
	validator "resuming/api/administrator/validator"
	"resuming/database"
	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func RemoveAdmin() echo.HandlerFunc {
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
		for i, user := range all_users {
			if user.PublicID.String() == public_user_id {
				all_users[i].UserType = sqlc.UserTypeClient
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
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to store user data"})
		}

		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id)

		return nil
	}
}

func ChangeAdminAccessibility() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func InvitationToBecomeAdmin() echo.HandlerFunc {
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
		retrieved_data, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("user_data").Build(),
		).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve cached data."})
		}

		var all_users []sqlc.User
		if err := json.Unmarshal([]byte(retrieved_data), &all_users); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse cached data."})
		}

		var target_user sqlc.User
		found := false
		for _, user := range all_users {
			if user.PublicID.String() == public_user_id {
				target_user = user
				found = true
				break
			}
		}

		if !found {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found."})
		}

		token := uuid.New().String()

		print(token)

		invite_data := database.InviteTokenDTO{
			PublicUserId: public_user_id,
			Used:         false,
		}
		serialised_invite, err := json.Marshal(invite_data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to create invite token."})
		}

		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("invite_token:"+token).
				Value(string(serialised_invite)).
				Ex(48*time.Hour).
				Build(),
		).Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to store invite token."})
		}

		if err := util.EmailInvitationToBecomeAdmin(target_user.Email, token); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send email invitation."})
		}

		return nil
	}
}

func AcceptToBecomeAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Param("token")
		if token == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Missing invite token."})
		}

		ctx := c.Request().Context()
		invite_raw, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("invite_token:"+token).Build(),
		).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Invalid or expired invite token."})
		}

		var invite_data database.InviteTokenDTO
		if err := json.Unmarshal([]byte(invite_raw), &invite_data); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse invite token."})
		}

		if invite_data.Used {
			return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invite token has already been used."})
		}

		invite_data.Used = true
		serialised_invite, err := json.Marshal(invite_data)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update invite token."})
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("invite_token:"+token).Value(string(serialised_invite)).
				Ex(48*time.Hour).
				Build(),
		).Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to mark invite as used."})
		}

		user_data_raw, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("user_data").Build(),
		).ToString()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "Failed to retrieve cached data."})
		}

		var all_users []sqlc.User
		if err := json.Unmarshal([]byte(user_data_raw), &all_users); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to parse cached data."})
		}

		found := false
		for i, user := range all_users {
			if user.PublicID.String() == invite_data.PublicUserId {
				all_users[i].UserType = sqlc.UserTypeAdmin
				found = true
				break
			}
		}

		if !found {
			return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found."})
		}

		serialised_users, err := json.Marshal(all_users)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to process updated data."})
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("user_data").Value(string(serialised_users)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update user data."})
		}

		return nil
	}
}

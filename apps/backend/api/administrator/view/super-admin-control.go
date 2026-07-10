package administrator_view

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	typing "resuming/api/administrator/typing"
	util "resuming/api/administrator/util"
	validator "resuming/api/administrator/validator"
	"resuming/database"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func RemoveAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request typing.UserControlRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve request."})
			return
		}

		polished_request, err := validator.ValidateUserControlRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		public_user_id := polished_request.PublicUserId

		ctx := c.Request.Context()
		group_data, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("user_data").Build(),
		).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve cached data."})
			return
		}

		var all_users []database.User
		if err := json.Unmarshal([]byte(group_data), &all_users); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse cached data."})
			return
		}

		var target_user *database.User
		for i, user := range all_users {
			if user.PublicId.String() == public_user_id {
				all_users[i].UserType = database.Client
				target_user = &all_users[i]
				break
			}
		}

		if target_user == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}

		serialised_group, err := json.Marshal(all_users)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process updated data."})
			return
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("user_data").Value(string(serialised_group)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user data."})
			return
		}

		individual_data, err := json.Marshal(target_user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to serialise user data"})
			return
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(public_user_id+":user_data").
				Value(string(individual_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store user data"})
			return
		}

		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id)
	}
}

func ChangeAdminAccessibility() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func InvitationToBecomeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request typing.UserControlRequest
		if err := c.BindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve request."})
			return
		}

		polished_request, err := validator.ValidateUserControlRequest(request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		public_user_id := polished_request.PublicUserId

		ctx := c.Request.Context()
		retrieved_data, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("user_data").Build(),
		).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve cached data."})
			return
		}

		var all_users []database.User
		if err := json.Unmarshal([]byte(retrieved_data), &all_users); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse cached data."})
			return
		}

		var target_user database.User
		found := false
		for _, user := range all_users {
			if user.PublicId.String() == public_user_id {
				target_user = user
				found = true
				break
			}
		}

		if !found {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}

		token := uuid.New().String()

		print(token)

		invite_data := database.InviteTokenDTO{
			PublicUserId: public_user_id,
			Used:         false,
		}
		serialised_invite, err := json.Marshal(invite_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create invite token."})
			return
		}

		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("invite_token:"+token).
				Value(string(serialised_invite)).
				Ex(48*time.Hour).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store invite token."})
			return
		}

		if err := util.EmailInvitationToBecomeAdmin(target_user.Email, token); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email invitation."})
			return
		}

	}
}

func AcceptToBecomeAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Missing invite token."})
			return
		}

		ctx := c.Request.Context()
		invite_raw, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("invite_token:"+token).Build(),
		).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid or expired invite token."})
			return
		}

		var invite_data database.InviteTokenDTO
		if err := json.Unmarshal([]byte(invite_raw), &invite_data); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse invite token."})
			return
		}

		if invite_data.Used {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invite token has already been used."})
			return
		}

		invite_data.Used = true
		serialised_invite, err := json.Marshal(invite_data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update invite token."})
			return
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("invite_token:"+token).Value(string(serialised_invite)).
				Ex(48*time.Hour).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to mark invite as used."})
			return
		}

		user_data_raw, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Get().Key("user_data").Build(),
		).ToString()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Failed to retrieve cached data."})
			return
		}

		var all_users []database.User
		if err := json.Unmarshal([]byte(user_data_raw), &all_users); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse cached data."})
			return
		}

		found := false
		for i, user := range all_users {
			if user.PublicId.String() == invite_data.PublicUserId {
				all_users[i].UserType = "admin"
				found = true
				break
			}
		}

		if !found {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found."})
			return
		}

		serialised_users, err := json.Marshal(all_users)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to process updated data."})
			return
		}
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key("user_data").Value(string(serialised_users)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user data."})
			return
		}

	}
}

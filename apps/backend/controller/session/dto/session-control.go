package session_dto

type SessionControlRequest struct {
	PublicUserId string `json:"public_user_id" binding:"required"`
}

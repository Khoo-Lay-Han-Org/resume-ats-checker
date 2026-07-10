package administrator_typing

type UserControlRequest struct {
	PublicUserId string `json:"public_user_id" binding:"required"`
}

type SessionControlRequest struct {
	PublicUserId string `json:"public_user_id" binding:"required"`
}

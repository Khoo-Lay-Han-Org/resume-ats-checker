package user_dto

type UserControlRequest struct {
	PublicUserId string `json:"public_user_id" binding:"required" mod:"trim" validate:"required"`
}

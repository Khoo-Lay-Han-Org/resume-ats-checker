package clientsupportmessage_dto

type ClientCommunicationReplyRequest struct {
	PublicId string `json:"public_id" binding:"required" mod:"trim" validate:"required"`
	Message  string `json:"message" binding:"required" mod:"trim" validate:"required"`
}

type ClientCommunicateRequest struct {
	Type    string `json:"type" binding:"required" mod:"trim" validate:"required"`
	Message string `json:"message" binding:"required" mod:"trim" validate:"required"`
}

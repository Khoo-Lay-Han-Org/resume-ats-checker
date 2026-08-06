package clientsupportmessage_dto

type ClientCommunicationReplyRequest struct {
	PublicId string `json:"public_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

type ClientCommunicateRequest struct {
	Type    string `json:"type" binding:"required"`
	Message string `json:"message" binding:"required"`
}

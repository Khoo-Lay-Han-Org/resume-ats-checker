package administrator_typing

type ClientCommunicationReplyRequest struct {
	PublicId string `json:"public_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

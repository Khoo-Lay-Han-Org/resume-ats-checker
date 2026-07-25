package client_support_dto

type ClientCommunicateRequest struct {
	Type    string `json:"type" binding:"required"`
	Message string `json:"message" binding:"required"`
}

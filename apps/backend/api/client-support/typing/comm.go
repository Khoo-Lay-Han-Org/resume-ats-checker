package client_support_typing

type ClientCommunicateRequest struct {
	Type    string `json:"type" binding:"required"`
	Message string `json:"message" binding:"required"`
}

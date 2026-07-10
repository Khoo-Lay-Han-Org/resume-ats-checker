package portfolio_typing

type ChooseTemplateRequest struct {
	TemplateId string `json:"template_id" binding:"required"`
}

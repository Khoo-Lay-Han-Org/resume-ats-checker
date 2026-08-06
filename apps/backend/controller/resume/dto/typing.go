package resume_dto

type ChooseTemplateRequest struct {
	TemplateId string `json:"template_id" binding:"required"`
}

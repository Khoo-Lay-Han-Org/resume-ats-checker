package api

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Displayname string `json:"displayname" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// OTPRequest represents the OTP verification request body
type OTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}

// ChooseResumeRequest represents the choose resume template request body
type ChooseResumeRequest struct {
	TemplateId string `json:"template_id" binding:"required"`
}

// ResumeDataResponse represents the resume data response
type ResumeDataResponse struct {
	TemplateId    int         `json:"template_id"`
	Name          []string    `json:"name"`
	Email         []string    `json:"email"`
	PhoneNumber   []string    `json:"phone_number"`
	Address       []string    `json:"address"`
	SocialMedia   []string    `json:"social_media"`
	JobExperience interface{} `json:"job_experience"`
	Education     interface{} `json:"education"`
	Skill         []string    `json:"skill"`
	Certificate   interface{} `json:"certificate"`
	Language      []string    `json:"language"`
	Project       interface{} `json:"project"`
}

// ResumeResponse represents the resume response
type ResumeResponse struct {
	Message string             `json:"message"`
	Data    ResumeDataResponse `json:"data"`
}

// ChoosePortfolioRequest represents the choose portfolio template request body
type ChoosePortfolioRequest struct {
	TemplateId string `json:"template_id" binding:"required"`
}

// ShowCaseRecordDeleteRequest represents the delete showcase record request body
type ShowCaseRecordDeleteRequest struct {
	SectionTitle string `json:"sectiontitle" binding:"required"`
	Index        string `json:"index" binding:"required"`
}

// ATSScoreResponse represents the ATS score response
type ATSScoreResponse struct {
	Message string `json:"message"`
	Data    int    `json:"data"`
}

// PortfolioDataResponse represents the portfolio data response
type PortfolioDataResponse struct {
	TemplateId    int         `json:"template_id"`
	Name          []string    `json:"name"`
	Email         []string    `json:"email"`
	PhoneNumber   []string    `json:"phone_number"`
	Address       []string    `json:"address"`
	SocialMedia   []string    `json:"social_media"`
	JobExperience interface{} `json:"job_experience"`
	Education     interface{} `json:"education"`
	Skill         []string    `json:"skill"`
	Certificate   interface{} `json:"certificate"`
	Language      []string    `json:"language"`
	Project       interface{} `json:"project"`
}

// PortfolioResponse represents the portfolio response
type PortfolioResponse struct {
	Message string                `json:"message"`
	Data    PortfolioDataResponse `json:"data"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// SelectedUser represents a user selection for admin invitation
type SelectedUser struct {
	PublicId string `json:"public_id" binding:"required"`
}

// ChangeUsernameRequest represents the change username request body
type ChangeUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

// ChangeDisplaynameRequest represents the change displayname request body
type ChangeDisplaynameRequest struct {
	Displayname string `json:"displayname" binding:"required"`
}

// ChangeEmailRequest represents the change email request body
type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

// ClientReportRequest represents the client report request body
type ClientReportRequest struct {
	TargetClientPublicUserId string `json:"target client_public_user_id" binding:"required"`
	ReportType               string `json:"report type" binding:"required"`
}

// ClientCommunicateRequest represents the client communication request body
type ClientCommunicateRequest struct {
	Type    string `json:"type" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// ClientCommunicationReplyRequest represents the admin reply to client communication
type ClientCommunicationReplyRequest struct {
	PublicId string `json:"public_id" binding:"required"`
	Message  string `json:"message" binding:"required"`
}

// UserControlRequest represents the user control request body for admin operations
type UserControlRequest struct {
	PublicUserId string `json:"public_user_id" binding:"required"`
}

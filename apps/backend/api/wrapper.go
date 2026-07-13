package api

// @title           Resume Builder API
// @version         1.0
// @description     API for resume building, portfolio management, and ATS scoring
// @host            localhost:5321
// @BasePath        /

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	administrator_view "resuming/api/administrator/view"
	ats_view "resuming/api/ats/view"
	auth_view "resuming/api/auth/view"
	client_support_view "resuming/api/client-support/view"
	portfolio_view "resuming/api/portfolio/view"
	resume_view "resuming/api/resume/view"
	setting_view "resuming/api/setting/view"
	showcaserecord_view "resuming/api/showcaserecord/view"
	"resuming/database"
)

//// AUTH

// PrepareRegistrationFlow initiates the registration process by sending an OTP
// @Summary Prepare registration
// @Description Validates registration data and sends OTP to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 200 {object} SuccessResponse
// @Failure 400,422,500 {object} ErrorResponse
// @Router /prepare-registeration [post]
func PrepareRegistrationFlow(c *gin.Context) {
	auth_view.PrepareRegistration()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully."})
}

// RegisterFlow completes user registration after OTP verification
// @Summary Complete registration
// @Description Verifies OTP and creates user account (path param: type-of-user = "client" or "admin")
// @Tags Authentication
// @Accept json
// @Produce json
// @Param type-of-user path string true "User type (client or admin)"
// @Param request body OTPRequest true "OTP verification"
// @Success 200 {object} SuccessResponse
// @Failure 400,404,422,500 {object} ErrorResponse
// @Router /register/{type-of-user} [post]
func RegisterFlow(c *gin.Context) {
	auth_view.Register()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created user."})
}

// PrepareLoginFlow initiates the login process by sending an OTP
// @Summary Prepare login
// @Description Validates login credentials and sends OTP to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} SuccessResponse
// @Failure 400,404,422,500 {object} ErrorResponse
// @Router /prepare-login [post]
func PrepareLoginFlow(c *gin.Context) {
	auth_view.PrepareLogin()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully."})
}

// LoginFlow completes user login after OTP verification and syncs all user data
// @Summary Complete login
// @Description Verifies OTP, creates session, and syncs user data
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body OTPRequest true "OTP verification"
// @Success 200 {object} SuccessResponse
// @Failure 400,404,422,500 {object} ErrorResponse
// @Router /login [post]
func LoginFlow(c *gin.Context) {
	auth_view.Login()(c)
	if c.IsAborted() {
		return
	}

	auth_view.SetSession()(c)
	if c.IsAborted() {
		return
	}

	if private_id, exists := c.Get("private_id"); exists {
		if public_user_id, exists2 := c.Get("public_user_id"); exists2 {
			if session_key, exists3 := c.Get("session_key"); exists3 {
				if signing_key, exists4 := c.Get("signing_key"); exists4 {
					if user, exists5 := c.Get("user"); exists5 {
						if err := database.SyncIndividualLoginDataToSessionStore(
							public_user_id.(string),
							session_key.(string),
							signing_key.(string),
							private_id.(int),
							user.(*database.User),
						); err != nil {
							log.Printf("Failed to sync login data: %v", err)
						}
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in."})
}

//// SHOWCASERECORD

// ShowCaseRecordAddFlow adds a new showcase record of the specified type
// @Summary Add showcase record
// @Description Adds a new showcase record (name, email, phone-number, address, social-media, job-experience, education, skill, certificate, language, project)
// @Tags ShowCaseRecords
// @Accept json
// @Produce json
// @Param type-of-data path string true "Type of data (name, email, phone-number, address, social-media, job-experience, education, skill, certificate, language, project)"
// @Param request body interface{} true "Showcase record data (structure depends on type-of-data)"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /showcaserecord-add/{type-of-data} [post]
func ShowCaseRecordAddFlow(c *gin.Context) {
	showcaserecord_view.AddShowCaseRecordData()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully added showcase record."})
}

// ShowCaseRecordDeleteFlow deletes a showcase record by section and index
// @Summary Delete showcase record
// @Description Deletes a showcase record at a specific index within a section
// @Tags ShowCaseRecords
// @Accept json
// @Produce json
// @Param request body ShowCaseRecordDeleteRequest true "Delete request"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,422,500 {object} ErrorResponse
// @Router /showcaserecord-delete [delete]
func ShowCaseRecordDeleteFlow(c *gin.Context) {
	showcaserecord_view.DeleteShowCaseRecordData()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted showcase record."})
}

// ShowCaseRecordEditFlow edits an existing showcase record of the specified type
// @Summary Edit showcase record
// @Description Edits a showcase record at a specific index (name, email, phone-number, address, social-media, job-experience, education, skill, certificate, language, project)
// @Tags ShowCaseRecords
// @Accept json
// @Produce json
// @Param type-of-data path string true "Type of data (name, email, phone-number, address, social-media, job-experience, education, skill, certificate, language, project)"
// @Param request body interface{} true "Updated showcase record data (structure depends on type-of-data)"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /showcaserecord-edit/{type-of-data} [patch]
func ShowCaseRecordEditFlow(c *gin.Context) {
	showcaserecord_view.EditShowCaseRecordData()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully edited showcase record."})
}

// ShowCaseRecordGetFlow retrieves all showcase records for the current user
// @Summary Get showcase records
// @Description Retrieves all showcase records for the authenticated user
// @Tags ShowCaseRecords
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401,404,500 {object} ErrorResponse
// @Router /showcaserecord-retrieve [get]
func ShowCaseRecordGetFlow(c *gin.Context) {
	showcaserecord_view.RetrieveShowCaseRecordData()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved showcase records.", "data": data})
}

//// PORTFOLIO

// ChoosePortfolioTemplateFlow selects a portfolio template
// @Summary Choose portfolio template
// @Description Sets the portfolio template for the user
// @Tags Portfolio
// @Accept json
// @Produce json
// @Param request body ChoosePortfolioRequest true "Template selection"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /choose-portfolio-template [patch]
func ChoosePortfolioTemplateFlow(c *gin.Context) {
	portfolio_view.ChooseTemplate()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed the template."})
}

// GetPortfolioContentFlow retrieves the portfolio content
// @Summary Get portfolio content
// @Description Retrieves the user's portfolio data including all sections
// @Tags Portfolio
// @Produce json
// @Success 200 {object} PortfolioResponse
// @Failure 401,404,500 {object} ErrorResponse
// @Router /get-portfolio-content [get]
func GetPortfolioContentFlow(c *gin.Context) {
	portfolio_view.RetrievePortfolioData()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved portfolio data.", "data": data})
}

//// RESUME

// ChooseResumeTemplateFlow selects a resume template
// @Summary Choose resume template
// @Description Sets the resume template for the user
// @Tags Resume
// @Accept json
// @Produce json
// @Param request body ChooseResumeRequest true "Template selection"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /choose-resume-template [patch]
func ChooseResumeTemplateFlow(c *gin.Context) {
	resume_view.ChooseTemplate()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed the template."})
}

// GetResumeContentFlow retrieves the resume content
// @Summary Get resume content
// @Description Retrieves the user's resume data including all sections
// @Tags Resume
// @Produce json
// @Success 200 {object} ResumeResponse
// @Failure 401,404,422,500 {object} ErrorResponse
// @Router /get-resume-content [get]
func GetResumeContentFlow(c *gin.Context) {
	resume_view.RetrieveResumeData()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved resume data.", "data": data})
}

//// ATS

// ATSScoreWebScrapeFlow calculates ATS score by web scraping job description
// @Summary ATS score with web scrape
// @Description Extracts resume, validates it, web scrapes job description from company/title, and calculates ATS score
// @Tags ATS
// @Accept multipart/form-data
// @Produce json
// @Param resume_file formData file true "Resume PDF file"
// @Param company formData string true "Company name for job description"
// @Param job_title formData string true "Job title for job description"
// @Success 200 {object} ATSScoreResponse
// @Failure 400,401,422,500 {object} ErrorResponse
// @Router /ats-score-webscrape [post]
func ATSScoreWebScrapeFlow(c *gin.Context) {
	ats_view.ExtractResume()(c)
	if c.IsAborted() {
		return
	}

	ats_view.ParseResume()(c)
	if c.IsAborted() {
		return
	}

	ats_view.SectionExistenceCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.FormattingCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.WebScrapeJobDesc()(c)
	if c.IsAborted() {
		return
	}

	ats_view.ResumeJobTypeCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.JobDescJobTypeCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.JobTypeRelevanceCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.ResumeSkillsCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.JobDescSkillsCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.OverallSkillsCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.OverallScore()(c)
	if c.IsAborted() {
		return
	}

	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully calculated ATS score.", "data": data})
}

// ATSScoreUserInputFlow calculates ATS score using user-provided job description
// @Summary ATS score with user input
// @Description Extracts resume, validates it, uses user-provided job description, and calculates ATS score
// @Tags ATS
// @Accept multipart/form-data
// @Produce json
// @Param resume_file formData file true "Resume PDF file"
// @Param job_desc formData string true "Job description text"
// @Success 200 {object} ATSScoreResponse
// @Failure 400,401,422,500 {object} ErrorResponse
// @Router /ats-score-user-input [post]
func ATSScoreUserInputFlow(c *gin.Context) {
	ats_view.ExtractResume()(c)
	if c.IsAborted() {
		return
	}

	ats_view.ParseResume()(c)
	if c.IsAborted() {
		return
	}

	ats_view.SectionExistenceCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.FormattingCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.UserInputJobDesc()(c)
	if c.IsAborted() {
		return
	}

	ats_view.ResumeJobTypeCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.JobDescJobTypeCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.JobTypeRelevanceCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.ResumeSkillsCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.JobDescSkillsCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.OverallSkillsCheck()(c)
	if c.IsAborted() {
		return
	}

	ats_view.OverallScore()(c)
	if c.IsAborted() {
		return
	}

	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully calculated ATS score.", "data": data})
}

//// SETTING

// ChangeUsernameFlow updates the user's username
// @Summary Change username
// @Description Updates the authenticated user's username
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body ChangeUsernameRequest true "New username"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /change-username [post]
func ChangeUsernameFlow(c *gin.Context) {
	setting_view.ChangeUsername()(c)
	if c.IsAborted() {
		return
	}
	if public_user_id, exists := c.Get("public_user_id"); exists {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed your username."})
}

// ChangeDisplaynameFlow updates the user's displayname
// @Summary Change displayname
// @Description Updates the authenticated user's displayname
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body ChangeDisplaynameRequest true "New displayname"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /change-displayname [post]
func ChangeDisplaynameFlow(c *gin.Context) {
	setting_view.ChangeDisplayname()(c)
	if c.IsAborted() {
		return
	}
	if public_user_id, exists := c.Get("public_user_id"); exists {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed your display name."})
}

// PrepareEmailChangeFlow initiates the email change process by sending OTP
// @Summary Prepare email change
// @Description Validates new email and sends OTP to current email for verification
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body ChangeEmailRequest true "New email"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,409,422,500 {object} ErrorResponse
// @Router /prepare-change-email [post]
func PrepareEmailChangeFlow(c *gin.Context) {
	setting_view.PrepareChangeEmail()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email."})
}

// ChangeEmailFlow completes the email change after OTP verification
// @Summary Change email
// @Description Verifies OTP and updates user's email address
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body OTPRequest true "OTP verification"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /change-email [post]
func ChangeEmailFlow(c *gin.Context) {
	setting_view.ChangeEmail()(c)
	if c.IsAborted() {
		return
	}
	if public_user_id, exists := c.Get("public_user_id"); exists {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed your email."})
}

// PreparePasswordChangeFlow initiates the password change process by sending OTP
// @Summary Prepare password change
// @Description Sends OTP to user's email for password change verification
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "New password"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /prepare-change-password [post]
func PreparePasswordChangeFlow(c *gin.Context) {
	setting_view.PrepareChangePassword()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email."})
}

// ChangePasswordFlow completes the password change after OTP verification
// @Summary Change password
// @Description Verifies OTP and updates user's password
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body OTPRequest true "OTP verification"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /change-password [post]
func ChangePasswordFlow(c *gin.Context) {
	setting_view.ChangePassword()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed your password."})
}

// PrepareDeleteAccountFlow initiates the account deletion process by sending OTP
// @Summary Prepare account deletion
// @Description Sends OTP to user's email for account deletion verification
// @Tags Setting
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401,422,500 {object} ErrorResponse
// @Router /prepare-delete-account [post]
func PrepareDeleteAccountFlow(c *gin.Context) {
	setting_view.PrepareDeleteAccount()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email."})
}

// DeleteAccountFlow completes the account deletion after OTP verification
// @Summary Delete account
// @Description Verifies OTP and deletes the authenticated user's account
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body OTPRequest true "OTP verification"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,404,422,500 {object} ErrorResponse
// @Router /delete-account [post]
func DeleteAccountFlow(c *gin.Context) {
	setting_view.DeleteAccount()(c)
	if c.IsAborted() {
		return
	}
	if public_user_id, exists := c.Get("public_user_id"); exists {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted your account."})
}

//// CLIENT SUPPORT

// ClientReportOtherClientFlow submits a report against another client
// @Summary Submit client report
// @Description Allows authenticated users to report another client
// @Tags ClientSupport
// @Accept json
// @Produce json
// @Param request body ClientReportRequest true "Report details"
// @Success 200 {object} SuccessResponse
// @Failure 401,422,500 {object} ErrorResponse
// @Router /client_report_other_client [post]
func ClientReportOtherClientFlow(c *gin.Context) {
	client_support_view.ClientReportOtherClient()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully submitted the report."})
}

// ClientCommunicateToAdminFlow submits a communication message to admin
// @Summary Submit client communication
// @Description Allows authenticated users to send a message to administrators
// @Tags ClientSupport
// @Accept json
// @Produce json
// @Param request body ClientCommunicateRequest true "Communication details"
// @Success 200 {object} SuccessResponse
// @Failure 401,422,500 {object} ErrorResponse
// @Router /client_comm_to_admin [post]
func ClientCommunicateToAdminFlow(c *gin.Context) {
	client_support_view.ClientCommunicateToAdmin()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully stored your message."})
}

///// ADMINISTRATOR

// BanClientFlow bans a client user
// @Summary Ban client
// @Description Bans a client user by their public session ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body UserControlRequest true "Ban request"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,403,500 {object} ErrorResponse
// @Router /ban_client [post]
func BanClientFlow(c *gin.Context) {
	administrator_view.BanClient()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully."})
}

// RemoveIndividualUserSessionFlow removes a specific user session
// @Summary Remove individual user session
// @Description Removes a specific user session by their public session ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body UserControlRequest true "Session removal request"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,403,500 {object} ErrorResponse
// @Router /remove_individual_session [post]
func RemoveIndividualUserSessionFlow(c *gin.Context) {
	administrator_view.RemoveIndividualUserSession()(c)
}

// RemoveAllClientSessionFlow removes all client sessions
// @Summary Remove all client sessions
// @Description Removes all client sessions from the system
// @Tags Admin
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401,500 {object} ErrorResponse
// @Router /remove_all_session [get]
func RemoveAllClientSessionFlow(c *gin.Context) {
	administrator_view.RemoveAllClientSession()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted sessions."})
}

// GetClientCommunicationLogsFlow retrieves all client communication logs
// @Summary Get client communication logs
// @Description Retrieves all client communication logs for admin review
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401,404,500 {object} ErrorResponse
// @Router /client_comm_log [get]
func GetClientCommunicationLogsFlow(c *gin.Context) {
	administrator_view.GetClientCommunicationLogs()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved the client communications data.", "data": data})
}

// ClientCommunicationReplyFlow sends a reply to client communication
// @Summary Reply to client communication
// @Description Sends a reply to a client communication message
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body ClientCommunicationReplyRequest true "Reply request"
// @Success 200 {object} SuccessResponse
// @Failure 401,404,500 {object} ErrorResponse
// @Router /client_comm_reply_log [post]
func ClientCommunicationReplyFlow(c *gin.Context) {
	administrator_view.ClientCommunicationReply()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully added the reply."})
}

// GetClientsFlow retrieves all client configurations
// @Summary Get all clients
// @Description Retrieves all client configurations for admin management
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401,404,500 {object} ErrorResponse
// @Router /get_all_clients [get]
func GetClientsFlow(c *gin.Context) {
	administrator_view.GetClients()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved the client configs data.", "data": data})
}

// GetAdminsFlow retrieves all admin configurations
// @Summary Get all admins
// @Description Retrieves all admin configurations for admin management
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 401,404,500 {object} ErrorResponse
// @Router /get_all_admins [get]
func GetAdminsFlow(c *gin.Context) {
	administrator_view.GetAdmins()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved the admin configs data.", "data": data})
}

// GetClientAuditLogsFlow retrieves client audit logs
// @Summary Get client audit logs
// @Description Retrieves audit logs for client actions
// @Tags Admin
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401,500 {object} ErrorResponse
// @Router /client_audit_logs [get]
func GetClientAuditLogsFlow(c *gin.Context) {
	administrator_view.GetClientAuditLogs()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved the client audit logs.", "data": data})
}

// GetAdminAuditLogsFlow retrieves admin audit logs
// @Summary Get admin audit logs
// @Description Retrieves audit logs for admin actions
// @Tags Admin
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401,500 {object} ErrorResponse
// @Router /admin_audit_logs [get]
func GetAdminAuditLogsFlow(c *gin.Context) {
	administrator_view.GetAdminAuditLogs()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved the admin audit logs.", "data": data})
}

// GetErrorAuditLogsFlow retrieves error audit logs
// @Summary Get error audit logs
// @Description Retrieves audit logs for system errors
// @Tags Admin
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 401,500 {object} ErrorResponse
// @Router /error_audit_logs [get]
func GetErrorAuditLogsFlow(c *gin.Context) {
	administrator_view.GetErrorAuditLogs()(c)
	if c.IsAborted() {
		return
	}
	data, _ := c.Get("response_data")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully retrieved the error logs.", "data": data})
}

// RemoveAdminFlow removes admin privileges from a user
// @Summary Remove admin
// @Description Removes admin privileges from a specified user
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body UserControlRequest true "Remove admin request"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,403,500 {object} ErrorResponse
// @Router /remove_admin [post]
func RemoveAdminFlow(c *gin.Context) {
	administrator_view.RemoveAdmin()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully removed admin privileges."})
}

// InvitationToBecomeAdminFlow sends an invitation email to a user to become admin
// @Summary Invite user to become admin
// @Description Sends an invitation email to a selected user to become an admin
// @Tags Admin
// @Accept json
// @Produce json
// @Param request body UserControlRequest true "Selected user session"
// @Success 200 {object} SuccessResponse
// @Failure 400,401,403,404,500 {object} ErrorResponse
// @Router /invite-become-admin [post]
func InvitationToBecomeAdminFlow(c *gin.Context) {
	administrator_view.InvitationToBecomeAdmin()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully sent email invitation."})
}

// AcceptToBecomeAdminFlow accepts the admin invitation and updates user type
// @Summary Accept admin invitation
// @Description Accepts the admin invitation using the token from the email link and updates the user type to admin
// @Tags Admin
// @Produce json
// @Param token path string true "Invite token"
// @Success 200 {object} SuccessResponse
// @Failure 400,404,500 {object} ErrorResponse
// @Router /accept-become-admin/{token} [get]
func AcceptToBecomeAdminFlow(c *gin.Context) {
	administrator_view.AcceptToBecomeAdmin()(c)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed the user type."})
}

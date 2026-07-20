package api

import (
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"resuming/api/middleware"
)

func APIConnect() *echo.Echo {
	router := echo.New()
	router.Use(echomw.Logger(), echomw.Recover())

	router.POST("/prepare-registeration", PrepareRegistrationFlow)
	router.POST("/register/:type-of-user", RegisterFlow)
	router.POST("/prepare-login", PrepareLoginFlow)
	router.POST("/login", LoginFlow)
	router.GET("/accept-become-admin/:token", AcceptToBecomeAdminFlow)

	authed := router.Group("")
	authed.Use(middleware.SessionCheck())
	{
		is_admin := authed.Group("")
		is_admin.Use(middleware.OnlyAdmin())
		{
			is_admin.POST("/ban_client", BanClientFlow)
			is_admin.POST("/remove_individual_session", RemoveIndividualUserSessionFlow)
			is_admin.POST("/remove_all_session", RemoveAllClientSessionFlow)
			is_admin.GET("/client_comm_log", GetSupportMessagesFlow)
			is_admin.POST("/client_comm_reply_log", ClientCommunicationReplyFlow)
			is_admin.GET("/get_all_clients", GetClientsFlow)
			is_admin.GET("/get_all_admins", GetAdminsFlow)
			is_admin.GET("/client_audit_logs", GetClientAuditLogsFlow)
			is_admin.GET("/admin_audit_logs", GetAdminAuditLogsFlow)
			is_admin.GET("/error_audit_logs", GetErrorAuditLogsFlow)

			is_super_admin := is_admin.Group("")
			is_super_admin.Use(middleware.OnlySuperAdmin())
			{
				is_super_admin.POST("/remove_admin", RemoveAdminFlow)
				is_super_admin.POST("/invite-become-admin", InvitationToBecomeAdminFlow)
			}
		}

		authed.POST("/showcaserecord-add/:type-of-data", ShowCaseRecordAddFlow)
		authed.DELETE("/showcaserecord-delete", ShowCaseRecordDeleteFlow)
		authed.PATCH("/showcaserecord-edit/:type-of-data", ShowCaseRecordEditFlow)
		authed.GET("/showcaserecord-retrieve", ShowCaseRecordGetFlow)

		authed.PATCH("/choose-portfolio-template", ChoosePortfolioTemplateFlow)
		authed.GET("/get-portfolio-content", GetPortfolioContentFlow)

		authed.PATCH("/choose-resume-template", ChooseResumeTemplateFlow)
		authed.GET("/get-resume-content", GetResumeContentFlow)

		authed.POST("/ats-score-webscrape", ATSScoreWebScrapeFlow)
		authed.POST("/ats-score-user-input", ATSScoreUserInputFlow)

		authed.POST("/change-username", ChangeUsernameFlow)
		authed.POST("/change-displayname", ChangeDisplaynameFlow)
		authed.POST("/prepare-change-email", PrepareEmailChangeFlow)
		authed.POST("/change-email", ChangeEmailFlow)
		authed.POST("/prepare-change-password", PreparePasswordChangeFlow)
		authed.POST("/change-password", ChangePasswordFlow)
		authed.POST("/prepare-delete-account", PrepareDeleteAccountFlow)
		authed.POST("/delete-account", DeleteAccountFlow)

		authed.POST("/client_comm_to_admin", ClientCommunicateToAdminFlow)
		authed.POST("/client_report_other_client", ClientReportOtherClientFlow)
	}

	return router
}

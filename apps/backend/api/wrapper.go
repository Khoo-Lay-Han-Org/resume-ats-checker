package api

// @title           Resume Builder API
// @version         1.0
// @description     API for resume building, portfolio management, and ATS scoring
// @host            localhost:5321
// @BasePath        /

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	administrator_view "resuming/api/administrator/view"
	ats_view "resuming/api/ats/view"
	auth_view "resuming/api/auth/view"
	client_support_view "resuming/api/client-support/view"
	portfolio_view "resuming/api/portfolio/view"
	resume_view "resuming/api/resume/view"
	setting_view "resuming/api/setting/view"
	showcaserecord_view "resuming/api/showcaserecord/view"
	"resuming/database"
	"resuming/database/sqlc"
)

//// AUTH

func PrepareRegistrationFlow(c echo.Context) error {
	return auth_view.PrepareRegistration()(c)
}

func RegisterFlow(c echo.Context) error {
	return auth_view.Register()(c)
}

func PrepareLoginFlow(c echo.Context) error {
	return auth_view.PrepareLogin()(c)
}

func LoginFlow(c echo.Context) error {
	if err := auth_view.Login()(c); err != nil {
		return err
	}

	auth_view.SetSession()(c)

	if private_id := c.Get("private_id"); private_id != nil {
		if public_user_id := c.Get("public_user_id"); public_user_id != nil {
			if session_key := c.Get("session_key"); session_key != nil {
				if signing_key := c.Get("signing_key"); signing_key != nil {
					if user := c.Get("user"); user != nil {
						if err := database.SyncIndividualLoginDataToSessionStore(
							public_user_id.(string),
							session_key.(string),
							signing_key.(string),
							int32(private_id.(int)),
							user.(*sqlc.User),
						); err != nil {
							log.Printf("Failed to sync login data: %v", err)
						}
					}
				}
			}
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged in."})
}

//// SHOWCASERECORD

func ShowCaseRecordAddFlow(c echo.Context) error {
	return showcaserecord_view.AddShowCaseRecordData()(c)
}

func ShowCaseRecordDeleteFlow(c echo.Context) error {
	return showcaserecord_view.DeleteShowCaseRecordData()(c)
}

func ShowCaseRecordEditFlow(c echo.Context) error {
	return showcaserecord_view.EditShowCaseRecordData()(c)
}

func ShowCaseRecordGetFlow(c echo.Context) error {
	if err := showcaserecord_view.RetrieveShowCaseRecordData()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved showcase records.", "data": data})
}

//// PORTFOLIO

func ChoosePortfolioTemplateFlow(c echo.Context) error {
	return portfolio_view.ChooseTemplate()(c)
}

func GetPortfolioContentFlow(c echo.Context) error {
	if err := portfolio_view.RetrievePortfolioData()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved portfolio data.", "data": data})
}

//// RESUME

func ChooseResumeTemplateFlow(c echo.Context) error {
	return resume_view.ChooseTemplate()(c)
}

func GetResumeContentFlow(c echo.Context) error {
	if err := resume_view.RetrieveResumeData()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved resume data.", "data": data})
}

//// ATS

func ATSScoreWebScrapeFlow(c echo.Context) error {
	if err := ats_view.ExtractResume()(c); err != nil {
		return err
	}
	if err := ats_view.ParseResume()(c); err != nil {
		return err
	}
	if err := ats_view.SectionExistenceCheck()(c); err != nil {
		return err
	}
	if err := ats_view.FormattingCheck()(c); err != nil {
		return err
	}
	if err := ats_view.WebScrapeJobDesc()(c); err != nil {
		return err
	}
	if err := ats_view.ResumeJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_view.JobDescJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_view.JobTypeRelevanceCheck()(c); err != nil {
		return err
	}
	if err := ats_view.ResumeSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_view.JobDescSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_view.OverallSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_view.OverallScore()(c); err != nil {
		return err
	}

	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully calculated ATS score.", "data": data})
}

func ATSScoreUserInputFlow(c echo.Context) error {
	if err := ats_view.ExtractResume()(c); err != nil {
		return err
	}
	if err := ats_view.ParseResume()(c); err != nil {
		return err
	}
	if err := ats_view.SectionExistenceCheck()(c); err != nil {
		return err
	}
	if err := ats_view.FormattingCheck()(c); err != nil {
		return err
	}
	if err := ats_view.UserInputJobDesc()(c); err != nil {
		return err
	}
	if err := ats_view.ResumeJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_view.JobDescJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_view.JobTypeRelevanceCheck()(c); err != nil {
		return err
	}
	if err := ats_view.ResumeSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_view.JobDescSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_view.OverallSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_view.OverallScore()(c); err != nil {
		return err
	}

	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully calculated ATS score.", "data": data})
}

//// SETTING

func ChangeUsernameFlow(c echo.Context) error {
	if err := setting_view.ChangeUsername()(c); err != nil {
		return err
	}
	if public_user_id := c.Get("public_user_id"); public_user_id != nil {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully changed your username."})
}

func ChangeDisplaynameFlow(c echo.Context) error {
	if err := setting_view.ChangeDisplayname()(c); err != nil {
		return err
	}
	if public_user_id := c.Get("public_user_id"); public_user_id != nil {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully changed your display name."})
}

func PrepareEmailChangeFlow(c echo.Context) error {
	return setting_view.PrepareChangeEmail()(c)
}

func ChangeEmailFlow(c echo.Context) error {
	if err := setting_view.ChangeEmail()(c); err != nil {
		return err
	}
	if public_user_id := c.Get("public_user_id"); public_user_id != nil {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully changed your email."})
}

func PreparePasswordChangeFlow(c echo.Context) error {
	return setting_view.PrepareChangePassword()(c)
}

func ChangePasswordFlow(c echo.Context) error {
	return setting_view.ChangePassword()(c)
}

func PrepareDeleteAccountFlow(c echo.Context) error {
	return setting_view.PrepareDeleteAccount()(c)
}

func DeleteAccountFlow(c echo.Context) error {
	if err := setting_view.DeleteAccount()(c); err != nil {
		return err
	}
	if public_user_id := c.Get("public_user_id"); public_user_id != nil {
		go func(psid string) {
			if err := database.SyncIndividualUserDataDatabase(psid); err != nil {
				log.Printf("Failed to sync user data: %v", err)
			}
		}(public_user_id.(string))
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully deleted your account."})
}

//// CLIENT SUPPORT

func ClientReportOtherClientFlow(c echo.Context) error {
	return client_support_view.ClientReportOtherClient()(c)
}

func ClientCommunicateToAdminFlow(c echo.Context) error {
	return client_support_view.ClientCommunicateToAdmin()(c)
}

///// ADMINISTRATOR

func BanClientFlow(c echo.Context) error {
	return administrator_view.BanClient()(c)
}

func RemoveIndividualUserSessionFlow(c echo.Context) error {
	return administrator_view.RemoveIndividualUserSession()(c)
}

func RemoveAllClientSessionFlow(c echo.Context) error {
	return administrator_view.RemoveAllClientSession()(c)
}

func GetSupportMessagesFlow(c echo.Context) error {
	if err := administrator_view.GetSupportMessages()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved support messages.", "data": data})
}

func ClientCommunicationReplyFlow(c echo.Context) error {
	return administrator_view.ClientCommunicationReply()(c)
}

func GetClientsFlow(c echo.Context) error {
	if err := administrator_view.GetClients()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the client configs data.", "data": data})
}

func GetAdminsFlow(c echo.Context) error {
	if err := administrator_view.GetAdmins()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the admin configs data.", "data": data})
}

func GetClientAuditLogsFlow(c echo.Context) error {
	if err := administrator_view.GetClientAuditLogs()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the client audit logs.", "data": data})
}

func GetAdminAuditLogsFlow(c echo.Context) error {
	if err := administrator_view.GetAdminAuditLogs()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the admin audit logs.", "data": data})
}

func GetErrorAuditLogsFlow(c echo.Context) error {
	if err := administrator_view.GetErrorAuditLogs()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the error logs.", "data": data})
}

func RemoveAdminFlow(c echo.Context) error {
	return administrator_view.RemoveAdmin()(c)
}

func InvitationToBecomeAdminFlow(c echo.Context) error {
	return administrator_view.InvitationToBecomeAdmin()(c)
}

func AcceptToBecomeAdminFlow(c echo.Context) error {
	return administrator_view.AcceptToBecomeAdmin()(c)
}

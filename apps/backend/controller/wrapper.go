package controller

// @title           Resume Builder API
// @version         1.0
// @description     API for resume building, portfolio management, and ATS scoring
// @host            localhost:5321
// @BasePath        /

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	adminauditlog_api "resuming/controller/adminauditlog/api"
	ats_api "resuming/controller/ats/api"
	clientauditlog_api "resuming/controller/clientauditlog/api"
	clientreportlog_api "resuming/controller/clientreportlog/api"
	clientsupportmessage_api "resuming/controller/clientsupportmessage/api"
	errorlog_api "resuming/controller/errorlog/api"
	portfolio_api "resuming/controller/portfolio/api"
	resume_api "resuming/controller/resume/api"
	session_api "resuming/controller/session/api"
	showcase_api "resuming/controller/showcase/api"
	user_api "resuming/controller/user/api"
	"resuming/database"
	"resuming/database/sqlc"
)

//// AUTH

func PrepareRegistrationFlow(c echo.Context) error {
	return user_api.PrepareRegistration()(c)
}

func RegisterFlow(c echo.Context) error {
	return user_api.Register()(c)
}

func PrepareLoginFlow(c echo.Context) error {
	return user_api.PrepareLogin()(c)
}

func LoginFlow(c echo.Context) error {
	if err := user_api.Login()(c); err != nil {
		return err
	}

	session_api.SetSession()(c)

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
	return showcase_api.AddShowCaseRecordData()(c)
}

func ShowCaseRecordDeleteFlow(c echo.Context) error {
	return showcase_api.DeleteShowCaseRecordData()(c)
}

func ShowCaseRecordEditFlow(c echo.Context) error {
	return showcase_api.EditShowCaseRecordData()(c)
}

func ShowCaseRecordGetFlow(c echo.Context) error {
	if err := showcase_api.RetrieveShowCaseRecordData()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved showcase records.", "data": data})
}

//// PORTFOLIO

func ChoosePortfolioTemplateFlow(c echo.Context) error {
	return portfolio_api.ChooseTemplate()(c)
}

func GetPortfolioContentFlow(c echo.Context) error {
	if err := portfolio_api.RetrievePortfolioData()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved portfolio data.", "data": data})
}

//// RESUME

func ChooseResumeTemplateFlow(c echo.Context) error {
	return resume_api.ChooseTemplate()(c)
}

func GetResumeContentFlow(c echo.Context) error {
	if err := resume_api.RetrieveResumeData()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved resume data.", "data": data})
}

//// ATS

func ATSScoreWebScrapeFlow(c echo.Context) error {
	if err := ats_api.ExtractResume()(c); err != nil {
		return err
	}
	if err := ats_api.ParseResume()(c); err != nil {
		return err
	}
	if err := ats_api.SectionExistenceCheck()(c); err != nil {
		return err
	}
	if err := ats_api.FormattingCheck()(c); err != nil {
		return err
	}
	if err := ats_api.WebScrapeJobDesc()(c); err != nil {
		return err
	}
	if err := ats_api.ResumeJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_api.JobDescJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_api.JobTypeRelevanceCheck()(c); err != nil {
		return err
	}
	if err := ats_api.ResumeSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_api.JobDescSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_api.OverallSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_api.OverallScore()(c); err != nil {
		return err
	}

	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully calculated ATS score.", "data": data})
}

func ATSScoreUserInputFlow(c echo.Context) error {
	if err := ats_api.ExtractResume()(c); err != nil {
		return err
	}
	if err := ats_api.ParseResume()(c); err != nil {
		return err
	}
	if err := ats_api.SectionExistenceCheck()(c); err != nil {
		return err
	}
	if err := ats_api.FormattingCheck()(c); err != nil {
		return err
	}
	if err := ats_api.UserInputJobDesc()(c); err != nil {
		return err
	}
	if err := ats_api.ResumeJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_api.JobDescJobTypeCheck()(c); err != nil {
		return err
	}
	if err := ats_api.JobTypeRelevanceCheck()(c); err != nil {
		return err
	}
	if err := ats_api.ResumeSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_api.JobDescSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_api.OverallSkillsCheck()(c); err != nil {
		return err
	}
	if err := ats_api.OverallScore()(c); err != nil {
		return err
	}

	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully calculated ATS score.", "data": data})
}

//// SETTING

func ChangeUsernameFlow(c echo.Context) error {
	if err := user_api.ChangeUsername()(c); err != nil {
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
	if err := user_api.ChangeDisplayname()(c); err != nil {
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
	return user_api.PrepareChangeEmail()(c)
}

func ChangeEmailFlow(c echo.Context) error {
	if err := user_api.ChangeEmail()(c); err != nil {
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
	return user_api.PrepareChangePassword()(c)
}

func ChangePasswordFlow(c echo.Context) error {
	return user_api.ChangePassword()(c)
}

func PrepareDeleteAccountFlow(c echo.Context) error {
	return user_api.PrepareDeleteAccount()(c)
}

func DeleteAccountFlow(c echo.Context) error {
	if err := user_api.DeleteAccount()(c); err != nil {
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
	return clientreportlog_api.ClientReportOtherClient()(c)
}

func ClientCommunicateToAdminFlow(c echo.Context) error {
	return clientsupportmessage_api.ClientCommunicateToAdmin()(c)
}

///// ADMINISTRATOR

func BanClientFlow(c echo.Context) error {
	return user_api.BanClient()(c)
}

func RemoveIndividualUserSessionFlow(c echo.Context) error {
	return session_api.RemoveIndividualUserSession()(c)
}

func RemoveAllClientSessionFlow(c echo.Context) error {
	return session_api.RemoveAllClientSession()(c)
}

func GetSupportMessagesFlow(c echo.Context) error {
	if err := clientsupportmessage_api.GetSupportMessages()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved support messages.", "data": data})
}

func ClientCommunicationReplyFlow(c echo.Context) error {
	return clientsupportmessage_api.ClientCommunicationReply()(c)
}

func GetClientsFlow(c echo.Context) error {
	if err := user_api.GetClients()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the client configs data.", "data": data})
}

func GetAdminsFlow(c echo.Context) error {
	if err := user_api.GetAdmins()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the admin configs data.", "data": data})
}

func GetClientAuditLogsFlow(c echo.Context) error {
	if err := clientauditlog_api.GetClientAuditLogs()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the client audit logs.", "data": data})
}

func GetAdminAuditLogsFlow(c echo.Context) error {
	if err := adminauditlog_api.GetAdminAuditLogs()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the admin audit logs.", "data": data})
}

func GetErrorAuditLogsFlow(c echo.Context) error {
	if err := errorlog_api.GetErrorAuditLogs()(c); err != nil {
		return err
	}
	data := c.Get("response_data")
	return c.JSON(http.StatusOK, echo.Map{"message": "Successfully retrieved the error logs.", "data": data})
}

func RemoveAdminFlow(c echo.Context) error {
	return user_api.RemoveAdmin()(c)
}

func InvitationToBecomeAdminFlow(c echo.Context) error {
	return user_api.InvitationToBecomeAdmin()(c)
}

func AcceptToBecomeAdminFlow(c echo.Context) error {
	return user_api.AcceptToBecomeAdmin()(c)
}

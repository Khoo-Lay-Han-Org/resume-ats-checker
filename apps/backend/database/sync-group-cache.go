package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	systemconfig "resuming/system-config"
	"resuming/tool"
)

func SyncGroupErrorLogSessionStore() error {
	var all_error_logs []ErrorLog
	result := NonExpiredErrorLog(DB).Find(&all_error_logs)
	if result.Error != nil {
		return fmt.Errorf("failed to sync all error logs: %w", result.Error)
	}

	var polished_all_error_log []ErrorLog

	for _, item := range all_error_logs {
		error_log_data := ErrorLog{
			UserId:    item.UserId,
			Type:      item.Type,
			PublicId:  item.PublicId,
			Message:   item.Message,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			ExpiresAt: item.ExpiresAt,
		}

		polished_all_error_log = append(polished_all_error_log, error_log_data)
	}

	serialised_error_log_data, err := json.Marshal(polished_all_error_log)
	if err != nil {
		return fmt.Errorf("failed to serialise error log data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("error_log_data").Value(string(serialised_error_log_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store error log data into session store: %w", err)
	}

	return nil
}

func SyncGroupClientAuditLogSessionStore() error {
	var all_client_audit_logs []ClientAuditLog
	result := NonExpiredClientAuditLog(DB).Find(&all_client_audit_logs)
	if result.Error != nil {
		return fmt.Errorf("failed to sync all client audit logs: %w", result.Error)
	}

	var polished_all_client_audit_log []ClientAuditLog

	for _, item := range all_client_audit_logs {
		client_autdit_log_data := ClientAuditLog{
			Type:      item.Type,
			PublicId:  item.PublicId,
			Message:   item.Message,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			ExpiresAt: item.ExpiresAt,
		}

		polished_all_client_audit_log = append(polished_all_client_audit_log, client_autdit_log_data)
	}

	serialised_client_audit_log_data, err := json.Marshal(polished_all_client_audit_log)
	if err != nil {
		return fmt.Errorf("failed to serialise client audit log data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("client_audit_log_data").Value(string(serialised_client_audit_log_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store client audit log data into session store: %w", err)
	}

	return nil
}

func SyncGroupAdminAuditLogSessionStore() error {
	var all_admin_audit_logs []AdminAuditLog
	result := NonExpiredAdminAuditLog(DB).Find(&all_admin_audit_logs)
	if result.Error != nil {
		return fmt.Errorf("failed to sync all admin audit logs: %w", result.Error)
	}

	var polished_all_admin_audit_log []AdminAuditLog

	for _, item := range all_admin_audit_logs {
		admin_autdit_log_data := AdminAuditLog{
			Type:      item.Type,
			PublicId:  item.PublicId,
			Message:   item.Message,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			ExpiresAt: item.ExpiresAt,
		}

		polished_all_admin_audit_log = append(polished_all_admin_audit_log, admin_autdit_log_data)
	}

	serialised_admin_audit_log_data, err := json.Marshal(polished_all_admin_audit_log)
	if err != nil {
		return fmt.Errorf("failed to serialise admin audit log data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("admin_audit_log_data").Value(string(serialised_admin_audit_log_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store admin audit log data into session store: %w", err)
	}

	return nil
}

func SyncGroupClientsConfigSessionStore() error {
	var all_clients []User
	result := DB.Preload("Session").Scopes(NonExpiredUser).Where("user_type = ?", "client").Find(&all_clients)
	if result.Error != nil {
		return fmt.Errorf("failed to sync all clients: %w", result.Error)
	}

	var polished_client_configs []ClientAdminConfigDTO
	for _, item := range all_clients {
		var deletedAt, bannedAt *time.Time
		if item.DeletedAt.Valid {
			t := item.DeletedAt.Time
			deletedAt = &t
		}
		if item.BannedAt != nil {
			bannedAt = item.BannedAt
		}
		client_config := ClientAdminConfigDTO{
			PublicId:     item.PublicId.String(),
			Username:     item.Username,
			Displayname:  item.Displayname,
			PublicUserId: item.Session.PublicId.String(),
			DeletedAt:    deletedAt,
			BannedAt:     bannedAt,
		}
		polished_client_configs = append(polished_client_configs, client_config)
	}

	serialised_client_configs, err := json.Marshal(polished_client_configs)
	if err != nil {
		return fmt.Errorf("failed to serialise client configs: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("client_configs").Value(string(serialised_client_configs)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store client configs: %w", err)
	}

	return nil
}

func SyncGroupAdminsConfigSessionStore() error {
	var all_admins []User
	result := DB.Preload("Session").Scopes(NonExpiredUser).Where("user_type = ?", "admin").Find(&all_admins)
	if result.Error != nil {
		return fmt.Errorf("failed to sync all admins: %w", result.Error)
	}

	var polished_admin_configs []ClientAdminConfigDTO
	for _, item := range all_admins {
		var deletedAt, bannedAt *time.Time
		if item.DeletedAt.Valid {
			t := item.DeletedAt.Time
			deletedAt = &t
		}
		if item.BannedAt != nil {
			bannedAt = item.BannedAt
		}
		admin_config := ClientAdminConfigDTO{
			PublicId:     item.PublicId.String(),
			Username:     item.Username,
			Displayname:  item.Displayname,
			PublicUserId: item.Session.PublicId.String(),
			DeletedAt:    deletedAt,
			BannedAt:     bannedAt,
		}
		polished_admin_configs = append(polished_admin_configs, admin_config)
	}

	serialised_admin_configs, err := json.Marshal(polished_admin_configs)
	if err != nil {
		return fmt.Errorf("failed to serialise admin configs: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("admin_configs").Value(string(serialised_admin_configs)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store admin configs: %w", err)
	}

	return nil
}

func SyncGroupClientReportLogSessionStore() error {
	var all_client_report_logs []ClientReportLog
	result := DB.Find(&all_client_report_logs)
	if result.Error != nil {
		return fmt.Errorf("failed to retrieve client error logs: %w", result.Error)
	}

	var polished_client_report_logs []ClientReportLogDTO
	for _, item := range all_client_report_logs {
		var reporting_user User
		if err := DB.Where("id = ?", item.ReportingUserId).First(&reporting_user).Error; err != nil {
			log.Printf("Failed to resolve reporting user for report sync %s: %v", item.PublicId, err)
			continue
		}

		var target_user User
		if err := DB.Where("id = ?", item.TargetUserId).First(&target_user).Error; err != nil {
			log.Printf("Failed to resolve target user for report sync %s: %v", item.PublicId, err)
			continue
		}

		client_report_log := ClientReportLogDTO{
			PublicId:              item.PublicId.String(),
			ReportingPublicUserId: reporting_user.PublicId.String(),
			TargetPublicUserId:    target_user.PublicId.String(),
			Type:                  string(item.Type),
		}

		polished_client_report_logs = append(polished_client_report_logs, client_report_log)
	}

	serialised_client_report_logs, err := json.Marshal(polished_client_report_logs)
	if err != nil {
		return fmt.Errorf("failed to serialise client report logs: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("client_report_logs").Value(string(serialised_client_report_logs)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store admin configs: %w", err)
	}

	return nil
}

func SyncGroupClientSupportMessagingSessionStore() error {
	var all_messages []ClientSupportMessaging
	if err := DB.Find(&all_messages).Error; err != nil {
		return fmt.Errorf("failed to get support messages: %w", err)
	}

	var polished []ClientSupportMessageDTO
	for _, item := range all_messages {
		var user User
		if err := DB.Where("id = ?", item.Content.UserId).First(&user).Error; err != nil {
			log.Printf("Failed to resolve user for support message %s: %v", item.PublicId, err)
			continue
		}

		sender_type := "client"
		if user.UserType == Admin || user.UserType == SuperAdmin {
			sender_type = "admin"
		}

		msg := ClientSupportMessageDTO{
			PublicId:              item.PublicId.String(),
			UserId:                user.PublicId.String(),
			Type:                  string(item.Type),
			Message:               item.Content.Text,
			SenderType:            sender_type,
			ClientCommLogPublicId: "",
			CreatedAt:             item.CreatedAt,
		}
		polished = append(polished, msg)
	}

	ctx := context.Background()
	serialised, err := json.Marshal(polished)
	if err != nil {
		return fmt.Errorf("failed to serialise support messages: %w", err)
	}
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("client_support_messages").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store support messages: %w", err)
	}

	return nil
}

func SyncGroupUsersSessionStore() error {
	var users []User
	result := DB.Scopes(NonExpiredUser).Preload("Resume").Preload("Portfolio").Preload("ShowcaseRecord").Preload("ATS").Find(&users)
	if result.Error != nil {
		return fmt.Errorf("failed to query users: %w", result.Error)
	}

	var polished []User
	for _, user := range users {
		polished = append(polished, User{
			PublicId:    user.PublicId,
			Username:    user.Username,
			Displayname: user.Displayname,
			Email:       user.Email,
			UserType:    user.UserType,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			BannedAt:    user.BannedAt,
			DeletedAt:   user.DeletedAt,
			ExpiresAt:   user.ExpiresAt,
		})
	}

	serialised, err := json.Marshal(polished)
	if err != nil {
		return fmt.Errorf("failed to serialise user data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("user_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store user data: %w", err)
	}

	return nil
}

func SyncGroupShowCaseRecordsSessionStore() error {
	var records []ShowcaseRecord
	result := DB.Scopes(NonExpiredShowcaseRecord).Find(&records)
	if result.Error != nil {
		return fmt.Errorf("failed to query showcase records: %w", result.Error)
	}

	var polished []ShowcaseRecord
	for _, record := range records {
		polished = append(polished, ShowcaseRecord{
			PublicId:      record.PublicId,
			UserId:        record.UserId,
			Name:          record.Name,
			Email:         record.Email,
			PhoneNumber:   record.PhoneNumber,
			Address:       record.Address,
			SocialMedia:   record.SocialMedia,
			JobExperience: record.JobExperience,
			Education:     record.Education,
			Skill:         record.Skill,
			Certificate:   record.Certificate,
			Language:      record.Language,
			Project:       record.Project,
			CreatedAt:     record.CreatedAt,
			UpdatedAt:     record.UpdatedAt,
			DeletedAt:     record.DeletedAt,
			ExpiresAt:     record.ExpiresAt,
		})
	}

	serialised, err := json.Marshal(polished)
	if err != nil {
		return fmt.Errorf("failed to serialise showcase record data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("showcaserecord_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store showcase record data: %w", err)
	}

	return nil
}

func SyncGroupPortfoliosSessionStore() error {
	var portfolios []Portfolio
	result := DB.Scopes(NonExpiredPortfolio).Find(&portfolios)
	if result.Error != nil {
		return fmt.Errorf("failed to query portfolios: %w", result.Error)
	}

	var polished []Portfolio
	for _, portfolio := range portfolios {
		polished = append(polished, Portfolio{
			PublicId:   portfolio.PublicId,
			UserId:     portfolio.UserId,
			TemplateId: portfolio.TemplateId,
			Detail:     portfolio.Detail,
			CreatedAt:  portfolio.CreatedAt,
			UpdatedAt:  portfolio.UpdatedAt,
			DeletedAt:  portfolio.DeletedAt,
			ExpiresAt:  portfolio.ExpiresAt,
		})
	}

	serialised, err := json.Marshal(polished)
	if err != nil {
		return fmt.Errorf("failed to serialise portfolio data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("portfolio_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store portfolio data: %w", err)
	}

	return nil
}

func SyncGroupResumesSessionStore() error {
	var resumes []Resume
	result := DB.Scopes(NonExpiredResume).Find(&resumes)
	if result.Error != nil {
		return fmt.Errorf("failed to query resumes: %w", result.Error)
	}

	var polished []Resume
	for _, resume := range resumes {
		polished = append(polished, Resume{
			PublicId:   resume.PublicId,
			UserId:     resume.UserId,
			TemplateId: resume.TemplateId,
			Detail:     resume.Detail,
			CreatedAt:  resume.CreatedAt,
			UpdatedAt:  resume.UpdatedAt,
			DeletedAt:  resume.DeletedAt,
			ExpiresAt:  resume.ExpiresAt,
		})
	}

	serialised, err := json.Marshal(polished)
	if err != nil {
		return fmt.Errorf("failed to serialise resume data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("resume_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store resume data: %w", err)
	}

	return nil
}

func SyncGroupAtsSessionStore() error {
	var ats_records []Ats
	result := DB.Scopes(NonExpiredAts).Find(&ats_records)
	if result.Error != nil {
		return fmt.Errorf("failed to query ATS records: %w", result.Error)
	}

	var polished []Ats
	for _, record := range ats_records {
		polished = append(polished, Ats{
			PublicId:  record.PublicId,
			UserId:    record.UserId,
			Score:     record.Score,
			Reasoning: record.Reasoning,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			DeletedAt: record.DeletedAt,
			ExpiresAt: record.ExpiresAt,
		})
	}

	serialised, err := json.Marshal(polished)
	if err != nil {
		return fmt.Errorf("failed to serialise ATS data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key("ats_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store ATS data: %w", err)
	}

	return nil
}

func SyncGroupSessionsSessionStore() error {
	var sessions []Session
	result := DB.Scopes(NonExpiredSession).Find(&sessions)
	if result.Error != nil {
		return fmt.Errorf("failed to query sessions: %w", result.Error)
	}

	var jwt_keys []JwtKey
	result = DB.Scopes(NonExpiredJwtKey).Find(&jwt_keys)
	if result.Error != nil {
		return fmt.Errorf("failed to query jwt keys: %w", result.Error)
	}

	jwt_by_user := make(map[int]JwtKey)
	for _, jk := range jwt_keys {
		jwt_by_user[jk.UserId] = jk
	}

	ctx := context.Background()
	for _, session := range sessions {
		psid := session.PublicId.String()

		session_data, _ := json.Marshal(map[string]string{"session_key": session.SessionKey})
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().
				Key(psid+":session_data").
				Value(string(session_data)).
				Ex(systemconfig.SessionExpiryDuration).
				Build(),
		).Error(); err != nil {
			return fmt.Errorf("failed to store session data for %s: %w", psid, err)
		}

		if jk, ok := jwt_by_user[session.UserId]; ok {
			jwt_data, err := json.Marshal(jk)
			if err != nil {
				return fmt.Errorf("failed to marshal jwt key for user %d: %w", session.UserId, err)
			}
			if err := tool.Valkey.Do(
				ctx,
				tool.Valkey.B().Set().
					Key(psid+":jwt_data").
					Value(string(jwt_data)).
					Ex(systemconfig.SessionExpiryDuration).
					Build(),
			).Error(); err != nil {
				return fmt.Errorf("failed to store jwt data for %s: %w", psid, err)
			}
		}
	}

	return nil
}

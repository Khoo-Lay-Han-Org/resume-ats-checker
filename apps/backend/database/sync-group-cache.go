package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func SyncGroupErrorLogSessionStore() error {
	logs, err := Queries.FindAllErrorLogs(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query error logs: %w", err)
	}

	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise error log data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("error_log_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupClientAuditLogSessionStore() error {
	logs, err := Queries.FindAllClientAuditLogs(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query client audit logs: %w", err)
	}

	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise client audit log data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("client_audit_log_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupAdminAuditLogSessionStore() error {
	logs, err := Queries.FindAllAdminAuditLogs(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query admin audit logs: %w", err)
	}

	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise admin audit log data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("admin_audit_log_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupClientsConfigSessionStore() error {
	rows, err := Queries.FindAllClientsWithSession(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query clients: %w", err)
	}

	type ClientConfig struct {
		PublicId     string `json:"public_id"`
		Username     string `json:"username"`
		Displayname  string `json:"displayname"`
		PublicUserId string `json:"public_user_id"`
	}

	var configs []ClientConfig
	for _, row := range rows {
		configs = append(configs, ClientConfig{
			PublicId:     row.PublicID.String(),
			Username:     row.Username,
			Displayname:  row.Displayname,
			PublicUserId: row.SessionPublicID.String(),
		})
	}

	serialised, err := json.Marshal(configs)
	if err != nil {
		return fmt.Errorf("failed to serialise client configs: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("client_configs").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupAdminsConfigSessionStore() error {
	rows, err := Queries.FindAllAdminsWithSession(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query admins: %w", err)
	}

	type AdminConfig struct {
		PublicId     string `json:"public_id"`
		Username     string `json:"username"`
		Displayname  string `json:"displayname"`
		PublicUserId string `json:"public_user_id"`
	}

	var configs []AdminConfig
	for _, row := range rows {
		configs = append(configs, AdminConfig{
			PublicId:     row.PublicID.String(),
			Username:     row.Username,
			Displayname:  row.Displayname,
			PublicUserId: row.SessionPublicID.String(),
		})
	}

	serialised, err := json.Marshal(configs)
	if err != nil {
		return fmt.Errorf("failed to serialise admin configs: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("admin_configs").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupClientReportLogSessionStore() error {
	rows, err := Queries.FindClientReportLogWithUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query client report logs: %w", err)
	}

	type ReportDTO struct {
		PublicId              string `json:"public_id"`
		ReportingPublicUserId string `json:"reporting_public_user_id"`
		TargetPublicUserId    string `json:"target_public_user_id"`
		Type                  string `json:"type"`
	}

	var dtos []ReportDTO
	for _, row := range rows {
		dtos = append(dtos, ReportDTO{
			PublicId:              row.PublicID.String(),
			ReportingPublicUserId: row.ReportingUserPublicID.String(),
			TargetPublicUserId:    row.TargetUserPublicID.String(),
			Type:                  row.Type,
		})
	}

	serialised, err := json.Marshal(dtos)
	if err != nil {
		return fmt.Errorf("failed to serialise client report logs: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("client_report_logs").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupClientSupportMessagingSessionStore() error {
	messages, err := Queries.FindAllClientSupportMessages(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query support messages: %w", err)
	}

	type MessageDTO struct {
		PublicId   string `json:"public_id"`
		UserId     string `json:"user_id"`
		Type       string `json:"type"`
		Message    string `json:"message"`
		SenderType string `json:"sender_type"`
		CreatedAt  string `json:"created_at"`
	}

	var dtos []MessageDTO
	for _, msg := range messages {
		var content struct {
			Text   string `json:"text"`
			UserId int    `json:"user_id"`
		}
		if err := json.Unmarshal(msg.Content, &content); err != nil {
			log.Printf("Failed to parse message content for %s: %v", msg.PublicID.String(), err)
			continue
		}

		user, err := Queries.FindUserById(context.Background(), int32(content.UserId))
		if err != nil {
			log.Printf("Failed to find user for message %s: %v", msg.PublicID.String(), err)
			continue
		}

		senderType := "client"
		if user.UserType == sqlc.UserTypeAdmin || user.UserType == sqlc.UserTypeSuperAdmin {
			senderType = "admin"
		}

		dtos = append(dtos, MessageDTO{
			PublicId:   msg.PublicID.String(),
			UserId:     user.PublicID.String(),
			Type:       string(msg.Type),
			Message:    content.Text,
			SenderType: senderType,
			CreatedAt:  msg.CreatedAt.Time.String(),
		})
	}

	serialised, err := json.Marshal(dtos)
	if err != nil {
		return fmt.Errorf("failed to serialise support messages: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("client_support_messages").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupUsersSessionStore() error {
	users, err := Queries.FindAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}

	type UserDTO struct {
		PublicId    string `json:"public_id"`
		Username    string `json:"username"`
		Displayname string `json:"displayname"`
		Email       string `json:"email"`
		UserType    string `json:"user_type"`
	}

	var dtos []UserDTO
	for _, u := range users {
		dtos = append(dtos, UserDTO{
			PublicId:    u.PublicID.String(),
			Username:    u.Username,
			Displayname: u.Displayname,
			Email:       u.Email,
			UserType:    string(u.UserType),
		})
	}

	serialised, err := json.Marshal(dtos)
	if err != nil {
		return fmt.Errorf("failed to serialise user data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("user_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupShowCaseRecordsSessionStore() error {
	records, err := Queries.FindAllShowcaseRecords(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query showcase records: %w", err)
	}

	serialised, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to serialise showcase record data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("showcaserecord_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupPortfoliosSessionStore() error {
	portfolios, err := Queries.FindAllPortfolios(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query portfolios: %w", err)
	}

	serialised, err := json.Marshal(portfolios)
	if err != nil {
		return fmt.Errorf("failed to serialise portfolio data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("portfolio_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupResumesSessionStore() error {
	resumes, err := Queries.FindAllResumes(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query resumes: %w", err)
	}

	serialised, err := json.Marshal(resumes)
	if err != nil {
		return fmt.Errorf("failed to serialise resume data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("resume_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupAtsSessionStore() error {
	ats, err := Queries.FindAllAts(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query ATS records: %w", err)
	}

	serialised, err := json.Marshal(ats)
	if err != nil {
		return fmt.Errorf("failed to serialise ATS data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().Key("ats_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).Build(),
	).Error()
}

func SyncGroupSessionsSessionStore() error {
	sessions, err := Queries.FindAllSessions(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query sessions: %w", err)
	}

	jwtKeys, err := Queries.FindAllJwtKeys(context.Background())
	if err != nil {
		return fmt.Errorf("failed to query jwt keys: %w", err)
	}

	jwtByUser := make(map[int32]sqlc.JwtKey)
	for _, jk := range jwtKeys {
		jwtByUser[jk.UserID] = jk
	}

	ctx := context.Background()
	for _, session := range sessions {
		psid := session.PublicID.String()

		sessionData, _ := json.Marshal(map[string]string{"session_key": session.SessionKey})
		if err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Set().Key(psid+":session_data").Value(string(sessionData)).
				Ex(systemconfig.SessionExpiryDuration).Build(),
		).Error(); err != nil {
			return fmt.Errorf("failed to store session data for %s: %w", psid, err)
		}

		if jk, ok := jwtByUser[session.UserID]; ok {
			jwtData, err := json.Marshal(jk)
			if err != nil {
				return fmt.Errorf("failed to marshal jwt key: %w", err)
			}
			if err := tool.Valkey.Do(
				ctx,
				tool.Valkey.B().Set().Key(psid+":jwt_data").Value(string(jwtData)).
					Ex(systemconfig.SessionExpiryDuration).Build(),
			).Error(); err != nil {
				return fmt.Errorf("failed to store jwt data for %s: %w", psid, err)
			}
		}
	}

	return nil
}

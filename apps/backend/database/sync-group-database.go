package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	valkey "github.com/valkey-io/valkey-go"
	"resuming/database/sqlc"
	"resuming/tool"
)

func SyncGroupErrorLogDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("error_log_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No error log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for error_log_data: %w", err)
	}

	var deserialised []sqlc.ErrorLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal error log data: %w", err)
	}

	for _, item := range deserialised {
		err := Queries.UpdateErrorLogByPublicId(ctx, sqlc.UpdateErrorLogByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
			Message:  item.Message,
		})
		if err != nil {
			Queries.CreateErrorLog(ctx, sqlc.CreateErrorLogParams{
				UserID:  item.UserID,
				Type:    item.Type,
				Message: item.Message,
			})
		}
	}

	log.Println("Successfully updated error logs.")
	return nil
}

func SyncGroupClientAuditLogDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_audit_log_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No client audit log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client_audit_log_data: %w", err)
	}

	var deserialised []sqlc.ClientAuditLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal client audit log data: %w", err)
	}

	for _, item := range deserialised {
		err := Queries.UpdateClientAuditLogByPublicId(ctx, sqlc.UpdateClientAuditLogByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
			Message:  item.Message,
		})
		if err != nil {
			Queries.CreateClientAuditLog(ctx, sqlc.CreateClientAuditLogParams{
				UserID:  item.UserID,
				Type:    item.Type,
				Message: item.Message,
			})
		}
	}

	log.Println("Successfully updated client audit logs.")
	return nil
}

func SyncGroupAdminAuditLogDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_audit_log_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No admin audit log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for admin_audit_log_data: %w", err)
	}

	var deserialised []sqlc.AdminAuditLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal admin audit log data: %w", err)
	}

	for _, item := range deserialised {
		err := Queries.UpdateAdminAuditLogByPublicId(ctx, sqlc.UpdateAdminAuditLogByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
			Message:  item.Message,
		})
		if err != nil {
			Queries.CreateAdminAuditLog(ctx, sqlc.CreateAdminAuditLogParams{
				UserID:  item.UserID,
				Type:    item.Type,
				Message: item.Message,
			})
		}
	}

	log.Println("Successfully updated admin audit logs.")
	return nil
}

func SyncGroupClientsConfigDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_configs").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No client config data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client_configs: %w", err)
	}

	type ClientConfig struct {
		PublicId     string `json:"public_id"`
		Username     string `json:"username"`
		Displayname  string `json:"displayname"`
		PublicUserId string `json:"public_user_id"`
	}

	var configs []ClientConfig
	if err := json.Unmarshal([]byte(data), &configs); err != nil {
		return fmt.Errorf("failed to unmarshal client configs: %w", err)
	}

	for _, cfg := range configs {
		uid := pgtype.UUID{}
		if err := uid.Scan(cfg.PublicId); err != nil {
			log.Printf("Invalid public_id in client config: %v", err)
			continue
		}

		userType := sqlc.UserTypeClient
		_, err := Queries.UpdateUserByPublicId(ctx, sqlc.UpdateUserByPublicIdParams{
			PublicID:    uid,
			Username:    cfg.Username,
			Displayname: cfg.Displayname,
			UserType:    userType,
		})
		if err != nil {
			log.Printf("Failed to update client %s: %v", cfg.PublicId, err)
		}
	}

	log.Println("Successfully updated client data.")
	return nil
}

func SyncGroupAdminsConfigDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_configs").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No admin config data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for admin_configs: %w", err)
	}

	type AdminConfig struct {
		PublicId     string `json:"public_id"`
		Username     string `json:"username"`
		Displayname  string `json:"displayname"`
		PublicUserId string `json:"public_user_id"`
	}

	var configs []AdminConfig
	if err := json.Unmarshal([]byte(data), &configs); err != nil {
		return fmt.Errorf("failed to unmarshal admin configs: %w", err)
	}

	for _, cfg := range configs {
		uid := pgtype.UUID{}
		if err := uid.Scan(cfg.PublicId); err != nil {
			log.Printf("Invalid public_id in admin config: %v", err)
			continue
		}

		userType := sqlc.UserTypeAdmin
		_, err := Queries.UpdateUserByPublicId(ctx, sqlc.UpdateUserByPublicIdParams{
			PublicID:    uid,
			Username:    cfg.Username,
			Displayname: cfg.Displayname,
			UserType:    userType,
		})
		if err != nil {
			log.Printf("Failed to update admin %s: %v", cfg.PublicId, err)
		}
	}

	log.Println("Successfully updated admin data.")
	return nil
}

func SyncGroupClientReportLogDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_report_logs").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No client report log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client_report_logs: %w", err)
	}

	type ReportDTO struct {
		PublicId              string `json:"public_id"`
		ReportingPublicUserId string `json:"reporting_public_user_id"`
		TargetPublicUserId    string `json:"target_public_user_id"`
		Type                  string `json:"type"`
	}

	var dtos []ReportDTO
	if err := json.Unmarshal([]byte(data), &dtos); err != nil {
		return fmt.Errorf("failed to parse client report logs: %w", err)
	}

	for _, item := range dtos {
		publicID := pgtype.UUID{}
		if err := publicID.Scan(item.PublicId); err != nil {
			log.Printf("Invalid public_id: %v", err)
			continue
		}

		_, err := Queries.FindClientReportLogByPublicId(ctx, publicID)
		if err == nil {
			continue
		}

		reportingUID := pgtype.UUID{}
		targetUID := pgtype.UUID{}
		reportingUID.Scan(item.ReportingPublicUserId)
		targetUID.Scan(item.TargetPublicUserId)

		reportingUser, err := Queries.FindUserByPublicId(ctx, reportingUID)
		if err != nil {
			log.Printf("Failed to resolve reporting user: %v", err)
			continue
		}

		targetUser, err := Queries.FindUserByPublicId(ctx, targetUID)
		if err != nil {
			log.Printf("Failed to resolve target user: %v", err)
			continue
		}

		Queries.CreateClientReportLog(ctx, sqlc.CreateClientReportLogParams{
			ReportingUserID: reportingUser.ID,
			TargetUserID:    targetUser.ID,
			Type:            item.Type,
		})
	}

	log.Println("Successfully synced client report logs.")
	return nil
}

func SyncGroupClientSupportMessagingDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_support_messages").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No support message data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client_support_messages: %w", err)
	}

	type MessageDTO struct {
		PublicId              string `json:"public_id"`
		UserId                string `json:"user_id"`
		Type                  string `json:"type"`
		Message               string `json:"message"`
		SenderType            string `json:"sender_type"`
		ClientCommLogPublicId string `json:"client_comm_log_public_id"`
		CreatedAt             string `json:"created_at"`
	}

	var dtos []MessageDTO
	if err := json.Unmarshal([]byte(data), &dtos); err != nil {
		return fmt.Errorf("failed to parse support message data: %w", err)
	}

	for _, item := range dtos {
		publicID := pgtype.UUID{}
		if err := publicID.Scan(item.PublicId); err != nil {
			log.Printf("Invalid public_id: %v", err)
			continue
		}

		_, err := Queries.FindClientSupportMessageByPublicId(ctx, publicID)
		if err == nil {
			continue
		}

		userUID := pgtype.UUID{}
		if err := userUID.Scan(item.UserId); err != nil {
			log.Printf("Invalid user_id: %v", err)
			continue
		}

		user, err := Queries.FindUserByPublicId(ctx, userUID)
		if err != nil {
			log.Printf("Failed to resolve user: %v", err)
			continue
		}

		content := map[string]any{
			"text":    item.Message,
			"user_id": user.ID,
			"time":    item.CreatedAt,
		}
		contentJSON, _ := json.Marshal(content)

		Queries.CreateClientSupportMessage(ctx, sqlc.CreateClientSupportMessageParams{
			Type:    item.Type,
			Content: contentJSON,
		})
	}

	log.Println("Successfully synced client support messages.")
	return nil
}

func SyncGroupShowCaseRecordsDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("showcaserecord_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No showcase record data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for showcaserecord_data: %w", err)
	}

	var deserialised []sqlc.ShowcaseRecord
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal showcase record data: %w", err)
	}

	for _, item := range deserialised {
		Queries.UpdateShowcaseRecordByPublicId(ctx, sqlc.UpdateShowcaseRecordByPublicIdParams{
			PublicID:      item.PublicID,
			Name:          item.Name,
			Email:         item.Email,
			PhoneNumber:   item.PhoneNumber,
			Address:       item.Address,
			SocialMedia:   item.SocialMedia,
			JobExperience: item.JobExperience,
			Education:     item.Education,
			Skill:         item.Skill,
			Certificate:   item.Certificate,
			Language:      item.Language,
			Project:       item.Project,
		})
	}

	log.Println("Successfully synced showcase record data.")
	return nil
}

func SyncGroupPortfoliosDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("portfolio_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No portfolio data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for portfolio_data: %w", err)
	}

	var deserialised []sqlc.Portfolio
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal portfolio data: %w", err)
	}

	for _, item := range deserialised {
		Queries.UpdatePortfolio(ctx, sqlc.UpdatePortfolioParams{
			UserID: item.UserID,
		})
	}

	log.Println("Successfully synced portfolio data.")
	return nil
}

func SyncGroupResumesDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("resume_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No resume data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for resume_data: %w", err)
	}

	var deserialised []sqlc.Resume
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal resume data: %w", err)
	}

	for _, item := range deserialised {
		Queries.UpdateResume(ctx, sqlc.UpdateResumeParams{
			UserID: item.UserID,
		})
	}

	log.Println("Successfully synced resume data.")
	return nil
}

func SyncGroupAtsDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("ats_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No ATS data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for ats_data: %w", err)
	}

	var deserialised []sqlc.At
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal ATS data: %w", err)
	}

	for _, item := range deserialised {
		score := item.Score
		reasoning := item.Reasoning
		Queries.UpdateAts(ctx, sqlc.UpdateAtsParams{
			UserID:    item.UserID,
			Score:     score,
			Reasoning: reasoning,
		})
	}

	log.Println("Successfully synced ATS data.")
	return nil
}

func SyncGroupUsersDatabase() error {
	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("user_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No user data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for user_data: %w", err)
	}

	type UserDTO struct {
		PublicId    string `json:"public_id"`
		Username    string `json:"username"`
		Displayname string `json:"displayname"`
		Email       string `json:"email"`
		UserType    string `json:"user_type"`
	}

	var dtos []UserDTO
	if err := json.Unmarshal([]byte(data), &dtos); err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	for _, dto := range dtos {
		uid := pgtype.UUID{}
		if err := uid.Scan(dto.PublicId); err != nil {
			log.Printf("Invalid public_id: %v", err)
			continue
		}

		_, err := Queries.UpdateUserByPublicId(ctx, sqlc.UpdateUserByPublicIdParams{
			PublicID:    uid,
			Username:    dto.Username,
			Email:       dto.Email,
			Displayname: dto.Displayname,
			UserType:    sqlc.UserType(dto.UserType),
		})
		if err != nil {
			log.Printf("Failed to update user %s: %v", dto.PublicId, err)
		}
	}

	if err := Queries.HardDeleteExpiredUsers(ctx); err != nil {
		return fmt.Errorf("failed to hard-delete expired users: %w", err)
	}

	log.Println("Successfully synced user data.")
	return nil
}

func SyncGroupSessionsDatabase() error {
	ctx := context.Background()

	sessions, err := Queries.FindAllSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to query sessions: %w", err)
	}

	for _, session := range sessions {
		psid := session.PublicID.String()

		exists, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Exists().Key(psid+":session_data").Build(),
		).AsInt64()
		if err != nil {
			return fmt.Errorf("failed to check session data in Valkey for %s: %w", psid, err)
		}

		if exists == 0 {
			if err := Queries.DeleteSessionByUserId(ctx, session.UserID); err != nil {
				return fmt.Errorf("failed to delete orphaned session: %w", err)
			}
			if err := Queries.DeleteJwtKeyByUserId(ctx, session.UserID); err != nil {
				return fmt.Errorf("failed to delete orphaned jwt key: %w", err)
			}
			log.Printf("Cleaned up orphaned session/jwt for user %d", session.UserID)
		}
	}

	log.Println("Successfully synced session data.")
	return nil
}

package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	valkey "github.com/valkey-io/valkey-go"
	"resuming/database/sqlc"
	"resuming/tool"
)

func SyncIndividualShowCaseRecordDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return err
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":showcaserecord_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No showcase record data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for showcase record data: %w", err)
	}

	var deserialised sqlc.ShowcaseRecord
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("showcase record data in Valkey is corrupted: %w", err)
	}

	name := deserialised.Name
	email := deserialised.Email
	phone := deserialised.PhoneNumber
	address := deserialised.Address
	social := deserialised.SocialMedia
	skill := deserialised.Skill
	language := deserialised.Language

	var jobExp, edu, cert, project []byte
	if deserialised.JobExperience != nil {
		jobExp = deserialised.JobExperience
	}
	if deserialised.Education != nil {
		edu = deserialised.Education
	}
	if deserialised.Certificate != nil {
		cert = deserialised.Certificate
	}
	if deserialised.Project != nil {
		project = deserialised.Project
	}

	if err := Queries.UpdateShowcaseRecord(ctx, sqlc.UpdateShowcaseRecordParams{
		UserID:        private_id,
		Name:          name,
		Email:         email,
		PhoneNumber:   phone,
		Address:       address,
		SocialMedia:   social,
		JobExperience: jobExp,
		Education:     edu,
		Skill:         skill,
		Certificate:   cert,
		Language:      language,
		Project:       project,
	}); err != nil {
		return fmt.Errorf("failed to update showcase record: %w", err)
	}

	log.Println("Successfully updated showcase records.")
	return nil
}

func SyncIndividualUserDataDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":user_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No user data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for user data: %w", err)
	}

	var deserialised map[string]any
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("user data in Valkey is corrupted: %w", err)
	}

	username, _ := deserialised["username"].(string)
	displayname, _ := deserialised["displayname"].(string)
	email, _ := deserialised["email"].(string)

	var userType sqlc.UserType
	if ut, ok := deserialised["user_type"].(string); ok {
		userType = sqlc.UserType(ut)
	}

	if err := Queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          private_id,
		Username:    username,
		Email:       email,
		Displayname: displayname,
		UserType:    userType,
	}); err != nil {
		return fmt.Errorf("failed to update user data: %w", err)
	}

	log.Println("Successfully updated user data.")
	return nil
}

func SyncIndividualSessionDataDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":session_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No session data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for session data: %w", err)
	}

	var deserialised map[string]any
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("session data in Valkey is corrupted: %w", err)
	}

	sessionKey, _ := deserialised["session_key"].(string)

	if err := Queries.UpdateSession(ctx, sqlc.UpdateSessionParams{
		UserID:     private_id,
		SessionKey: sessionKey,
	}); err != nil {
		return fmt.Errorf("failed to update session data: %w", err)
	}

	log.Println("Successfully updated session data.")
	return nil
}

func SyncIndividualResumeDataDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":resume_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No resume data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for resume data: %w", err)
	}

	var deserialised map[string]any
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("resume data in Valkey is corrupted: %w", err)
	}

	if err := Queries.UpdateResume(ctx, sqlc.UpdateResumeParams{
		UserID: private_id,
	}); err != nil {
		return fmt.Errorf("failed to update resume data: %w", err)
	}

	log.Println("Successfully updated resume data.")
	return nil
}

func SyncIndividualPortfolioDataDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":portfolio_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No portfolio data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for portfolio data: %w", err)
	}

	var deserialised map[string]any
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("portfolio data in Valkey is corrupted: %w", err)
	}

	if err := Queries.UpdatePortfolio(ctx, sqlc.UpdatePortfolioParams{
		UserID: private_id,
	}); err != nil {
		return fmt.Errorf("failed to update portfolio data: %w", err)
	}

	log.Println("Successfully updated portfolio data.")
	return nil
}

func SyncIndividualJWTDataDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":jwt_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No JWT data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for JWT data: %w", err)
	}

	var deserialised sqlc.JwtKey
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("JWT data in Valkey is corrupted: %w", err)
	}

	if err := Queries.UpdateJwtKey(ctx, sqlc.UpdateJwtKeyParams{
		UserID: private_id,
		Key:    deserialised.Key,
	}); err != nil {
		return fmt.Errorf("failed to update JWT data: %w", err)
	}

	log.Println("Successfully updated JWT data.")
	return nil
}

func SyncIndividualClientAuditLogDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":client_audit_log_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No client audit log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client audit log data: %w", err)
	}

	var deserialised []sqlc.ClientAuditLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("client audit log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		if err := Queries.UpdateClientAuditLogByPublicId(ctx, sqlc.UpdateClientAuditLogByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
			Message:  item.Message,
		}); err != nil {
			return fmt.Errorf("failed to update client audit log: %w", err)
		}
	}

	log.Println("Successfully updated client audit logs.")
	return nil
}

func SyncIndividualAdminAuditLogDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":admin_audit_log_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No admin audit log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for admin audit log data: %w", err)
	}

	var deserialised []sqlc.AdminAuditLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("admin audit log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		if err := Queries.UpdateAdminAuditLogByPublicId(ctx, sqlc.UpdateAdminAuditLogByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
			Message:  item.Message,
		}); err != nil {
			return fmt.Errorf("failed to update admin audit log: %w", err)
		}
	}

	log.Println("Successfully updated admin audit logs.")
	return nil
}

func SyncIndividualClientReportLogDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":client_report_log_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No client report log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client report log data: %w", err)
	}

	var deserialised []sqlc.ClientReportLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("client report log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		if err := Queries.UpdateClientReportLogByPublicId(ctx, sqlc.UpdateClientReportLogByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
		}); err != nil {
			return fmt.Errorf("failed to update client report log: %w", err)
		}
	}

	log.Println("Successfully updated client report logs.")
	return nil
}

func SyncIndividualErrorLogDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":error_log_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No error log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for error log data: %w", err)
	}

	var deserialised []sqlc.ErrorLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("error log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		if err := Queries.UpdateErrorLogByPublicId(ctx, sqlc.UpdateErrorLogByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
			Message:  item.Message,
		}); err != nil {
			return fmt.Errorf("failed to update error log: %w", err)
		}
	}

	log.Println("Successfully updated error logs.")
	return nil
}

func SyncIndividualClientSupportMessagingDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":client_support_messages").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No support message data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for support message data: %w", err)
	}

	var deserialised []sqlc.ClientSupportMessaging
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("support message data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		if err := Queries.UpdateClientSupportMessageByPublicId(ctx, sqlc.UpdateClientSupportMessageByPublicIdParams{
			PublicID: item.PublicID,
			Type:     item.Type,
			Content:  item.Content,
		}); err != nil {
			return fmt.Errorf("failed to update support message: %w", err)
		}
	}

	log.Println("Successfully updated support messages.")
	return nil
}

func SyncIndividualATSDataDatabase(public_user_id string) error {
	private_id, err := resolveUserId(public_user_id)
	if err != nil {
		return fmt.Errorf("failed to resolve user ID: %w", err)
	}
	if private_id == 0 {
		return nil
	}

	ctx := context.Background()
	data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key(public_user_id+":ats_data").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No ATS data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for ATS data: %w", err)
	}

	var deserialised sqlc.At
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("ATS data in Valkey is corrupted: %w", err)
	}

	if err := Queries.UpdateAts(ctx, sqlc.UpdateAtsParams{
		UserID:    private_id,
		Score:     deserialised.Score,
		Reasoning: deserialised.Reasoning,
	}); err != nil {
		return fmt.Errorf("failed to update ATS data: %w", err)
	}

	log.Println("Successfully updated ATS data.")
	return nil
}

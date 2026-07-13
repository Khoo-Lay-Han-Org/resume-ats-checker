package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	valkey "github.com/valkey-io/valkey-go"
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
		return fmt.Errorf("valkey GET failed for showcase record data (key: %s:showcaserecord_data): %w", public_user_id, err)
	}

	json_data := []byte(data)
	var deserialised_data ShowcaseRecord
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("showcase record data in Valkey is corrupted: %w", err)
	}

	result := DB.Model(&ShowcaseRecord{}).Where("user_id = ?", private_id).Updates(ShowcaseRecord{
		Name:          deserialised_data.Name,
		Email:         deserialised_data.Email,
		PhoneNumber:   deserialised_data.PhoneNumber,
		Address:       deserialised_data.Address,
		SocialMedia:   deserialised_data.SocialMedia,
		JobExperience: deserialised_data.JobExperience,
		Education:     deserialised_data.Education,
		Skill:         deserialised_data.Skill,
		Certificate:   deserialised_data.Certificate,
		Language:      deserialised_data.Language,
		Project:       deserialised_data.Project,
		CreatedAt:     deserialised_data.CreatedAt,
		UpdatedAt:     deserialised_data.UpdatedAt,
		DeletedAt:     deserialised_data.DeletedAt,
		ExpiresAt:     deserialised_data.ExpiresAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update showcase record: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for user data (key: %s:user_data): %w", public_user_id, err)
	}

	json_data := []byte(data)
	var deserialised_data User
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("user data in Valkey is corrupted: %w", err)
	}

	result := DB.Model(&User{}).Where("id = ?", private_id).Updates(User{
		PublicId:    deserialised_data.PublicId,
		Username:    deserialised_data.Username,
		Displayname: deserialised_data.Displayname,
		Email:       deserialised_data.Email,
		UserType:    deserialised_data.UserType,
		CreatedAt:   deserialised_data.CreatedAt,
		UpdatedAt:   deserialised_data.UpdatedAt,
		BannedAt:    deserialised_data.BannedAt,
		DeletedAt:   deserialised_data.DeletedAt,
		ExpiresAt:   deserialised_data.ExpiresAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update user data: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for session data (key: %s:session_data): %w", public_user_id, err)
	}

	json_data := []byte(data)
	var deserialised_data map[string]any
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("session data in Valkey is corrupted: %w", err)
	}

	session_key, _ := deserialised_data["session_key"].(string)

	result := DB.Model(&Session{}).Where("user_id = ?", private_id).Updates(Session{
		SessionKey: session_key,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update session data: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for resume data (key: %s:resume_data): %w", public_user_id, err)
	}

	json_data := []byte(data)
	var deserialised_data Resume
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("resume data in Valkey is corrupted: %w", err)
	}

	result := DB.Model(&Resume{}).Where("user_id = ?", private_id).Updates(Resume{
		TemplateId: deserialised_data.TemplateId,
		Detail:     deserialised_data.Detail,
		CreatedAt:  deserialised_data.CreatedAt,
		UpdatedAt:  deserialised_data.UpdatedAt,
		DeletedAt:  deserialised_data.DeletedAt,
		ExpiresAt:  deserialised_data.ExpiresAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update resume data: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for portfolio data (key: %s:portfolio_data): %w", public_user_id, err)
	}

	json_data := []byte(data)
	var deserialised_data Portfolio
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("portfolio data in Valkey is corrupted: %w", err)
	}

	result := DB.Model(&Portfolio{}).Where("user_id = ?", private_id).Updates(Portfolio{
		TemplateId: deserialised_data.TemplateId,
		Detail:     deserialised_data.Detail,
		CreatedAt:  deserialised_data.CreatedAt,
		UpdatedAt:  deserialised_data.UpdatedAt,
		DeletedAt:  deserialised_data.DeletedAt,
		ExpiresAt:  deserialised_data.ExpiresAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update portfolio data: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for JWT data (key: %s:jwt_data): %w", public_user_id, err)
	}

	json_data := []byte(data)
	var deserialised_data JwtKey
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("JWT data in Valkey is corrupted: %w", err)
	}

	result := DB.Model(&JwtKey{}).Where("user_id = ?", private_id).Updates(JwtKey{
		Key:       deserialised_data.Key,
		CreatedAt: deserialised_data.CreatedAt,
		UpdatedAt: deserialised_data.UpdatedAt,
		ExpiresAt: deserialised_data.ExpiresAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update JWT data: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for client audit log data (key: %s:client_audit_log_data): %w", public_user_id, err)
	}

	var deserialised []ClientAuditLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("client audit log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&ClientAuditLog{}).Where("public_id = ?", item.PublicId).Updates(ClientAuditLog{
			Type:    item.Type,
			Message: item.Message,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update client audit log: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for admin audit log data (key: %s:admin_audit_log_data): %w", public_user_id, err)
	}

	var deserialised []AdminAuditLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("admin audit log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&AdminAuditLog{}).Where("public_id = ?", item.PublicId).Updates(AdminAuditLog{
			Type:    item.Type,
			Message: item.Message,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update admin audit log: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for client report log data (key: %s:client_report_log_data): %w", public_user_id, err)
	}

	var deserialised []ClientReportLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("client report log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&ClientReportLog{}).Where("public_id = ?", item.PublicId).Updates(ClientReportLog{
			Type: item.Type,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update client report log: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for error log data (key: %s:error_log_data): %w", public_user_id, err)
	}

	var deserialised []ErrorLog
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("error log data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&ErrorLog{}).Where("public_id = ?", item.PublicId).Updates(ErrorLog{
			Type:    item.Type,
			Message: item.Message,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update error log: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for support message data (key: %s:client_support_messages): %w", public_user_id, err)
	}

	var deserialised []ClientSupportMessaging
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("support message data in Valkey is corrupted: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&ClientSupportMessaging{}).Where("public_id = ?", item.PublicId).Updates(ClientSupportMessaging{
			Type:    item.Type,
			Content: item.Content,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update support message: %w", result.Error)
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
		return fmt.Errorf("valkey GET failed for ATS data (key: %s:ats_data): %w", public_user_id, err)
	}

	json_data := []byte(data)
	var deserialised_data Ats
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("ATS data in Valkey is corrupted: %w", err)
	}

	result := DB.Model(&Ats{}).Where("user_id = ?", private_id).Updates(Ats{
		Score:     deserialised_data.Score,
		Reasoning: deserialised_data.Reasoning,
		CreatedAt: deserialised_data.CreatedAt,
		UpdatedAt: deserialised_data.UpdatedAt,
		DeletedAt: deserialised_data.DeletedAt,
		ExpiresAt: deserialised_data.ExpiresAt,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update ATS data: %w", result.Error)
	}

	log.Println("Successfully updated ATS data.")
	return nil
}

package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	valkey "github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
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

	json_data := []byte(data)
	var deserialised_data []ErrorLog
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal error log data: %w", err)
	}

	for _, item := range deserialised_data {
		result := DB.Model(&ErrorLog{}).Where("public_id = ?", item.PublicId).Updates(ErrorLog{
			UserId:    item.UserId,
			Type:      item.Type,
			PublicId:  item.PublicId,
			Message:   item.Message,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			ExpiresAt: item.ExpiresAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update error log: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			new_error_log := ErrorLog{
				UserId:   item.UserId,
				PublicId: item.PublicId,
				Type:     item.Type,
				Message:  item.Message,
			}
			if err := DB.Create(&new_error_log).Error; err != nil {
				return fmt.Errorf("failed to create error log: %w", err)
			}
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

	json_data := []byte(data)
	var deserialised_data []ClientAuditLog
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal client audit log data: %w", err)
	}

	for _, item := range deserialised_data {
		result := DB.Model(&ClientAuditLog{}).Where("public_id = ?", item.PublicId).Updates(ClientAuditLog{
			UserId:    item.UserId,
			Type:      item.Type,
			PublicId:  item.PublicId,
			Message:   item.Message,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			ExpiresAt: item.ExpiresAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update client audit log: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			new_client_audit_log := ClientAuditLog{
				UserId:   item.UserId,
				PublicId: item.PublicId,
				Type:     item.Type,
				Message:  item.Message,
			}
			if err := DB.Create(&new_client_audit_log).Error; err != nil {
				return fmt.Errorf("failed to create client audit log: %w", err)
			}
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

	json_data := []byte(data)
	var deserialised_data []AdminAuditLog
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal admin audit log data: %w", err)
	}

	for _, item := range deserialised_data {
		result := DB.Model(&AdminAuditLog{}).Where("public_id = ?", item.PublicId).Updates(AdminAuditLog{
			UserId:    item.UserId,
			Type:      item.Type,
			PublicId:  item.PublicId,
			Message:   item.Message,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			ExpiresAt: item.ExpiresAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update admin audit log: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			new_admin_audit_log := AdminAuditLog{
				UserId:   item.UserId,
				PublicId: item.PublicId,
				Type:     item.Type,
				Message:  item.Message,
			}
			if err := DB.Create(&new_admin_audit_log).Error; err != nil {
				return fmt.Errorf("failed to create admin audit log: %w", err)
			}
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

	json_data := []byte(data)
	var deserialised_data []ClientAdminConfigDTO
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal client configs data: %w", err)
	}

	for _, item := range deserialised_data {
		public_id := uuid.MustParse(item.PublicId)
		var deletedAt gorm.DeletedAt
		if item.DeletedAt != nil {
			deletedAt = gorm.DeletedAt{Time: *item.DeletedAt, Valid: true}
		}
		result := DB.Model(&User{}).Where("public_id = ?", public_id).Updates(User{
			Username:    item.Username,
			Displayname: item.Displayname,
			DeletedAt:   deletedAt,
			BannedAt:    item.BannedAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update user: %w", result.Error)
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

	json_data := []byte(data)
	var deserialised_data []ClientAdminConfigDTO
	err = json.Unmarshal(json_data, &deserialised_data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal admin configs data: %w", err)
	}

	for _, item := range deserialised_data {
		public_id := uuid.MustParse(item.PublicId)
		var deletedAt gorm.DeletedAt
		if item.DeletedAt != nil {
			deletedAt = gorm.DeletedAt{Time: *item.DeletedAt, Valid: true}
		}

		result := DB.Model(&User{}).Where("public_id = ?", public_id).Updates(User{
			Username:    item.Username,
			Displayname: item.Displayname,
			DeletedAt:   deletedAt,
			BannedAt:    item.BannedAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update user: %w", result.Error)
		}
	}

	log.Println("Successfully updated admin data.")
	return nil
}

func SyncGroupClientReportLogDatabase() error {
	ctx := context.Background()
	retrieved_data, err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Get().
			Key("client_report_logs").
			Build(),
	).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No client report log data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client_report_logs: %w", err)
	}

	var data []ClientReportLogDTO
	err = json.Unmarshal([]byte(retrieved_data), &data)
	if err != nil {
		return fmt.Errorf("failed to parse client report logs: %w", err)
	}

	for _, item := range data {
		public_id := uuid.MustParse(item.PublicId)
		reporting_public_user_id := uuid.MustParse(item.ReportingPublicUserId)
		target_public_user_id := uuid.MustParse(item.TargetPublicUserId)
		type_of_report := item.Type

		var report_type ClientBannedReasonType
		switch type_of_report {
		case "profanity":
			report_type = Profanity
		case "explicit content":
			report_type = ExplicitContent
		}

		var reporting_user User
		if err := DB.Where("public_id = ?", reporting_public_user_id).First(&reporting_user).Error; err != nil {
			log.Printf("Failed to resolve reporting user for report %s: %v", item.PublicId, err)
			continue
		}

		var target_user User
		if err := DB.Where("public_id = ?", target_public_user_id).First(&target_user).Error; err != nil {
			log.Printf("Failed to resolve target user for report %s: %v", item.PublicId, err)
			continue
		}

		var existing_report ClientReportLog
		if err := DB.Where("public_id = ?", public_id).First(&existing_report).Error; err == nil {
			continue
		}

		DB.Create(&ClientReportLog{
			PublicId:        public_id,
			ReportingUserId: reporting_user.Id,
			TargetUserId:    target_user.Id,
			Type:            report_type,
		})
	}

	log.Println("Successfully synced client report logs.")
	return nil
}

func SyncGroupClientCommLogDatabase() error {
	ctx := context.Background()
	retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("client_comms").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No client communication data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for client_comms: %w", err)
	}

	var all_client_comm []ClientCommDTO
	err = json.Unmarshal([]byte(retrieved_data), &all_client_comm)
	if err != nil {
		return fmt.Errorf("failed to serialise all client communications data: %w", err)
	}

	for _, item := range all_client_comm {
		public_id := uuid.MustParse(item.PublicId)
		user_id := item.UserId
		if user_id == 0 {
			var user User
			if err := DB.Where("public_id = ?", item.PublicUserId).First(&user).Error; err != nil {
				log.Printf("Failed to resolve user for comm %s: %v", item.PublicId, err)
				continue
			}
			user_id = user.Id
		}
		retrieved_message_type := item.Type
		message := item.Message

		var message_type ClientCommType
		switch retrieved_message_type {
		case "technical support":
			message_type = TechnicalSupport
		case "feature improvement":
			message_type = FeatureImprovement
		case "billing management":
			message_type = BillingManagement
		case "service and operation":
			message_type = ServiceAndOperation
		case "onboarding support":
			message_type = OnboardingSupport
		case "complaint":
			message_type = Complaint
		}

		var client_comm ClientCommLog
		err := DB.Where("public_id = ?", public_id).First(&client_comm).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				new_client_comm := ClientCommLog{
					PublicId: public_id,
					UserId:   user_id,
					Type:     message_type,
					Message:  message,
				}

				err = DB.Create(&new_client_comm).Error
				if err != nil {
					return fmt.Errorf("failed to store new client messages: %w", err)
				}
			} else {
				return fmt.Errorf("failed to serialise all client communications data: %w", err)
			}
		}
	}

	log.Println("Successfully synced client communication logs.")
	return nil
}

func SyncGroupAdminCommLogDatabase() error {
	ctx := context.Background()
	retrieved_data, err := tool.Valkey.Do(ctx, tool.Valkey.B().Get().Key("admin_comms").Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			log.Println("No admin communication data in Valkey — nothing to sync.")
			return nil
		}
		return fmt.Errorf("valkey GET failed for admin_comms: %w", err)
	}

	var data []AdminCommDTO
	err = json.Unmarshal([]byte(retrieved_data), &data)
	if err != nil {
		return fmt.Errorf("failed to parse admin communications data: %w", err)
	}

	for _, item := range data {
		admin_comm_public_id := uuid.MustParse(item.AdminCommPublicId)
		message := item.Message

		var admin_user User
		if err := DB.Where("public_id = ?", item.AdminUserPublicId).First(&admin_user).Error; err != nil {
			log.Printf("Failed to resolve admin user for admin comm %s: %v", item.AdminCommPublicId, err)
			continue
		}

		var client_comm ClientCommLog
		if err := DB.Where("public_id = ?", item.ClientCommLogPublicId).First(&client_comm).Error; err != nil {
			log.Printf("Failed to resolve client comm for admin comm %s: %v", item.AdminCommPublicId, err)
			continue
		}

		var admin_comm AdminCommLog
		err := DB.Where("public_id = ?", admin_comm_public_id).First(&admin_comm).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := DB.Create(&AdminCommLog{
					PublicId:        admin_comm_public_id,
					UserId:          admin_user.Id,
					ClientCommLogId: client_comm.Id,
					Message:         message,
				}).Error; err != nil {
					return fmt.Errorf("failed to store new admin messages: %w", err)
				}
			} else {
				return fmt.Errorf("failed to serialise all admin communications data: %w", err)
			}
		}
	}

	log.Println("Successfully synced admin communication logs.")
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

	var deserialised []ShowcaseRecord
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal showcase record data: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&ShowcaseRecord{}).Where("public_id = ?", item.PublicId).Updates(ShowcaseRecord{
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
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
			DeletedAt:     item.DeletedAt,
			ExpiresAt:     item.ExpiresAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update showcase record: %w", result.Error)
		}
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

	var deserialised []Portfolio
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal portfolio data: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&Portfolio{}).Where("public_id = ?", item.PublicId).Updates(Portfolio{
			TemplateId: item.TemplateId,
			Detail:     item.Detail,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			DeletedAt:  item.DeletedAt,
			ExpiresAt:  item.ExpiresAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update portfolio: %w", result.Error)
		}
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

	var deserialised []Resume
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal resume data: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&Resume{}).Where("public_id = ?", item.PublicId).Updates(Resume{
			TemplateId: item.TemplateId,
			Detail:     item.Detail,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			DeletedAt:  item.DeletedAt,
			ExpiresAt:  item.ExpiresAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update resume: %w", result.Error)
		}
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

	var deserialised []Ats
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal ATS data: %w", err)
	}

	for _, item := range deserialised {
		result := DB.Model(&Ats{}).Where("public_id = ?", item.PublicId).Updates(Ats{
			Score:     item.Score,
			Reasoning: item.Reasoning,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: item.DeletedAt,
			ExpiresAt: item.ExpiresAt,
		})
		if result.Error != nil {
			return fmt.Errorf("failed to update ATS record: %w", result.Error)
		}
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

	var deserialised []User
	if err := json.Unmarshal([]byte(data), &deserialised); err != nil {
		return fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	for _, item := range deserialised {
		var deletedAt gorm.DeletedAt
		if item.DeletedAt.Valid {
			deletedAt = gorm.DeletedAt{Time: item.DeletedAt.Time, Valid: true}
		}

		if err := DB.Model(&User{}).Where("public_id = ?", item.PublicId).Updates(User{
			Username:    item.Username,
			Displayname: item.Displayname,
			Email:       item.Email,
			UserType:    item.UserType,
			BannedAt:    item.BannedAt,
			DeletedAt:   deletedAt,
			ExpiresAt:   item.ExpiresAt,
		}).Error; err != nil {
			return fmt.Errorf("failed to update user %s: %w", item.PublicId.String(), err)
		}
	}

	if err := DB.Unscoped().Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&User{}).Error; err != nil {
		return fmt.Errorf("failed to hard-delete expired users: %w", err)
	}

	log.Println("Successfully synced user data.")
	return nil
}

func SyncGroupSessionsDatabase() error {
	var sessions []Session
	result := DB.Scopes(NonExpiredSession).Find(&sessions)
	if result.Error != nil {
		return fmt.Errorf("failed to query sessions: %w", result.Error)
	}

	ctx := context.Background()
	for _, session := range sessions {
		psid := session.PublicId.String()

		exists, err := tool.Valkey.Do(
			ctx,
			tool.Valkey.B().Exists().Key(psid+":session_data").Build(),
		).AsInt64()
		if err != nil {
			return fmt.Errorf("failed to check session data in Valkey for %s: %w", psid, err)
		}

		if exists == 0 {
			if err := DB.Where("user_id = ?", session.UserId).Delete(&Session{}).Error; err != nil {
				return fmt.Errorf("failed to delete orphaned session for user %d: %w", session.UserId, err)
			}
			if err := DB.Where("user_id = ?", session.UserId).Delete(&JwtKey{}).Error; err != nil {
				return fmt.Errorf("failed to delete orphaned jwt key for user %d: %w", session.UserId, err)
			}
			log.Printf("cleaned up orphaned session/jwt for user %d (session %s)", session.UserId, psid)
		}
	}

	log.Println("Successfully synced session data.")
	return nil
}

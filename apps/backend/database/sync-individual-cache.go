package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func SyncIndividualUserDataSessionStore(public_user_id string, user *User) error {
	user_data := User{
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
	}

	serialised_user_data, err := json.Marshal(user_data)
	if err != nil {
		return fmt.Errorf("failed to serialise user data for user data sync: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":user_data").Value(string(serialised_user_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store user data into session store: %w", err)
	}

	return nil
}

func SyncIndividualShowCaseRecordDataSessionStore(public_user_id string, user *User) error {
	user_data := map[string]any{
		"name":           user.ShowcaseRecord.Name,
		"email":          user.ShowcaseRecord.Email,
		"phone_number":   user.ShowcaseRecord.PhoneNumber,
		"address":        user.ShowcaseRecord.Address,
		"social_media":   user.ShowcaseRecord.SocialMedia,
		"job_experience": user.ShowcaseRecord.JobExperience,
		"education":      user.ShowcaseRecord.Education,
		"skill":          user.ShowcaseRecord.Skill,
		"certificate":    user.ShowcaseRecord.Certificate,
		"language":       user.ShowcaseRecord.Language,
		"project":        user.ShowcaseRecord.Project,
	}

	serialised_showcaserecord_data, err := json.Marshal(user_data)
	if err != nil {
		return fmt.Errorf("failed to serialise showcase records data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":showcaserecord_data").Value(string(serialised_showcaserecord_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store showcase records data: %w", err)
	}

	return nil
}

func SyncIndividualSessionDataSessionStore(public_user_id string, session_key string) error {
	session_data := map[string]string{"session_key": session_key}

	serialised_session_data, err := json.Marshal(session_data)
	if err != nil {
		return fmt.Errorf("failed to serialise session data for session data sync: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":session_data").Value(string(serialised_session_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store session data into session store: %w", err)
	}

	return nil
}

func SyncIndividualResumeDataSessionStore(public_user_id string, user *User) error {
	user_data := Resume{
		PublicId:   user.Resume.PublicId,
		TemplateId: user.Resume.TemplateId,
		Detail:     user.Resume.Detail,
		CreatedAt:  user.Resume.CreatedAt,
		UpdatedAt:  user.Resume.UpdatedAt,
		DeletedAt:  user.Resume.DeletedAt,
		ExpiresAt:  user.Resume.ExpiresAt,
	}

	serialised_resume_data, err := json.Marshal(user_data)
	if err != nil {
		return fmt.Errorf("failed to serialise resume data for resume data sync: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":resume_data").Value(string(serialised_resume_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store resume data into session store: %w", err)
	}

	return nil
}

func SyncIndividualPortfolioDataSessionStore(public_user_id string, user *User) error {
	user_data := Portfolio{
		PublicId:   user.Portfolio.PublicId,
		TemplateId: user.Portfolio.TemplateId,
		Detail:     user.Portfolio.Detail,
		CreatedAt:  user.Portfolio.CreatedAt,
		UpdatedAt:  user.Portfolio.UpdatedAt,
		DeletedAt:  user.Portfolio.DeletedAt,
		ExpiresAt:  user.Portfolio.ExpiresAt,
	}

	serialised_portfolio_data, err := json.Marshal(user_data)
	if err != nil {
		return fmt.Errorf("failed to serialise portfolio data for portfolio data sync: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":portfolio_data").Value(string(serialised_portfolio_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store portfolio data into session store: %w", err)
	}

	return nil
}

func SyncIndividualJWTDataSessionStore(public_user_id string, jwt_key *JwtKey) error {
	serialised_jwt_data, err := json.Marshal(jwt_key)
	if err != nil {
		return fmt.Errorf("failed to serialise JWT data for JWT data sync: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":jwt_data").Value(string(serialised_jwt_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store JWT data into session store: %w", err)
	}

	return nil
}

func SyncIndividualUserSessionMapping(upid string, psid string) error {
	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(upid+":session_id").Value(psid).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
}

func SyncIndividualLoginDataToSessionStore(psid string, sessionKey string, signingKey string, privateId int, user *User) error {
	if err := SyncIndividualUserDataSessionStore(psid, user); err != nil {
		return err
	}
	if err := SyncIndividualUserDataSessionStore(user.PublicId.String(), user); err != nil {
		return err
	}
	new_jwt_key := JwtKey{
		PublicId: uuid.New(),
		UserId:   privateId,
		Key:      signingKey,
	}
	if err := SyncIndividualJWTDataSessionStore(psid, &new_jwt_key); err != nil {
		return err
	}
	if err := SyncIndividualShowCaseRecordDataSessionStore(psid, user); err != nil {
		return err
	}
	if err := SyncIndividualShowCaseRecordDataSessionStore(user.PublicId.String(), user); err != nil {
		return err
	}
	if err := SyncIndividualResumeDataSessionStore(psid, user); err != nil {
		return err
	}
	if err := SyncIndividualResumeDataSessionStore(user.PublicId.String(), user); err != nil {
		return err
	}
	if err := SyncIndividualPortfolioDataSessionStore(psid, user); err != nil {
		return err
	}
	if err := SyncIndividualPortfolioDataSessionStore(user.PublicId.String(), user); err != nil {
		return err
	}
	if err := SyncIndividualATSDataSessionStore(psid, user); err != nil {
		return err
	}
	if err := SyncIndividualATSDataSessionStore(user.PublicId.String(), user); err != nil {
		return err
	}
	if err := SyncIndividualUserSessionMapping(user.PublicId.String(), psid); err != nil {
		return err
	}

	go func() {
		DB.Where("user_id = ?", privateId).Delete(&JwtKey{})
		DB.Where("user_id = ?", privateId).Delete(&Session{})
		DB.Create(&new_jwt_key)
		DB.Create(&Session{
			PublicId:   uuid.MustParse(psid),
			UserId:     privateId,
			SessionKey: sessionKey,
		})
	}()

	return nil
}

func SyncIndividualClientAuditLogSessionStore(public_user_id string, logs []ClientAuditLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise client audit log data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":client_audit_log_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store client audit log data: %w", err)
	}

	return nil
}

func SyncIndividualAdminAuditLogSessionStore(public_user_id string, logs []AdminAuditLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise admin audit log data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":admin_audit_log_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store admin audit log data: %w", err)
	}

	return nil
}

func SyncIndividualClientReportLogSessionStore(public_user_id string, logs []ClientReportLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise client report log data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":client_report_log_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store client report log data: %w", err)
	}

	return nil
}

func SyncIndividualErrorLogSessionStore(public_user_id string, logs []ErrorLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise error log data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":error_log_data").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store error log data: %w", err)
	}

	return nil
}

func SyncIndividualClientSupportMessagingSessionStore(public_user_id string, messages []ClientSupportMessaging) error {
	serialised, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to serialise support message data: %w", err)
	}

	ctx := context.Background()
	if err := tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":client_support_messages").
			Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error(); err != nil {
		return fmt.Errorf("failed to store support message data: %w", err)
	}

	return nil
}

func SyncIndividualATSDataSessionStore(public_user_id string, user *User) error {
	user_data := Ats{
		PublicId:  user.ATS.PublicId,
		Score:     user.ATS.Score,
		Reasoning: user.ATS.Reasoning,
		CreatedAt: user.ATS.CreatedAt,
		UpdatedAt: user.ATS.UpdatedAt,
		DeletedAt: user.ATS.DeletedAt,
		ExpiresAt: user.ATS.ExpiresAt,
	}

	serialised_ats_data, err := json.Marshal(user_data)
	if err != nil {
		return fmt.Errorf("failed to serialise ATS data for ATS data sync: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":ats_data").Value(string(serialised_ats_data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store ATS data into session store: %w", err)
	}

	return nil
}

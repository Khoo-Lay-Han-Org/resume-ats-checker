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

func SyncIndividualUserDataSessionStore(public_user_id string, user *sqlc.User) error {
	userData := map[string]any{
		"public_id":   user.PublicID.Bytes,
		"username":    user.Username,
		"displayname": user.Displayname,
		"email":       user.Email,
		"user_type":   string(user.UserType),
		"created_at":  user.CreatedAt.Time,
		"updated_at":  user.UpdatedAt.Time,
		"banned_at":   nil,
		"deleted_at":  nil,
		"expires_at":  nil,
	}
	if user.BannedAt.Valid {
		userData["banned_at"] = user.BannedAt.Time
	}
	if user.DeletedAt.Valid {
		userData["deleted_at"] = user.DeletedAt.Time
	}
	if user.ExpiresAt.Valid {
		userData["expires_at"] = user.ExpiresAt.Time
	}

	serialised, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("failed to serialise user data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":user_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store user data into session store: %w", err)
	}

	return nil
}

func SyncIndividualShowCaseRecordDataSessionStore(public_user_id string, showcase *sqlc.ShowcaseRecord) error {
	userData := map[string]any{
		"name":           showcase.Name,
		"email":          showcase.Email,
		"phone_number":   showcase.PhoneNumber,
		"address":        showcase.Address,
		"social_media":   showcase.SocialMedia,
		"job_experience": showcase.JobExperience,
		"education":      showcase.Education,
		"skill":          showcase.Skill,
		"certificate":    showcase.Certificate,
		"language":       showcase.Language,
		"project":        showcase.Project,
	}

	serialised, err := json.Marshal(userData)
	if err != nil {
		return fmt.Errorf("failed to serialise showcase records data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":showcaserecord_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store showcase records data: %w", err)
	}

	return nil
}

func SyncIndividualSessionDataSessionStore(public_user_id string, session_key string) error {
	sessionData := map[string]string{"session_key": session_key}
	serialised, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to serialise session data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":session_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store session data: %w", err)
	}

	return nil
}

func SyncIndividualJWTDataSessionStore(public_user_id string, jwt_key *sqlc.JwtKey) error {
	serialised, err := json.Marshal(jwt_key)
	if err != nil {
		return fmt.Errorf("failed to serialise JWT data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":jwt_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store JWT data: %w", err)
	}

	return nil
}

func SyncIndividualATSDataSessionStore(public_user_id string, ats *sqlc.At) error {
	serialised, err := json.Marshal(ats)
	if err != nil {
		return fmt.Errorf("failed to serialise ATS data: %w", err)
	}

	ctx := context.Background()
	err = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":ats_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
	if err != nil {
		return fmt.Errorf("failed to store ATS data: %w", err)
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

func SyncIndividualLoginDataToSessionStore(psid string, sessionKey string, signingKey string, privateId int32, user *sqlc.User) error {
	if err := SyncIndividualUserDataSessionStore(psid, user); err != nil {
		return err
	}
	if err := SyncIndividualUserDataSessionStore(user.PublicID.String(), user); err != nil {
		return err
	}

	newJwtKey := &sqlc.JwtKey{
		Key: signingKey,
	}
	if err := SyncIndividualJWTDataSessionStore(psid, newJwtKey); err != nil {
		return err
	}

	showcase, err := Queries.FindShowcaseRecordByUserId(context.Background(), privateId)
	if err != nil {
		log.Printf("Failed to find showcase record for login sync: %v", err)
	} else {
		if err := SyncIndividualShowCaseRecordDataSessionStore(psid, &showcase); err != nil {
			return err
		}
		if err := SyncIndividualShowCaseRecordDataSessionStore(user.PublicID.String(), &showcase); err != nil {
			return err
		}
	}

	resume, err := Queries.FindResumeByUserId(context.Background(), privateId)
	if err == nil {
		resumeData := map[string]any{
			"public_id":   resume.PublicID.Bytes,
			"template_id": resume.TemplateID,
			"detail":      resume.Detail,
		}
		serialised, _ := json.Marshal(resumeData)
		storeInValkey(psid+":resume_data", serialised)
		storeInValkey(user.PublicID.String()+":resume_data", serialised)
	}

	portfolio, err := Queries.FindPortfolioByUserId(context.Background(), privateId)
	if err == nil {
		portfolioData := map[string]any{
			"public_id":   portfolio.PublicID.Bytes,
			"template_id": portfolio.TemplateID,
			"detail":      portfolio.Detail,
		}
		serialised, _ := json.Marshal(portfolioData)
		storeInValkey(psid+":portfolio_data", serialised)
		storeInValkey(user.PublicID.String()+":portfolio_data", serialised)
	}

	ats, err := Queries.FindAtsByUserId(context.Background(), privateId)
	if err == nil {
		serialised, _ := json.Marshal(ats)
		storeInValkey(psid+":ats_data", serialised)
		storeInValkey(user.PublicID.String()+":ats_data", serialised)
	}

	if err := SyncIndividualUserSessionMapping(user.PublicID.String(), psid); err != nil {
		return err
	}

	if err := Queries.DeleteJwtKeyByUserId(context.Background(), privateId); err != nil {
		log.Printf("Failed to delete old JWT key: %v", err)
	}
	if err := Queries.DeleteSessionByUserId(context.Background(), privateId); err != nil {
		log.Printf("Failed to delete old session: %v", err)
	}

	newJwt, err := Queries.CreateJwtKey(context.Background(), sqlc.CreateJwtKeyParams{
		UserID: privateId,
		Key:    signingKey,
	})
	if err != nil {
		log.Printf("Failed to create new JWT key: %v", err)
	} else {
		newJwtKey = &newJwt
	}

	newSession, err := Queries.CreateSession(context.Background(), sqlc.CreateSessionParams{
		UserID:     privateId,
		SessionKey: sessionKey,
	})
	if err != nil {
		log.Printf("Failed to create new session: %v", err)
	} else {
		_ = newSession
	}

	return nil
}

func storeInValkey(key string, data []byte) {
	ctx := context.Background()
	_ = tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(key).Value(string(data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
}

func SyncIndividualResumeDataSessionStore(public_user_id string, resume *sqlc.Resume) error {
	resumeData := map[string]any{
		"public_id":   resume.PublicID.Bytes,
		"template_id": resume.TemplateID,
		"detail":      resume.Detail,
	}
	serialised, err := json.Marshal(resumeData)
	if err != nil {
		return fmt.Errorf("failed to serialise resume data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":resume_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
}

func SyncIndividualPortfolioDataSessionStore(public_user_id string, portfolio *sqlc.Portfolio) error {
	portfolioData := map[string]any{
		"public_id":   portfolio.PublicID.Bytes,
		"template_id": portfolio.TemplateID,
		"detail":      portfolio.Detail,
	}
	serialised, err := json.Marshal(portfolioData)
	if err != nil {
		return fmt.Errorf("failed to serialise portfolio data: %w", err)
	}

	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(public_user_id+":portfolio_data").Value(string(serialised)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
}

func SyncIndividualClientAuditLogSessionStore(public_user_id string, logs []sqlc.ClientAuditLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise client audit log data: %w", err)
	}
	return storeInValkeyWithTTL(public_user_id+":client_audit_log_data", serialised)
}

func SyncIndividualAdminAuditLogSessionStore(public_user_id string, logs []sqlc.AdminAuditLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise admin audit log data: %w", err)
	}
	return storeInValkeyWithTTL(public_user_id+":admin_audit_log_data", serialised)
}

func SyncIndividualClientReportLogSessionStore(public_user_id string, logs []sqlc.ClientReportLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise client report log data: %w", err)
	}
	return storeInValkeyWithTTL(public_user_id+":client_report_log_data", serialised)
}

func SyncIndividualErrorLogSessionStore(public_user_id string, logs []sqlc.ErrorLog) error {
	serialised, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to serialise error log data: %w", err)
	}
	return storeInValkeyWithTTL(public_user_id+":error_log_data", serialised)
}

func SyncIndividualClientSupportMessagingSessionStore(public_user_id string, messages []sqlc.ClientSupportMessaging) error {
	serialised, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to serialise support message data: %w", err)
	}
	return storeInValkeyWithTTL(public_user_id+":client_support_messages", serialised)
}

func storeInValkeyWithTTL(key string, data []byte) error {
	ctx := context.Background()
	return tool.Valkey.Do(
		ctx,
		tool.Valkey.B().Set().
			Key(key).Value(string(data)).
			Ex(systemconfig.SessionExpiryDuration).
			Build(),
	).Error()
}

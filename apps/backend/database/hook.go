package database

import (
	"gorm.io/gorm"
)

// user

func (u *User) AfterCreate(tx *gorm.DB) error {
	return tx.Create(&ShowcaseRecord{
		UserId: u.Id,
	}).Error
}

// showcaserecord

func (s *ShowcaseRecord) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (s *ShowcaseRecord) AfterDelete(tx *gorm.DB) error {
	return nil
}

// portfolio

func (p *Portfolio) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (p *Portfolio) AfterDelete(tx *gorm.DB) error {
	return nil
}

// resume

func (r *Resume) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (r *Resume) AfterDelete(tx *gorm.DB) error {
	return nil
}

// ats

func (a *Ats) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (a *Ats) AfterDelete(tx *gorm.DB) error {
	return nil
}

// jwt key

func (j *JwtKey) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (j *JwtKey) AfterDelete(tx *gorm.DB) error {
	return nil
}

// session

func (s *Session) AfterCreate(tx *gorm.DB) error {
	public_user_id := s.PublicId.String()

	var user User
	if err := tx.Preload("ShowcaseRecord").
		Preload("Resume").
		Preload("Portfolio").
		Preload("ATS").
		Preload("JWTKey").
		Where("id = ?", s.UserId).First(&user).Error; err != nil {
		return err
	}

	if err := SyncIndividualUserDataSessionStore(public_user_id, &user); err != nil {
		return err
	}
	if err := SyncIndividualJWTDataSessionStore(public_user_id, &user.JWTKey); err != nil {
		return err
	}
	if err := SyncIndividualShowCaseRecordDataSessionStore(public_user_id, &user); err != nil {
		return err
	}
	if err := SyncIndividualSessionDataSessionStore(public_user_id, s.SessionKey); err != nil {
		return err
	}
	if err := SyncIndividualResumeDataSessionStore(public_user_id, &user); err != nil {
		return err
	}
	if err := SyncIndividualPortfolioDataSessionStore(public_user_id, &user); err != nil {
		return err
	}
	if err := SyncIndividualATSDataSessionStore(public_user_id, &user); err != nil {
		return err
	}

	return nil
}

func (s *Session) AfterDelete(tx *gorm.DB) error {
	public_user_id := s.PublicId.String()
	if err := SyncIndividualUserDataDatabase(public_user_id); err != nil {
		return err
	}
	if err := SyncIndividualJWTDataDatabase(public_user_id); err != nil {
		return err
	}
	if err := SyncIndividualShowCaseRecordDatabase(public_user_id); err != nil {
		return err
	}
	if err := SyncIndividualSessionDataDatabase(public_user_id); err != nil {
		return err
	}
	if err := SyncIndividualResumeDataDatabase(public_user_id); err != nil {
		return err
	}
	if err := SyncIndividualPortfolioDataDatabase(public_user_id); err != nil {
		return err
	}
	if err := SyncIndividualATSDataDatabase(public_user_id); err != nil {
		return err
	}

	return nil
}

// client audit log

func (c *ClientAuditLog) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (c *ClientAuditLog) AfterDelete(tx *gorm.DB) error {
	return nil
}

// admin audit log

func (a *AdminAuditLog) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (a *AdminAuditLog) AfterDelete(tx *gorm.DB) error {
	return nil
}

// error log

func (e *ErrorLog) AfterCreate(tx *gorm.DB) error {
	return nil
}

func (e *ErrorLog) AfterDelete(tx *gorm.DB) error {
	return nil
}

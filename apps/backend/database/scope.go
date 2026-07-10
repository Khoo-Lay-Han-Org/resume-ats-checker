package database

import (
	"time"

	"gorm.io/gorm"
)

func NonExpiredUser(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredShowcaseRecord(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredPortfolio(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredResume(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredAts(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredJwtKey(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredSession(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredClientAuditLog(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredAdminAuditLog(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

func NonExpiredErrorLog(tx *gorm.DB) *gorm.DB {
	return tx.Where("expires_at IS NULL OR expires_at > ?", time.Now())
}

package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Capitalise the table and field names so that GORM can read and write them
type User struct {
	Id          int        `gorm:"type:int;primaryKey;index"`
	PublicId    uuid.UUID  `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	Username    string     `gorm:"type:varchar(255);not null;unique;"`
	Email       string     `gorm:"type:varchar(255);not null;unique;"`
	Password    []byte     `gorm:"type:bytea;not null;"`
	Displayname string     `gorm:"type:varchar(255);not null;unique;"`
	UserType    UserType   `gorm:"type:varchar(20);not null;"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt   time.Time  `gorm:"type:timestamptz;autoUpdateTime;index"`
	BannedAt    *time.Time `gorm:"type:timestamptz;index"`
	// GORM automatically detects this
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index"`
	ExpiresAt *time.Time     `gorm:"type:timestamptz;index"`

	ShowcaseRecord   ShowcaseRecord    `gorm:"foreignKey:UserId;references:Id"`
	Portfolio        Portfolio         `gorm:"foreignKey:UserId;references:Id"`
	Resume           Resume            `gorm:"foreignKey:UserId;references:Id"`
	ATS              Ats               `gorm:"foreignKey:UserId;references:Id"`
	JWTKey           JwtKey            `gorm:"foreignKey:UserId;references:Id"`
	Session          Session           `gorm:"foreignKey:UserId;references:Id"`
	ClientAuditLog   ClientAuditLog    `gorm:"foreignKey:UserId;references:Id"`
	AdminAuditLog    AdminAuditLog     `gorm:"foreignKey:UserId;references:Id"`
	ReportingReports []ClientReportLog `gorm:"foreignKey:ReportingUserId;references:Id"`
	TargetReports    []ClientReportLog `gorm:"foreignKey:TargetUserId;references:Id"`
	ErrorLog         ErrorLog          `gorm:"foreignKey:UserId;references:Id"`
}

type ShowcaseRecord struct {
	Id            int            `gorm:"type:int;primaryKey;"`
	PublicId      uuid.UUID      `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId        int            `gorm:"type:int;index;not null"`
	Name          pq.StringArray `gorm:"type:varchar(255)[];"`
	Email         pq.StringArray `gorm:"type:varchar(255)[];"`
	PhoneNumber   pq.StringArray `gorm:"type:varchar(255)[];"`
	Address       pq.StringArray `gorm:"type:varchar(255)[];"`
	SocialMedia   pq.StringArray `gorm:"type:varchar(255)[];"`
	JobExperience datatypes.JSON `gorm:"type:jsonb;"`
	Education     datatypes.JSON `gorm:"type:jsonb;"`
	Skill         pq.StringArray `gorm:"type:varchar(255)[];"`
	Certificate   datatypes.JSON `gorm:"type:jsonb;"`
	Language      pq.StringArray `gorm:"type:varchar(255)[];"`
	Project       datatypes.JSON `gorm:"type:jsonb;"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;autoUpdateTime;index"`
	DeletedAt     gorm.DeletedAt `gorm:"type:timestamptz;index"`
	ExpiresAt     *time.Time     `gorm:"type:timestamptz;index"`
}

type Portfolio struct {
	Id         int            `gorm:"type:int;primaryKey;"`
	PublicId   uuid.UUID      `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId     int            `gorm:"type:int;index;not null"`
	TemplateId int            `gorm:"type:int;index;not null"`
	Detail     datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt  time.Time      `gorm:"type:timestamptz;autoUpdateTime;index"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamptz;index"`
	ExpiresAt  *time.Time     `gorm:"type:timestamptz;index"`
}

type Resume struct {
	Id         int            `gorm:"type:int;primaryKey;"`
	PublicId   uuid.UUID      `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId     int            `gorm:"type:int;index;not null"`
	TemplateId int            `gorm:"type:int;index;not null"`
	Detail     datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt  time.Time      `gorm:"type:timestamptz;autoUpdateTime;index"`
	DeletedAt  gorm.DeletedAt `gorm:"type:timestamptz;index"`
	ExpiresAt  *time.Time     `gorm:"type:timestamptz;index"`
}

type Ats struct {
	Id        int            `gorm:"type:int;primaryKey;"`
	PublicId  uuid.UUID      `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId    int            `gorm:"type:int;index;not null"`
	Score     int            `gorm:"type:int;not null"`
	Reasoning string         `gorm:"type:varchar(355);not null"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime;index"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index"`
	ExpiresAt *time.Time     `gorm:"type:timestamptz;index"`
}

type JwtKey struct {
	Id        int        `gorm:"type:int;primaryKey;"`
	PublicId  uuid.UUID  `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId    int        `gorm:"type:int;index;not null;unique"`
	Key       string     `gorm:"type:varchar(355);not null;"`
	CreatedAt time.Time  `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;autoUpdateTime;index"`
	ExpiresAt *time.Time `gorm:"type:timestamptz;index"`
}

type Session struct {
	Id         int        `gorm:"type:int;primaryKey;"`
	PublicId   uuid.UUID  `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId     int        `gorm:"type:int;index;not null;unique"`
	SessionKey string     `gorm:"type:varchar(355);not null;"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt  time.Time  `gorm:"type:timestamptz;autoUpdateTime;index"`
	ExpiresAt  *time.Time `gorm:"type:timestamptz;index;default:CURRENT_TIMESTAMP + INTERVAL '3 days'"`
}

type ClientAuditLog struct {
	Id        int                `gorm:"type:int;primaryKey;"`
	PublicId  uuid.UUID          `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId    int                `gorm:"type:int;index;not null"`
	Type      ClientAuditLogType `gorm:"type:varchar(355);not null;"`
	Message   string             `gorm:"type:varchar(355);not null;"`
	CreatedAt time.Time          `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt time.Time          `gorm:"type:timestamptz;autoUpdateTime;index"`
	ExpiresAt *time.Time         `gorm:"type:timestamptz;index"`
}

type AdminAuditLog struct {
	Id        int               `gorm:"type:int;primaryKey;"`
	PublicId  uuid.UUID         `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId    int               `gorm:"type:int;index;not null"`
	Type      AdminAuditLogType `gorm:"type:varchar(355);not null;"`
	Message   string            `gorm:"type:varchar(355);not null;"`
	CreatedAt time.Time         `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt time.Time         `gorm:"type:timestamptz;autoUpdateTime;index"`
	ExpiresAt *time.Time        `gorm:"type:timestamptz;index"`
}

type ClientReportLog struct {
	Id              int                    `gorm:"type:int;primaryKey;"`
	PublicId        uuid.UUID              `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	ReportingUserId int                    `gorm:"type:int;index;not null"`
	TargetUserId    int                    `gorm:"type:int;index;not null"`
	Type            ClientBannedReasonType `gorm:"type:string;not null;"`
	CreatedAt       time.Time              `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt       time.Time              `gorm:"type:timestamptz;autoUpdateTime;index"`
	ExpiresAt       *time.Time             `gorm:"type:timestamptz;index"`
}

type ErrorLog struct {
	Id       int       `gorm:"type:int;primaryKey;"`
	PublicId uuid.UUID `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	UserId   int       `gorm:"type:int;index;not null"`
	// Use the status code
	Type      int        `gorm:"type:int;not null;"`
	Message   string     `gorm:"type:varchar(355);not null;"`
	CreatedAt time.Time  `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;autoUpdateTime;index"`
	ExpiresAt *time.Time `gorm:"type:timestamptz;index"`
}

type ClientSupportMessaging struct {
	Id        int                          `gorm:"type:int;primaryKey;"`
	PublicId  uuid.UUID                    `gorm:"type:uuid;not null;index;unique;default:gen_random_uuid()"`
	Type      ClientSupportMessagingTyping `gorm:"type:varchar(355);not null;"`
	Content   ClientSupportContentTyping   `gorm:"type:jsonb;not null;"`
	CreatedAt time.Time                    `gorm:"type:timestamptz;autoCreateTime;index"`
	UpdatedAt time.Time                    `gorm:"type:timestamptz;autoUpdateTime;index"`
	ExpiresAt *time.Time                   `gorm:"type:timestamptz;index"`
}

func TableConnect() error {
	if err := DB.AutoMigrate(
		&User{},
		&ShowcaseRecord{},
		&Portfolio{},
		&Resume{},
		&Ats{},
		&JwtKey{},
		&Session{},
		&ClientAuditLog{},
		&AdminAuditLog{},
		&ClientReportLog{},
		&ErrorLog{},
		&ClientSupportMessaging{},
	); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}
	return nil
}

package database

import "time"

type ClientAuditLogType string

const (
	// who, when
	NewClient ClientAuditLogType = "new client"
	// who, when
	PortfolioUpdate ClientAuditLogType = "portfolio update"
	// who, when
	ResumeUpdate ClientAuditLogType = "resume update"
	// who, when
	ShowcaseRecordUpdate ClientAuditLogType = "showcase record update"
	// who, when, old, new
	UsernameUpdate ClientAuditLogType = "username update"
	// who, when, old, new
	DisplaynameUpdate ClientAuditLogType = "displayname update"
	// who, when
	PasswordUpdate ClientAuditLogType = "password update"
	// who, when, old, new
	EmailUpdate ClientAuditLogType = "email update"
	// who, when
	AccountDeletion ClientAuditLogType = "account deletion"
)

type AdminAuditLogType string

const (
	// who, when
	NewAdmin AdminAuditLogType = "new admin"
	// who, when, title
	NewAnnouncement AdminAuditLogType = "new announcement"
	// who, when, to who, why
	ClientBanned AdminAuditLogType = "client banned"
	// who, when, to who, why
	AdminBanned AdminAuditLogType = "admin banned"
	// who, when, to who, satisfaction rate?
	CustomerSupported AdminAuditLogType = "customer supported"
)

type ClientBannedReasonType string

const (
	Profanity       ClientBannedReasonType = "profanity"
	ExplicitContent ClientBannedReasonType = "explicit content"
)

type AdminBannedReasonType string

const (
	PriviledgeAbuse        AdminBannedReasonType = "priviledge abuse"
	UnauthorisedDisclosure AdminBannedReasonType = "unauthorised disclosure"
)

type UserType string

const (
	SuperAdmin UserType = "super-admin"
	Admin      UserType = "admin"
	Client     UserType = "client"
)

type ClientSupportMessagingTyping string

const (
	TechnicalSupport    ClientSupportMessagingTyping = "technical support"
	FeatureImprovement  ClientSupportMessagingTyping = "feature improvement"
	BillingManagement   ClientSupportMessagingTyping = "billing management"
	ServiceAndOperation ClientSupportMessagingTyping = "service and operation"
	OnboardingSupport   ClientSupportMessagingTyping = "onboarding support"
	Complaint           ClientSupportMessagingTyping = "complaint"
)

type ClientSupportContentTyping struct {
	Text   string    `json:"text"`
	UserId int       `json:"user_id"`
	Time   time.Time `json:"time"`
}

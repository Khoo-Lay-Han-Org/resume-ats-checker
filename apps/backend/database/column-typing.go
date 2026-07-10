package database

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

type ClientCommType string

const (
	TechnicalSupport    ClientCommType = "technical support"
	FeatureImprovement  ClientCommType = "feature improvement"
	BillingManagement   ClientCommType = "billing management"
	ServiceAndOperation ClientCommType = "service and operation"
	OnboardingSupport   ClientCommType = "onboarding support"
	Complaint           ClientCommType = "complaint"
)

package database

import "time"

type InviteTokenDTO struct {
	PublicUserId string `json:"public_user_id"`
	Used         bool   `json:"used"`
}

type ClientAdminConfigDTO struct {
	PublicId     string     `json:"public_id"`
	Username     string     `json:"username"`
	Displayname  string     `json:"displayname"`
	PublicUserId string     `json:"public_user_id"`
	DeletedAt    *time.Time `json:"deleted_at"`
	BannedAt     *time.Time `json:"banned_at"`
}

type ClientReportLogDTO struct {
	PublicId              string `json:"public_id"`
	ReportingPublicUserId string `json:"reporting_public_user_id"`
	TargetPublicUserId    string `json:"target_public_user_id"`
	Type                  string `json:"type"`
}

type ClientCommDTO struct {
	UserId       int    `json:"user_id"`
	PublicId     string `json:"public_id"`
	PublicUserId string `json:"public_user_id"`
	Type         string `json:"type"`
	Message      string `json:"message"`
}

type AdminCommDTO struct {
	AdminCommPublicId     string `json:"admin_comm_public_id"`
	AdminUserPublicId     string `json:"admin_user_public_id"`
	ClientCommLogPublicId string `json:"client_comm_log_public_id"`
	Message               string `json:"message"`
}

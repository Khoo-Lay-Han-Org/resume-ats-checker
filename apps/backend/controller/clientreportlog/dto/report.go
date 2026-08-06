package clientreportlog_dto

type ClientReportRequest struct {
	TargetClientPublicUserId string `json:"target_client_public_user_id" binding:"required" mod:"trim" validate:"required"`
	ReportType               string `json:"report_type" binding:"required" mod:"trim" validate:"required"`
}

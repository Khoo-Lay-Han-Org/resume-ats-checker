package support_dto

type ClientReportRequest struct {
	TargetClientPublicUserId string `json:"target_client_public_user_id" binding:"required"`
	ReportType               string `json:"report_type" binding:"required"`
}

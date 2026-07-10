package setting_typing

type ChangeUsernameRequest struct {
	Username string `json:"username" binding:"required"`
}

type ChangeDisplaynameRequest struct {
	Displayname string `json:"displayname" binding:"required"`
}

type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

type OTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}

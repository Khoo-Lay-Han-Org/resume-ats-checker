package user_dto

type ChangeUsernameRequest struct {
	Username string `json:"username" binding:"required" mod:"trim" validate:"required,min=4,max=255"`
}

type ChangeDisplaynameRequest struct {
	Displayname string `json:"displayname" binding:"required" mod:"trim" validate:"required,min=4,max=255"`
}

type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required" mod:"trim,lcase" validate:"required,email,emailmx"`
}

type ChangePhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" mod:"trim" validate:"required"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required" validate:"required,min=8,max=20"`
}

type OTPRequest struct {
	OTP string `json:"otp" binding:"required"`
}

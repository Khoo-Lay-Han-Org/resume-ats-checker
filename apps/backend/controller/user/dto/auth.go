package user_dto

type Register struct {
	Displayname string `json:"displayname" binding:"required" mod:"trim" validate:"required,min=4,max=255"`
	Username    string `json:"username" binding:"required" mod:"trim" validate:"required,min=4,max=255"`
	Email       string `json:"email" binding:"required" mod:"trim,lcase" validate:"required,email,emailmx"`
	Password    string `json:"password" binding:"required" validate:"required,min=8,max=20"`
}

type Login struct {
	Email    string `json:"email" binding:"required" mod:"trim,lcase" validate:"required"`
	Password string `json:"password" binding:"required" mod:"trim" validate:"required"`
}

type OTP struct {
	OTP string `json:"otp" binding:"required"`
}

type OAuthResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

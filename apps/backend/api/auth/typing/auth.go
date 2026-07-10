package auth_typing

type Register struct {
	Displayname string `json:"displayname" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type OTP struct {
	OTP string `json:"otp" binding:"required"`
}

type OAuthResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

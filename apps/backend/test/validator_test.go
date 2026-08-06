package test

import (
	"testing"

	clientreportlog_dto "resuming/controller/clientreportlog/dto"
	clientreportlog_validator "resuming/controller/clientreportlog/validator"
	clientsupportmessage_dto "resuming/controller/clientsupportmessage/dto"
	clientsupportmessage_validator "resuming/controller/clientsupportmessage/validator"
	session_dto "resuming/controller/session/dto"
	session_validator "resuming/controller/session/validator"
	user_dto "resuming/controller/user/dto"
	user_validator "resuming/controller/user/validator"
)

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name    string
		input   user_dto.Register
		wantErr bool
	}{
		{
			name: "success",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			input: user_dto.Register{
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "username too short",
			input: user_dto.Register{
				Username:    "ab",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "missing displayname",
			input: user_dto.Register{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "displayname too short",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "ab",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "not-an-email",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "password too short",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "short",
			},
			wantErr: true,
		},
		{
			name: "password too long",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "thispasswordiswaytoolong20",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "all fields empty",
			input: user_dto.Register{
				Username:    "",
				Displayname: "",
				Email:       "",
				Password:    "",
			},
			wantErr: true,
		},
		{
			name: "email with leading/trailing spaces",
			input: user_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "  TEST@Example.com  ",
				Password:    "password123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user_validator.ValidateRegistration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegistration() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name    string
		input   user_dto.Login
		wantErr bool
	}{
		{
			name: "success",
			input: user_dto.Login{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			input: user_dto.Login{
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			input: user_dto.Login{
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: user_dto.Login{
				Email:    "",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "email with spaces trimmed",
			input: user_dto.Login{
				Email:    "  TEST@Example.com  ",
				Password: "password123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user_validator.ValidateLogin(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogin() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserControlRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   user_dto.UserControlRequest
		wantErr bool
	}{
		{
			name: "success",
			input: user_dto.UserControlRequest{
				PublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: false,
		},
		{
			name: "missing public user id",
			input: user_dto.UserControlRequest{
				PublicUserId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user_validator.ValidateUserControlRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserControlRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSessionControlRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   session_dto.SessionControlRequest
		wantErr bool
	}{
		{
			name: "success",
			input: session_dto.SessionControlRequest{
				PublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: false,
		},
		{
			name: "empty public user id",
			input: session_dto.SessionControlRequest{
				PublicUserId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := session_validator.ValidateSessionControlRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSessionControlRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateClientCommunicationReplyRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   clientsupportmessage_dto.ClientCommunicationReplyRequest
		wantErr bool
	}{
		{
			name: "success",
			input: clientsupportmessage_dto.ClientCommunicationReplyRequest{
				PublicId: "abc123",
				Message:  "Thank you for your message",
			},
			wantErr: false,
		},
		{
			name: "missing public id",
			input: clientsupportmessage_dto.ClientCommunicationReplyRequest{
				Message: "Thank you",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			input: clientsupportmessage_dto.ClientCommunicationReplyRequest{
				PublicId: "abc123",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: clientsupportmessage_dto.ClientCommunicationReplyRequest{
				PublicId: "",
				Message:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := clientsupportmessage_validator.ValidateClientCommunicationReplyRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClientCommunicationReplyRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateClientCommunicateRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   clientsupportmessage_dto.ClientCommunicateRequest
		wantErr bool
	}{
		{
			name: "success",
			input: clientsupportmessage_dto.ClientCommunicateRequest{
				Type:    "complaint",
				Message: "This is a complaint",
			},
			wantErr: false,
		},
		{
			name: "missing type",
			input: clientsupportmessage_dto.ClientCommunicateRequest{
				Message: "This is a complaint",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			input: clientsupportmessage_dto.ClientCommunicateRequest{
				Type: "complaint",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: clientsupportmessage_dto.ClientCommunicateRequest{
				Type:    "",
				Message: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := clientsupportmessage_validator.ValidateClientCommunicateRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClientCommunicateRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateClientReportRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   clientreportlog_dto.ClientReportRequest
		wantErr bool
	}{
		{
			name: "success",
			input: clientreportlog_dto.ClientReportRequest{
				TargetClientPublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
				ReportType:               "profanity",
			},
			wantErr: false,
		},
		{
			name: "missing target client",
			input: clientreportlog_dto.ClientReportRequest{
				ReportType: "profanity",
			},
			wantErr: true,
		},
		{
			name: "missing report type",
			input: clientreportlog_dto.ClientReportRequest{
				TargetClientPublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: clientreportlog_dto.ClientReportRequest{
				TargetClientPublicUserId: "",
				ReportType:               "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := clientreportlog_validator.ValidateClientReportRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClientReportRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUsernameRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   user_dto.ChangeUsernameRequest
		wantErr bool
	}{
		{
			name: "success",
			input: user_dto.ChangeUsernameRequest{
				Username: "newusername",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: user_dto.ChangeUsernameRequest{
				Username: "ab",
			},
			wantErr: true,
		},
		{
			name: "too long",
			input: user_dto.ChangeUsernameRequest{
				Username: string(make([]byte, 256)),
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: user_dto.ChangeUsernameRequest{
				Username: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user_validator.ValidateUsernameRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsernameRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDisplaynameRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   user_dto.ChangeDisplaynameRequest
		wantErr bool
	}{
		{
			name: "success",
			input: user_dto.ChangeDisplaynameRequest{
				Displayname: "New Display Name",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: user_dto.ChangeDisplaynameRequest{
				Displayname: "ab",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: user_dto.ChangeDisplaynameRequest{
				Displayname: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user_validator.ValidateDisplaynameRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDisplaynameRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmailRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   user_dto.ChangeEmailRequest
		wantErr bool
	}{
		{
			name: "success",
			input: user_dto.ChangeEmailRequest{
				Email: "newemail@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			input: user_dto.ChangeEmailRequest{
				Email: "not-an-email",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: user_dto.ChangeEmailRequest{
				Email: "",
			},
			wantErr: true,
		},
		{
			name: "email with spaces and uppercase",
			input: user_dto.ChangeEmailRequest{
				Email: "  NEW@Example.com  ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user_validator.ValidateEmailRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmailRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePasswordRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   user_dto.ChangePasswordRequest
		wantErr bool
	}{
		{
			name: "success",
			input: user_dto.ChangePasswordRequest{
				Password: "newpassword123",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: user_dto.ChangePasswordRequest{
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "too long",
			input: user_dto.ChangePasswordRequest{
				Password: "thispasswordiswaytoolong20",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: user_dto.ChangePasswordRequest{
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := user_validator.ValidatePasswordRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

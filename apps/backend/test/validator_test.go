package test

import (
	"testing"

	auth_validator "resuming/api/auth/validator"
	administrator_validator "resuming/api/administrator/validator"
	client_support_validator "resuming/api/client-support/validator"
	setting_validator "resuming/api/setting/validator"
	auth_typing "resuming/api/auth/typing"
	administrator_typing "resuming/api/administrator/typing"
	client_support_typing "resuming/api/client-support/typing"
	setting_typing "resuming/api/setting/typing"
)

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name    string
		input   auth_typing.Register
		wantErr bool
	}{
		{
			name: "success",
			input: auth_typing.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			input: auth_typing.Register{
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "username too short",
			input: auth_typing.Register{
				Username:    "ab",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "missing displayname",
			input: auth_typing.Register{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "displayname too short",
			input: auth_typing.Register{
				Username:    "testuser",
				Displayname: "ab",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			input: auth_typing.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "not-an-email",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			input: auth_typing.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "password too short",
			input: auth_typing.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "short",
			},
			wantErr: true,
		},
		{
			name: "password too long",
			input: auth_typing.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "thispasswordiswaytoolong20",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			input: auth_typing.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "all fields empty",
			input: auth_typing.Register{
				Username:    "",
				Displayname: "",
				Email:       "",
				Password:    "",
			},
			wantErr: true,
		},
		{
			name: "email with leading/trailing spaces",
			input: auth_typing.Register{
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
			_, err := auth_validator.ValidateRegistration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRegistration() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name    string
		input   auth_typing.Login
		wantErr bool
	}{
		{
			name: "success",
			input: auth_typing.Login{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			input: auth_typing.Login{
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			input: auth_typing.Login{
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: auth_typing.Login{
				Email:    "",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "email with spaces trimmed",
			input: auth_typing.Login{
				Email:    "  TEST@Example.com  ",
				Password: "password123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := auth_validator.ValidateLogin(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogin() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserControlRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   administrator_typing.UserControlRequest
		wantErr bool
	}{
		{
			name: "success",
			input: administrator_typing.UserControlRequest{
				PublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: false,
		},
		{
			name: "missing public user id",
			input: administrator_typing.UserControlRequest{
				PublicUserId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := administrator_validator.ValidateUserControlRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserControlRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateSessionControlRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   administrator_typing.SessionControlRequest
		wantErr bool
	}{
		{
			name: "success",
			input: administrator_typing.SessionControlRequest{
				PublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: false,
		},
		{
			name: "empty public user id",
			input: administrator_typing.SessionControlRequest{
				PublicUserId: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := administrator_validator.ValidateSessionControlRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSessionControlRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateClientCommunicationReplyRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   administrator_typing.ClientCommunicationReplyRequest
		wantErr bool
	}{
		{
			name: "success",
			input: administrator_typing.ClientCommunicationReplyRequest{
				PublicId: "abc123",
				Message:  "Thank you for your message",
			},
			wantErr: false,
		},
		{
			name: "missing public id",
			input: administrator_typing.ClientCommunicationReplyRequest{
				Message: "Thank you",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			input: administrator_typing.ClientCommunicationReplyRequest{
				PublicId: "abc123",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: administrator_typing.ClientCommunicationReplyRequest{
				PublicId: "",
				Message:  "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := administrator_validator.ValidateClientCommunicationReplyRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClientCommunicationReplyRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateClientCommunicateRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   client_support_typing.ClientCommunicateRequest
		wantErr bool
	}{
		{
			name: "success",
			input: client_support_typing.ClientCommunicateRequest{
				Type:    "complaint",
				Message: "This is a complaint",
			},
			wantErr: false,
		},
		{
			name: "missing type",
			input: client_support_typing.ClientCommunicateRequest{
				Message: "This is a complaint",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			input: client_support_typing.ClientCommunicateRequest{
				Type: "complaint",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: client_support_typing.ClientCommunicateRequest{
				Type:    "",
				Message: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client_support_validator.ValidateClientCommunicateRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClientCommunicateRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateClientReportRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   client_support_typing.ClientReportRequest
		wantErr bool
	}{
		{
			name: "success",
			input: client_support_typing.ClientReportRequest{
				TargetClientPublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
				ReportType:               "profanity",
			},
			wantErr: false,
		},
		{
			name: "missing target client",
			input: client_support_typing.ClientReportRequest{
				ReportType: "profanity",
			},
			wantErr: true,
		},
		{
			name: "missing report type",
			input: client_support_typing.ClientReportRequest{
				TargetClientPublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: client_support_typing.ClientReportRequest{
				TargetClientPublicUserId: "",
				ReportType:               "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client_support_validator.ValidateClientReportRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClientReportRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUsernameRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   setting_typing.ChangeUsernameRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_typing.ChangeUsernameRequest{
				Username: "newusername",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: setting_typing.ChangeUsernameRequest{
				Username: "ab",
			},
			wantErr: true,
		},
		{
			name: "too long",
			input: setting_typing.ChangeUsernameRequest{
				Username: string(make([]byte, 256)),
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_typing.ChangeUsernameRequest{
				Username: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := setting_validator.ValidateUsernameRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsernameRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateDisplaynameRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   setting_typing.ChangeDisplaynameRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_typing.ChangeDisplaynameRequest{
				Displayname: "New Display Name",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: setting_typing.ChangeDisplaynameRequest{
				Displayname: "ab",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_typing.ChangeDisplaynameRequest{
				Displayname: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := setting_validator.ValidateDisplaynameRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDisplaynameRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmailRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   setting_typing.ChangeEmailRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_typing.ChangeEmailRequest{
				Email: "newemail@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			input: setting_typing.ChangeEmailRequest{
				Email: "not-an-email",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_typing.ChangeEmailRequest{
				Email: "",
			},
			wantErr: true,
		},
		{
			name: "email with spaces and uppercase",
			input: setting_typing.ChangeEmailRequest{
				Email: "  NEW@Example.com  ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := setting_validator.ValidateEmailRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmailRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePasswordRequest(t *testing.T) {
	tests := []struct {
		name    string
		input   setting_typing.ChangePasswordRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_typing.ChangePasswordRequest{
				Password: "newpassword123",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: setting_typing.ChangePasswordRequest{
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "too long",
			input: setting_typing.ChangePasswordRequest{
				Password: "thispasswordiswaytoolong20",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_typing.ChangePasswordRequest{
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := setting_validator.ValidatePasswordRequest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordRequest() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

package test

import (
	"testing"

	administrator_dto "resuming/backend-api/administrator/dto"
	administrator_validator "resuming/backend-api/administrator/validator"
	auth_dto "resuming/backend-api/auth/dto"
	auth_validator "resuming/backend-api/auth/validator"
	client_support_dto "resuming/backend-api/client-support/dto"
	client_support_validator "resuming/backend-api/client-support/validator"
	setting_dto "resuming/backend-api/setting/dto"
	setting_validator "resuming/backend-api/setting/validator"
)

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name    string
		input   auth_dto.Register
		wantErr bool
	}{
		{
			name: "success",
			input: auth_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			input: auth_dto.Register{
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "username too short",
			input: auth_dto.Register{
				Username:    "ab",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "missing displayname",
			input: auth_dto.Register{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "displayname too short",
			input: auth_dto.Register{
				Username:    "testuser",
				Displayname: "ab",
				Email:       "test@example.com",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			input: auth_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "not-an-email",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			input: auth_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Password:    "password123",
			},
			wantErr: true,
		},
		{
			name: "password too short",
			input: auth_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "short",
			},
			wantErr: true,
		},
		{
			name: "password too long",
			input: auth_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
				Password:    "thispasswordiswaytoolong20",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			input: auth_dto.Register{
				Username:    "testuser",
				Displayname: "Test User",
				Email:       "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "all fields empty",
			input: auth_dto.Register{
				Username:    "",
				Displayname: "",
				Email:       "",
				Password:    "",
			},
			wantErr: true,
		},
		{
			name: "email with leading/trailing spaces",
			input: auth_dto.Register{
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
		input   auth_dto.Login
		wantErr bool
	}{
		{
			name: "success",
			input: auth_dto.Login{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "missing email",
			input: auth_dto.Login{
				Password: "password123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			input: auth_dto.Login{
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: auth_dto.Login{
				Email:    "",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "email with spaces trimmed",
			input: auth_dto.Login{
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
		input   administrator_dto.UserControlRequest
		wantErr bool
	}{
		{
			name: "success",
			input: administrator_dto.UserControlRequest{
				PublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: false,
		},
		{
			name: "missing public user id",
			input: administrator_dto.UserControlRequest{
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
		input   administrator_dto.SessionControlRequest
		wantErr bool
	}{
		{
			name: "success",
			input: administrator_dto.SessionControlRequest{
				PublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: false,
		},
		{
			name: "empty public user id",
			input: administrator_dto.SessionControlRequest{
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
		input   administrator_dto.ClientCommunicationReplyRequest
		wantErr bool
	}{
		{
			name: "success",
			input: administrator_dto.ClientCommunicationReplyRequest{
				PublicId: "abc123",
				Message:  "Thank you for your message",
			},
			wantErr: false,
		},
		{
			name: "missing public id",
			input: administrator_dto.ClientCommunicationReplyRequest{
				Message: "Thank you",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			input: administrator_dto.ClientCommunicationReplyRequest{
				PublicId: "abc123",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: administrator_dto.ClientCommunicationReplyRequest{
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
		input   client_support_dto.ClientCommunicateRequest
		wantErr bool
	}{
		{
			name: "success",
			input: client_support_dto.ClientCommunicateRequest{
				Type:    "complaint",
				Message: "This is a complaint",
			},
			wantErr: false,
		},
		{
			name: "missing type",
			input: client_support_dto.ClientCommunicateRequest{
				Message: "This is a complaint",
			},
			wantErr: true,
		},
		{
			name: "missing message",
			input: client_support_dto.ClientCommunicateRequest{
				Type: "complaint",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: client_support_dto.ClientCommunicateRequest{
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
		input   client_support_dto.ClientReportRequest
		wantErr bool
	}{
		{
			name: "success",
			input: client_support_dto.ClientReportRequest{
				TargetClientPublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
				ReportType:               "profanity",
			},
			wantErr: false,
		},
		{
			name: "missing target client",
			input: client_support_dto.ClientReportRequest{
				ReportType: "profanity",
			},
			wantErr: true,
		},
		{
			name: "missing report type",
			input: client_support_dto.ClientReportRequest{
				TargetClientPublicUserId: "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantErr: true,
		},
		{
			name: "both empty",
			input: client_support_dto.ClientReportRequest{
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
		input   setting_dto.ChangeUsernameRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_dto.ChangeUsernameRequest{
				Username: "newusername",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: setting_dto.ChangeUsernameRequest{
				Username: "ab",
			},
			wantErr: true,
		},
		{
			name: "too long",
			input: setting_dto.ChangeUsernameRequest{
				Username: string(make([]byte, 256)),
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_dto.ChangeUsernameRequest{
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
		input   setting_dto.ChangeDisplaynameRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_dto.ChangeDisplaynameRequest{
				Displayname: "New Display Name",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: setting_dto.ChangeDisplaynameRequest{
				Displayname: "ab",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_dto.ChangeDisplaynameRequest{
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
		input   setting_dto.ChangeEmailRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_dto.ChangeEmailRequest{
				Email: "newemail@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			input: setting_dto.ChangeEmailRequest{
				Email: "not-an-email",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_dto.ChangeEmailRequest{
				Email: "",
			},
			wantErr: true,
		},
		{
			name: "email with spaces and uppercase",
			input: setting_dto.ChangeEmailRequest{
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
		input   setting_dto.ChangePasswordRequest
		wantErr bool
	}{
		{
			name: "success",
			input: setting_dto.ChangePasswordRequest{
				Password: "newpassword123",
			},
			wantErr: false,
		},
		{
			name: "too short",
			input: setting_dto.ChangePasswordRequest{
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "too long",
			input: setting_dto.ChangePasswordRequest{
				Password: "thispasswordiswaytoolong20",
			},
			wantErr: true,
		},
		{
			name: "empty",
			input: setting_dto.ChangePasswordRequest{
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

package test

import (
	"net/http"
	"testing"
)

func TestSettingRoutes_RequireAuth(t *testing.T) {
	endpoints := []string{
		"/change-username",
		"/change-displayname",
		"/prepare-change-email",
		"/change-email",
		"/prepare-change-password",
		"/change-password",
		"/prepare-delete-account",
		"/delete-account",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			resp, err := makePostRequest(ep, map[string]string{})
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != http.StatusUnauthorized {
				t.Errorf("%s: status = %d, want 401, body = %v", ep, resp.statusCode, resp.body)
			}
		})
	}
}

func TestChangeUsername_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "username too short",
			body: map[string]string{
				"username": "ab",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "username empty",
			body: map[string]string{
				"username": "",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing username field",
			body:       map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makePostRequest("/change-username", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

func TestChangeDisplayname_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "displayname too short",
			body: map[string]string{
				"displayname": "ab",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "displayname empty",
			body: map[string]string{
				"displayname": "",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing displayname field",
			body:       map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makePostRequest("/change-displayname", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

func TestChangeEmail_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "invalid email format",
			body: map[string]string{
				"email": "not-an-email",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "email empty",
			body: map[string]string{
				"email": "",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing email field",
			body:       map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makePostRequest("/prepare-change-email", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

func TestChangePassword_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "password too short",
			body: map[string]string{
				"password": "short",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "password too long",
			body: map[string]string{
				"password": "thispasswordiswaytoolong20",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "password empty",
			body: map[string]string{
				"password": "",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing password field",
			body:       map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makePostRequest("/prepare-change-password", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

func TestDeleteAccount_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name:       "missing fields",
			body:       map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed JSON",
			body:       "not-json",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makePostRequest("/prepare-delete-account", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

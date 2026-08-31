package test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPrepareRegistration_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "missing required field",
			body: map[string]string{
				"displayname": "Test",
				"email":       "test@example.com",
				"password":    "password123",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "empty required field",
			body: map[string]string{
				"username":    "",
				"displayname": "Test",
				"email":       "test@example.com",
				"password":    "password123",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid email format",
			body: map[string]string{
				"username":    "testuser",
				"displayname": "Test",
				"email":       "not-an-email",
				"password":    "password123",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "password too short",
			body: map[string]string{
				"username":    "testuser",
				"displayname": "Test",
				"email":       "test@example.com",
				"password":    "short",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "malformed JSON",
			body:       "not-json",
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := MakePostRequest("/prepare-registeration/email", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

func TestRegisterFlow_RequiresCookie(t *testing.T) {
	tests := []struct {
		name string
		body any
		path string
	}{
		{
			name: "client register without cookie",
			body: map[string]string{"otp": "123456"},
			path: "/register",
		},
		{
			name: "admin register without cookie",
			body: map[string]string{"otp": "123456"},
			path: "/register-admin",
		},
		{
			name: "malformed JSON",
			body: "not-json",
			path: "/register",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := MakePostRequest(tt.path, tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want 400, body = %v", resp.statusCode, resp.body)
			}
		})
	}
}

func TestPrepareLogin_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "missing required field",
			body: map[string]string{
				"password": "password123",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "empty required field",
			body: map[string]string{
				"email":    "",
				"password": "password123",
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "malformed JSON",
			body:       "not-json",
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := MakePostRequest("/prepare-login/email", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

func TestLoginFlow_RequiresCookie(t *testing.T) {
	resp, err := MakePostRequest("/login", "not-json")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.statusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400, body = %v", resp.statusCode, resp.body)
	}
}

func TestAuthSuccessFlow(t *testing.T) {
	SkipIfNoInfra(t)

	email := fmt.Sprintf("testuser_%d@example.com", EpochMs())
	username := fmt.Sprintf("testuser%d", EpochMs())

	resp, err := MakePostRequest("/prepare-registeration/email", map[string]string{
		"username":    username,
		"displayname": "Test User",
		"email":       email,
		"password":    "password123",
	})
	if err != nil {
		t.Fatalf("prepare registration failed: %v", err)
	}
	if resp.statusCode != http.StatusOK {
		t.Fatalf("prepare registration status = %d, want 200, body = %v", resp.statusCode, resp.body)
	}
}

func TestDuplicateRegister(t *testing.T) {
	SkipIfNoInfra(t)

	email := fmt.Sprintf("dupe_%d@example.com", EpochMs())
	username := fmt.Sprintf("dupe%d", EpochMs())

	for i := 0; i < 2; i++ {
		resp, err := MakePostRequest("/prepare-registeration/email", map[string]string{
			"username":    username,
			"displayname": "Test User",
			"email":       email,
			"password":    "password123",
		})
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		if resp.statusCode != http.StatusOK {
			t.Errorf("request %d status = %d, want 200, body = %v", i, resp.statusCode, resp.body)
		}
	}
}

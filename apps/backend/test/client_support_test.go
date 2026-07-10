package test

import (
	"net/http"
	"testing"
)

func TestClientSupportRoutes_RequireAuth(t *testing.T) {
	endpoints := []string{
		"/client_comm_to_admin",
		"/client_report_other_client",
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

func TestClientCommunicate_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "missing type",
			body: map[string]string{
				"message": "Hello admin",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing message",
			body: map[string]string{
				"type": "complaint",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "both empty",
			body: map[string]string{
				"type":    "",
				"message": "",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "malformed JSON",
			body: "not-json",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makePostRequest("/client_comm_to_admin", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

func TestClientReport_Validation(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		wantStatus int
	}{
		{
			name: "missing target",
			body: map[string]string{
				"report_type": "profanity",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing report type",
			body: map[string]string{
				"target_client_public_user_id": "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "both empty",
			body: map[string]string{
				"target_client_public_user_id": "",
				"report_type":                  "",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "malformed JSON",
			body: "not-json",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := makePostRequest("/client_report_other_client", tt.body)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d, body = %v", resp.statusCode, tt.wantStatus, resp.body)
			}
		})
	}
}

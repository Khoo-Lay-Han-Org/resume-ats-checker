package test

import (
	"net/http"
	"testing"
)

func TestAdminPostRoutes_RequireAuth(t *testing.T) {
	postEndpoints := []string{
		"/ban_client",
		"/remove_individual_session",
		"/remove_all_session",
		"/client_comm_reply_log",
		"/remove_admin",
		"/invite-become-admin",
	}

	for _, ep := range postEndpoints {
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

func TestAdminGetRoutes_RequireAuth(t *testing.T) {
	getEndpoints := []string{
		"/client_comm_log",
		"/get_all_clients",
		"/get_all_admins",
		"/client_audit_logs",
		"/admin_audit_logs",
		"/error_audit_logs",
	}

	for _, ep := range getEndpoints {
		t.Run(ep, func(t *testing.T) {
			resp, err := makeRequest(http.MethodGet, ep, nil)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if resp.statusCode != http.StatusUnauthorized {
				t.Errorf("%s: status = %d, want 401, body = %v", ep, resp.statusCode, resp.body)
			}
		})
	}
}

func TestBanClient_RequiresAuth(t *testing.T) {
	resp, err := makePostRequest("/ban_client", map[string]string{
		"public_user_id": "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.statusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401, body = %v", resp.statusCode, resp.body)
	}
}

func TestRemoveIndividualSession_RequiresAuth(t *testing.T) {
	resp, err := makePostRequest("/remove_individual_session", map[string]string{
		"public_user_id": "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.statusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401, body = %v", resp.statusCode, resp.body)
	}
}

func TestRemoveAllSession_RequiresAuth(t *testing.T) {
	resp, err := makePostRequest("/remove_all_session", map[string]string{
		"public_user_id": "aaf33fc6-e1a1-4c95-946c-436dd68a7fbd",
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.statusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401, body = %v", resp.statusCode, resp.body)
	}
}

func TestClientCommReply_RequiresAuth(t *testing.T) {
	resp, err := makePostRequest("/client_comm_reply_log", map[string]string{
		"public_id": "abc123",
		"message":   "Thank you",
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.statusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401, body = %v", resp.statusCode, resp.body)
	}
}

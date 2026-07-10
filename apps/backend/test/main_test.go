package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"resuming/api"
	"resuming/database"
	"resuming/tool"
)

var router *gin.Engine
var infraAvailable bool

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		log.Printf("Failed to change working directory: %v", err)
	}
	if err := tool.SetupValkey(); err != nil {
		log.Printf("Valkey not available, skipping integration tests: %v", err)
		infraAvailable = false
	} else if err := database.DatabaseConnect(); err != nil {
		log.Printf("Database not available, skipping integration tests: %v", err)
		infraAvailable = false
	} else if err := database.TableConnect(); err != nil {
		log.Printf("Table migration failed, skipping integration tests: %v", err)
		infraAvailable = false
	} else {
		infraAvailable = true
	}

	gin.SetMode(gin.TestMode)
	router = api.APIConnect()

	os.Exit(m.Run())
}

func skipIfNoInfra(t *testing.T) {
	t.Helper()
	if !infraAvailable {
		t.Skip("Skipping: Valkey or PostgreSQL not available")
	}
}

type testResponse struct {
	statusCode int
	body       map[string]any
}

func makeRequest(method, path string, body any, cookies ...*http.Cookie) (*testResponse, error) {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var respBody map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		respBody = map[string]any{"raw": w.Body.String()}
	}

	return &testResponse{statusCode: w.Code, body: respBody}, nil
}

func makePostRequest(path string, body any, cookies ...*http.Cookie) (*testResponse, error) {
	return makeRequest(http.MethodPost, path, body, cookies...)
}

func epochMs() int64 {
	return time.Now().UnixMilli()
}

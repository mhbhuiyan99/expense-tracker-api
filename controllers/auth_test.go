package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestMain(m *testing.M) {
	// Move working directory up to project root so Beego can locate conf/app.conf
	os.Chdir("../")

	// Set unified testing variables used by ALL controller tests
	beego.AppConfig.Set("datadir", "testdata")
	beego.AppConfig.Set("userfile", "users_controller_test.csv")
	beego.AppConfig.Set("expensefile", "expenses_controller_test.csv")

	beego.InitBeegoBeforeIndex()

	code := m.Run()

	os.RemoveAll("testdata")
	os.Exit(code)
}


// resetTestData wipes out tracking artifacts between sub-tests
func resetTestData() {
	os.RemoveAll("testdata")
	os.MkdirAll("testdata", os.ModePerm)
}

func TestAuthController_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{} // Can take raw strings (bad JSON) or structs
		expectedStatus int
	}{
		{
			name: "success - valid registration",
			requestBody: map[string]string{
				"name":     "John Doe",
				"email":    "john@example.com",
				"password": "secretpassword",
			},
			expectedStatus: http.StatusCreated, // 201
		},
		{
			name: "failure - invalid json format",
			requestBody:    "{invalid-json-body",
			expectedStatus: http.StatusBadRequest, // 400
		},
		{
			name: "failure - duplicate email validation",
			requestBody: map[string]string{
				"name":     "Clone User",
				"email":    "john@example.com", // Already created in the first test item
				"password": "anotherpassword",
			},
			expectedStatus: http.StatusConflict, // 409
		},
		// NOTE: If your BaseController.ValidateInput contains your structural checks, 
		// these cases will catch them directly at the HTTP layer:
		{
			name: "failure - missing name validation",
			requestBody: map[string]string{
				"name":     "",
				"email":    "noname@example.com",
				"password": "secretpassword",
			},
			expectedStatus: http.StatusBadRequest, // 400
		},
	}

	resetTestData()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert payload variant interface to JSON byte streams
			var bodyBytes []byte
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			// Generate a mock HTTP lifecycle transaction target
			req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("Failed to create HTTP request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			// Dispatch the request execution directly into Beego's Router handler
			beego.BeeApp.Handlers.ServeHTTP(resp, req)

			// Validate response assertions
			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but received %d. Response: %s", tt.expectedStatus, resp.Code, resp.Body.String())
			}
		})
	}
}

func TestAuthController_Login(t *testing.T) {
	resetTestData()

	// Pre-seed a dummy user account through the Beego HTTP lifecycle
	seedBody, _ := json.Marshal(map[string]string{
		"name":     "Auth Test",
		"email":    "auth@example.com",
		"password": "securepassword123",
	})
	seedReq, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(seedBody))
	seedReq.Header.Set("Content-Type", "application/json")
	beego.BeeApp.Handlers.ServeHTTP(httptest.NewRecorder(), seedReq)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "success - valid login credentials",
			requestBody: map[string]string{
				"email":    "auth@example.com",
				"password": "securepassword123",
			},
			expectedStatus: http.StatusOK, // 200
			expectToken:    true,
		},
		{
			name: "failure - wrong password credential matching",
			requestBody: map[string]string{
				"email":    "auth@example.com",
				"password": "incorrectpassword",
			},
			expectedStatus: http.StatusUnauthorized, // 401
			expectToken:    false,
		},
		{
			name: "failure - nonexistent target email matching",
			requestBody: map[string]string{
				"email":    "missing@example.com",
				"password": "securepassword123",
			},
			expectedStatus: http.StatusUnauthorized, // 401
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but received %d. Response: %s", tt.expectedStatus, resp.Code, resp.Body.String())
				return
			}

			// If login succeeds, inspect the data payload layout to verify data structural integrity
			if tt.expectToken {
				var responseMap map[string]interface{}
				if err := json.Unmarshal(resp.Body.Bytes(), &responseMap); err != nil {
					t.Fatalf("Failed to parse response body structure: %v", err)
				}
				
				// Confirm that Beego returns your response wrapping envelope (e.g., success keys or data blocks)
				if responseMap["data"] == nil {
					t.Errorf("Expected populated data object block inside response payload, got: %v", resp.Body.String())
				}
			}
		})
	}
}
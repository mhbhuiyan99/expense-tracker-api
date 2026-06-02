package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

func TestMain(m *testing.M) {
	os.Chdir("../")

	beego.AppConfig.Set("datadir", "testdata")
	beego.AppConfig.Set("userfile", "users_controller_test.csv")
	beego.AppConfig.Set("expensefile", "expenses_controller_test.csv")

	// MANUALLY REGISTER ROUTES TO AVOID THE IMPORT CYCLE MATRIX
	// This mirrors what is inside your routers/router.go file!
	beego.Router("/api/v1/auth/register", &AuthController{}, "post:Register")
	beego.Router("/api/v1/auth/login", &AuthController{}, "post:Login")
	
	// Register summary BEFORE the parameter wildcard to prevent path conflicts
	beego.Router("/api/v1/expenses/summary", &ExpenseController{}, "get:Summary")
	beego.Router("/api/v1/expenses", &ExpenseController{}, "get:List;post:Create")
	beego.Router("/api/v1/expenses/:id", &ExpenseController{}, "get:GetOne;put:Update;delete:Remove")

	beego.InitBeegoBeforeIndex()

	code := m.Run()

	os.RemoveAll("testdata")
	os.Exit(code)
}

// resetExpenseTestData wipes out testing artifacts between sub-tests
func resetExpenseTestData() {
	os.RemoveAll("testdata")
	os.MkdirAll("testdata", os.ModePerm)
}

// helper to simulate a request through Beego's routing engine
func executeTestRequest(req *http.NewRequest) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(resp, req)
	return resp
}

func TestExpenseController_AuthRequired(t *testing.T) {
	resetExpenseTestData()

	// Verify that sending requests without the X-User-ID header yields a 401 Unauthorized
	req, _ := http.NewRequest("GET", "/api/v1/expenses", nil)
	resp := executeTestRequest(req)

	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for unauthenticated request, got %d", resp.Code)
	}
}

func TestExpenseController_Create(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name:   "success - valid expense creation",
			userID: "1",
			requestBody: map[string]interface{}{
				"title":        "Lunch",
				"amount":       350.50,
				"category":     "Food",
				"expense_date": "2026-06-10",
			},
			expectedStatus: http.StatusCreated, // 201
		},
		{
			name:           "failure - invalid json payload",
			userID:         "1",
			requestBody:    "{bad-json",
			expectedStatus: http.StatusBadRequest, // 400
		},
		{
			name:   "failure - service validation rejecting zero amount",
			userID: "1",
			requestBody: map[string]interface{}{
				"title":        "Free snack",
				"amount":       0,
				"category":     "Food",
				"expense_date": "2026-06-10",
			},
			expectedStatus: http.StatusBadRequest, // 400
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetExpenseTestData() // Start with a fresh mock CSV

			var bodyBytes []byte
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest("POST", "/api/v1/expenses", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", tt.userID)

			resp := executeTestRequest(req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, resp.Code, resp.Body.String())
			}
		})
	}
}

func TestExpenseController_CRUD_and_Filtering_Workflow(t *testing.T) {
	resetExpenseTestData()
	userID := "1"

	// 1. Seed two valid expense items via the API lifecycle
	exp1Body, _ := json.Marshal(map[string]interface{}{
		"title":        "Lunch",
		"amount":       150.00,
		"category":     "Food",
		"expense_date": "2026-06-01",
	})
	req1, _ := http.NewRequest("POST", "/api/v1/expenses", bytes.NewBuffer(exp1Body))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-User-ID", userID)
	resp1 := executeTestRequest(req1)
	
	// Unmarshal response to capture the generated dynamic ID
	var createdExpense map[string]interface{}
	var envelope map[string]interface{}
	json.Unmarshal(resp1.Body.Bytes(), &envelope)
	// Adapts cleanly if your RespondCreated / RespondSuccess wraps data inside a "data" property
	if dataField, exists := envelope["data"]; exists && dataField != nil {
		createdExpense = dataField.(map[string]interface{})
	} else {
		createdExpense = envelope
	}
	
	expenseID := int(createdExpense["id"].(float64))

	// Seed item 2 (Transport)
	exp2Body, _ := json.Marshal(map[string]interface{}{
		"title":        "Taxi",
		"amount":       500.00,
		"category":     "Transport",
		"expense_date": "2026-06-02",
	})
	req2, _ := http.NewRequest("POST", "/api/v1/expenses", bytes.NewBuffer(exp2Body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-User-ID", userID)
	executeTestRequest(req2)

	// 2. Test GET /api/v1/expenses (With filters and sorting)
	t.Run("List with filtering and sorting", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/expenses?category=Food&sort_by=amount&sort_order=desc", nil)
		req.Header.Set("X-User-ID", userID)
		resp := executeTestRequest(req)

		if resp.Code != http.StatusOK {
			t.Errorf("Expected GET listing to return 200, got %d", resp.Code)
		}
	})

	// 3. Test GET /api/v1/expenses/:id
	t.Run("Get single expense path validations", func(t *testing.T) {
		// Valid item retrieval
		reqValid, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/expenses/%d", expenseID), nil)
		reqValid.Header.Set("X-User-ID", userID)
		respValid := executeTestRequest(reqValid)
		if respValid.Code != http.StatusOK {
			t.Errorf("Expected 200 for valid expense item read, got %d", respValid.Code)
		}

		// Non-existent item (404)
		reqMissing, _ := http.NewRequest("GET", "/api/v1/expenses/99999", nil)
		reqMissing.Header.Set("X-User-ID", userID)
		respMissing := executeTestRequest(reqMissing)
		if respMissing.Code != http.StatusNotFound {
			t.Errorf("Expected 404 for missing resource identifier, got %d", respMissing.Code)
		}

		// Wrong User attempting to get item (Should return 404/403 according to spec ownership protections)
		reqWrongUser, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/expenses/%d", expenseID), nil)
		reqWrongUser.Header.Set("X-User-ID", "2") // Alternative User Context
		respWrongUser := executeTestRequest(reqWrongUser)
		if respWrongUser.Code != http.StatusNotFound {
			t.Errorf("Expected 404 ownership rejection when cross-accessing resource, got %d", respWrongUser.Code)
		}
	})

	// 4. Test PUT /api/v1/expenses/:id
	t.Run("Update expense path validation", func(t *testing.T) {
		updateTitle := "Business Team Lunch"
		updateBody, _ := json.Marshal(map[string]interface{}{
			"title": &updateTitle,
		})
		
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/expenses/%d", expenseID), bytes.NewBuffer(updateBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-ID", userID)
		
		resp := executeTestRequest(req)
		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200 for successful resource path alteration, got %d. Body: %s", resp.Code, resp.Body.String())
		}
	})

	// 5. Test GET /api/v1/expenses/summary
	t.Run("Get expenses summary statistics", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/expenses/summary?date_from=2026-06-01&date_to=2026-06-30", nil)
		req.Header.Set("X-User-ID", userID)
		
		resp := executeTestRequest(req)
		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200 for summary collection calculation, got %d. Body: %s", resp.Code, resp.Body.String())
		}
	})

	// 6. Test DELETE /api/v1/expenses/:id
	t.Run("Delete expense path cleanup verification", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/expenses/%d", expenseID), nil)
		req.Header.Set("X-User-ID", userID)
		
		resp := executeTestRequest(req)
		if resp.Code != http.StatusOK {
			t.Errorf("Expected 200 status verification on entry delete, got %d", resp.Code)
		}

		// Ensure consecutive execution returns a 404 resource missing error status code
		respSecondTry := executeTestRequest(req)
		if respSecondTry.Code != http.StatusNotFound {
			t.Errorf("Expected 404 validation response on duplicate delete executions, got %d", respSecondTry.Code)
		}
	})
}

func TestExpenseController_InvalidRouteID(t *testing.T) {
	resetExpenseTestData()
	
	// Validates that alphanumeric wildcards fail cleanly inside conversion functions
	req, _ := http.NewRequest("GET", "/api/v1/expenses/abc-invalid-id", nil)
	req.Header.Set("X-User-ID", "1")
	resp := executeTestRequest(req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for alphabetical resource string parameters, got %d", resp.Code)
	}
}
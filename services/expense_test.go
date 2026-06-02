package services

import (
	"testing"

	"expense-tracker-api/models"
)

// mustCreateExpense creates an expense and fails the test immediately if it errors.
func mustCreateExpense(t *testing.T, userID int, req models.CreateExpenseRequest) *models.Expense {
	t.Helper()
	expense, err := CreateExpense(userID, req)
	if err != nil {
		t.Fatalf("setup: failed to create expense %q: %v", req.Title, err)
	}
	return expense
}

func setupExpenseTest(t *testing.T) int {
	cleanTestData()
	// register a user and return their ID
	RegisterUser(models.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "secret123",
	})
	data, err := LoginUser(models.LoginRequest{
		Email:    "test@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	return data.UserID
}

func TestCreateExpense(t *testing.T) {
	tests := []struct {
		name    string
		input   models.CreateExpenseRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid expense",
			input: models.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			wantErr: false,
		},
		{
			name: "missing title",
			input: models.CreateExpenseRequest{
				Amount:      350.50,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "zero amount",
			input: models.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      0,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			wantErr: true,
			errMsg:  "amount must be greater than zero",
		},
		{
			name: "negative amount",
			input: models.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      -100,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			wantErr: true,
			errMsg:  "amount must be greater than zero",
		},
		{
			name: "invalid category",
			input: models.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "InvalidCat",
				ExpenseDate: "2025-06-10",
			},
			wantErr: true,
			errMsg:  "invalid category",
		},
		{
			name: "missing date",
			input: models.CreateExpenseRequest{
				Title:    "Lunch",
				Amount:   350.50,
				Category: "Food",
			},
			wantErr: true,
			errMsg:  "expense_date is required",
		},
	}

	userID := setupExpenseTest(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expense, err := CreateExpense(userID, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}
			if expense.ID == 0 {
				t.Errorf("expected valid expense ID, got 0")
			}
		})
	}
}

func TestGetExpenses(t *testing.T) {
	userID := setupExpenseTest(t)

	mustCreateExpense(t, userID, models.CreateExpenseRequest{
		Title: "Lunch", Amount: 350.50,
		Category: "Food", ExpenseDate: "2025-06-10",
	})
	mustCreateExpense(t, userID, models.CreateExpenseRequest{
		Title: "Bus", Amount: 50.00,
		Category: "Transport", ExpenseDate: "2025-06-11",
	})
	mustCreateExpense(t, userID, models.CreateExpenseRequest{
		Title: "Rent", Amount: 15000.00,
		Category: "Housing", ExpenseDate: "2025-07-01",
	})

	tests := []struct {
		name      string
		params    ListExpensesParams
		wantCount int
	}{
		{
			name:      "no filters returns all",
			params:    ListExpensesParams{},
			wantCount: 3,
		},
		{
			name:      "filter by category Food",
			params:    ListExpensesParams{Category: "Food"},
			wantCount: 1,
		},
		{
			name:      "filter by date range June only",
			params:    ListExpensesParams{DateFrom: "2025-06-01", DateTo: "2025-06-30"},
			wantCount: 2,
		},
		{
			name:      "sort by amount desc",
			params:    ListExpensesParams{SortBy: "amount", SortOrder: "desc"},
			wantCount: 3,
		},
		{
			name:      "sort by amount asc",
			params:    ListExpensesParams{SortBy: "amount", SortOrder: "asc"},
			wantCount: 3,
		},
		{
			name:      "limit 1",
			params:    ListExpensesParams{Limit: 1},
			wantCount: 1,
		},
		{
			name:      "category not found returns empty",
			params:    ListExpensesParams{Category: "Healthcare"},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenses, err := GetExpenses(userID, tt.params)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(expenses) != tt.wantCount {
				t.Errorf("expected %d expenses, got %d", tt.wantCount, len(expenses))
			}
		})
	}
}

func TestGetSummary(t *testing.T) {
	userID := setupExpenseTest(t)

	mustCreateExpense(t, userID, models.CreateExpenseRequest{
		Title: "Lunch", Amount: 350.50,
		Category: "Food", ExpenseDate: "2025-06-10",
	})
	mustCreateExpense(t, userID, models.CreateExpenseRequest{
		Title: "Bus", Amount: 50.00,
		Category: "Transport", ExpenseDate: "2025-06-11",
	})

	tests := []struct {
		name        string
		dateFrom    string
		dateTo      string
		wantErr     bool
		errMsg      string
		wantTotal   float64
		wantCount   int
	}{
		{
			name:      "valid summary",
			dateFrom:  "2025-06-01",
			dateTo:    "2025-06-30",
			wantTotal: 400.50,
			wantCount: 2,
		},
		{
			name:    "missing date_from",
			dateTo:  "2025-06-30",
			wantErr: true,
			errMsg:  "date_from is required",
		},
		{
			name:     "missing date_to",
			dateFrom: "2025-06-01",
			wantErr:  true,
			errMsg:   "date_to is required",
		},
		{
			name:     "date_from after date_to",
			dateFrom: "2025-06-30",
			dateTo:   "2025-06-01",
			wantErr:  true,
			errMsg:   "date_from cannot be after date_to",
		},
		{
			name:      "empty range returns zero",
			dateFrom:  "2025-01-01",
			dateTo:    "2025-01-31",
			wantTotal: 0,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetSummary(userID, tt.dateFrom, tt.dateTo)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result.TotalAmount != tt.wantTotal {
				t.Errorf("expected total %v, got %v", tt.wantTotal, result.TotalAmount)
			}
			if result.TotalCount != tt.wantCount {
				t.Errorf("expected count %d, got %d", tt.wantCount, result.TotalCount)
			}
		})
	}
}

func TestGetExpense(t *testing.T) {
	userID := setupExpenseTest(t)
	created := mustCreateExpense(t, userID, models.CreateExpenseRequest{
		Title: "Lunch", Amount: 350.50,
		Category: "Food", ExpenseDate: "2025-06-10",
	})

	tests := []struct {
		name    string
		id      int
		userID  int
		wantErr bool
	}{
		{
			name:    "valid get",
			id:      created.ID,
			userID:  userID,
			wantErr: false,
		},
		{
			name:    "wrong user cannot get",
			id:      created.ID,
			userID:  999,
			wantErr: true,
		},
		{
			name:    "id does not exist",
			id:      999,
			userID:  userID,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expense, err := GetExpense(tt.id, tt.userID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if expense.ID != tt.id {
				t.Errorf("expected expense ID %d, got %d", tt.id, expense.ID)
			}
		})
	}
}

func TestUpdateExpense(t *testing.T) {
	userID := setupExpenseTest(t)
	created := mustCreateExpense(t, userID, models.CreateExpenseRequest{
		Title: "Lunch", Amount: 350.50,
		Category: "Food", ExpenseDate: "2025-06-10",
	})

	newTitle := "Team Lunch"
	newAmount := 500.00
	newCategory := "Food"
	wrongCategory := "InvalidCat"

	tests := []struct {
		name    string
		id      int
		userID  int
		req     models.UpdateExpenseRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid update title",
			id:      created.ID,
			userID:  userID,
			req:     models.UpdateExpenseRequest{Title: &newTitle},
			wantErr: false,
		},
		{
			name:    "valid update amount",
			id:      created.ID,
			userID:  userID,
			req:     models.UpdateExpenseRequest{Amount: &newAmount},
			wantErr: false,
		},
		{
			name:    "invalid category",
			id:      created.ID,
			userID:  userID,
			req:     models.UpdateExpenseRequest{Category: &wrongCategory},
			wantErr: true,
			errMsg:  "invalid category",
		},
		{
			name:    "valid category update",
			id:      created.ID,
			userID:  userID,
			req:     models.UpdateExpenseRequest{Category: &newCategory},
			wantErr: false,
		},
		{
			name:   "negative amount",
			id:     created.ID,
			userID: userID,
			req: models.UpdateExpenseRequest{
				Amount: func() *float64 { v := -100.0; return &v }(),
			},
			wantErr: true,
			errMsg:  "amount must be greater than zero",
		},
		{
			name:    "expense not found",
			id:      999,
			userID:  userID,
			req:     models.UpdateExpenseRequest{Title: &newTitle},
			wantErr: true,
			errMsg:  "expense not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expense, err := UpdateExpense(tt.id, tt.userID, tt.req)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if expense == nil {
				t.Errorf("expected expense, got nil")
			}
		})
	}
}

func TestDeleteExpense(t *testing.T) {
	userID := setupExpenseTest(t)

	expense, _ := CreateExpense(userID, models.CreateExpenseRequest{
		Title: "Lunch", Amount: 350.50,
		Category: "Food", ExpenseDate: "2025-06-10",
	})

	tests := []struct {
		name    string
		id      int
		userID  int
		wantErr bool
	}{
		{
			name:    "valid delete",
			id:      expense.ID,
			userID:  userID,
			wantErr: false,
		},
		{
			name:    "already deleted",
			id:      expense.ID,
			userID:  userID,
			wantErr: true,
		},
		{
			name:    "wrong user cannot delete",
			id:      999,
			userID:  99,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteExpense(tt.id, tt.userID)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
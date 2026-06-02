package models

import (
	"testing"
)

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{
			name:     "valid Food category",
			category: "Food",
			want:     true,
		},
		{
			name:     "valid Transport category",
			category: "Transport",
			want:     true,
		},
		{
			name:     "valid Housing category",
			category: "Housing",
			want:     true,
		},
		{
			name:     "valid Entertainment category",
			category: "Entertainment",
			want:     true,
		},
		{
			name:     "valid Shopping category",
			category: "Shopping",
			want:     true,
		},
		{
			name:     "valid Healthcare category",
			category: "Healthcare",
			want:     true,
		},
		{
			name:     "valid Education category",
			category: "Education",
			want:     true,
		},
		{
			name:     "valid Utilities category",
			category: "Utilities",
			want:     true,
		},
		{
			name:     "valid Other category",
			category: "Other",
			want:     true,
		},
		{
			name:     "invalid lowercase food",
			category: "food",
			want:     false,
		},
		{
			name:     "invalid random category",
			category: "Random",
			want:     false,
		},
		{
			name:     "invalid empty string",
			category: "",
			want:     false,
		},
		{
			name:     "invalid partial match",
			category: "Foo",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidCategory(tt.category)
			if got != tt.want {
				t.Errorf("IsValidCategory(%q) = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}

func TestAllowedCategoriesLength(t *testing.T) {
	want := 9
	got := len(AllowedCategories)
	if got != want {
		t.Errorf("len(AllowedCategories) = %d, want %d", got, want)
	}
}

func TestAllowedCategoriesContent(t *testing.T) {
	expectedCategories := map[string]bool{
		"Food":          true,
		"Transport":     true,
		"Housing":       true,
		"Entertainment": true,
		"Shopping":      true,
		"Healthcare":    true,
		"Education":     true,
		"Utilities":     true,
		"Other":         true,
	}

	for _, category := range AllowedCategories {
		if !expectedCategories[category] {
			t.Errorf("unexpected category in AllowedCategories: %q", category)
		}
	}

	if len(AllowedCategories) != len(expectedCategories) {
		t.Errorf("AllowedCategories has %d items, expected %d", len(AllowedCategories), len(expectedCategories))
	}
}

func TestExpenseStruct(t *testing.T) {
	tests := []struct {
		name     string
		expense  Expense
		checkID  bool
		wantID   int
		checkAll bool
	}{
		{
			name: "complete expense struct",
			expense: Expense{
				ID:          1,
				UserID:      100,
				Title:       "Lunch",
				Amount:      25.50,
				Category:    "Food",
				Note:        "Office lunch",
				ExpenseDate: "2025-06-10",
				CreatedAt:   "2025-06-10T12:00:00Z",
			},
			checkID:  true,
			wantID:   1,
			checkAll: true,
		},
		{
			name: "expense with zero ID",
			expense: Expense{
				ID:          0,
				UserID:      100,
				Title:       "Test",
				Amount:      10.00,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			checkID: true,
			wantID:  0,
		},
		{
			name: "expense with empty note",
			expense: Expense{
				ID:          1,
				UserID:      100,
				Title:       "Test",
				Amount:      10.00,
				Category:    "Food",
				Note:        "",
				ExpenseDate: "2025-06-10",
			},
			checkAll: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.checkID && tt.expense.ID != tt.wantID {
				t.Errorf("Expense.ID = %d, want %d", tt.expense.ID, tt.wantID)
			}
			if tt.checkAll {
				if tt.expense.Title == "" {
					t.Error("Expense.Title is empty")
				}
				if tt.expense.Amount <= 0 {
					t.Error("Expense.Amount is not positive")
				}
				if tt.expense.Category == "" {
					t.Error("Expense.Category is empty")
				}
			}
		})
	}
}

func TestCreateExpenseRequest(t *testing.T) {
	tests := []struct {
		name    string
		request CreateExpenseRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "Food",
				Note:        "Office lunch",
				ExpenseDate: "2025-06-10",
			},
			wantErr: false,
		},
		{
			name: "request with empty note",
			request: CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "Food",
				Note:        "",
				ExpenseDate: "2025-06-10",
			},
			wantErr: false,
		},
		{
			name: "request missing title",
			request: CreateExpenseRequest{
				Title:       "",
				Amount:      350.50,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			wantErr: true,
		},
		{
			name: "request with zero amount",
			request: CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      0,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate basic fields
			hasError := tt.request.Title == "" || tt.request.Amount <= 0

			if tt.wantErr && !hasError {
				t.Errorf("expected error but validation passed")
			}
			if !tt.wantErr && hasError {
				t.Errorf("expected no error but validation failed")
			}
		})
	}
}

func TestUpdateExpenseRequest(t *testing.T) {
	tests := []struct {
		name     string
		newTitle *string
		newAmount *float64
		newCat   *string
	}{
		{
			name:      "update title only",
			newTitle:  func() *string { s := "New Title"; return &s }(),
			newAmount: nil,
			newCat:    nil,
		},
		{
			name:      "update amount only",
			newTitle:  nil,
			newAmount: func() *float64 { f := 500.00; return &f }(),
			newCat:    nil,
		},
		{
			name:      "update all fields",
			newTitle:  func() *string { s := "New"; return &s }(),
			newAmount: func() *float64 { f := 100.0; return &f }(),
			newCat:    func() *string { s := "Transport"; return &s }(),
		},
		{
			name:      "all nil",
			newTitle:  nil,
			newAmount: nil,
			newCat:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := UpdateExpenseRequest{
				Title:  tt.newTitle,
				Amount: tt.newAmount,
				Category: tt.newCat,
			}

			if tt.newTitle != nil {
				if *req.Title != *tt.newTitle {
					t.Errorf("Title mismatch")
				}
			} else {
				if req.Title != nil {
					t.Error("expected nil Title")
				}
			}

			if tt.newAmount != nil {
				if *req.Amount != *tt.newAmount {
					t.Errorf("Amount mismatch")
				}
			} else {
				if req.Amount != nil {
					t.Error("expected nil Amount")
				}
			}

			if tt.newCat != nil {
				if *req.Category != *tt.newCat {
					t.Errorf("Category mismatch")
				}
			} else {
				if req.Category != nil {
					t.Error("expected nil Category")
				}
			}
		})
	}
}

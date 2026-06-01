package models

// Expense represents a single expense record
type Expense struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Title string `json:"title"`
	Amount float64 `json:"amount"`
	Category string `json:"category"`
	Note string `json:"note,omitempty"`
	ExpenseDate string `json:"expense_date"`
	CreatedAt string `json:"created_at"`
}

// AllowedCategories defines the valid categories for expenses
var AllowedCategories = []string{
	"Food",
	"Transportation",
	"Entertainment",
	"Utilities",
	"Healthcare",
	"Education",
	"Housing",
	"Shopping",
	"Other",
}

// IsValidCategory checks if the provided category is valid
func IsValidCategory(category string) bool {
	for _, c := range AllowedCategories {
		if c == category {
			return true
		}
	}
	return false
}

// CreateExpenseRequest is the expected JSON body for creating a new expense
type CreateExpenseRequest struct {
	Title       string  `json:"title" valid:"Required"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category" valid:"Required"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date" valid:"Required;Date"`
}

// UpdateExpenseRequest is the expected JSON body for updating an existing expense
type UpdateExpenseRequest struct {
	Title       *string  `json:"title"`
	Amount      *float64 `json:"amount"`
	Category    *string  `json:"category"`
	Note        *string  `json:"note"`
	ExpenseDate *string  `json:"expense_date"`
}
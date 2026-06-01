package services

import(
	"expense-tracker-api/models"
	"fmt"
	"sort"
)

type ListExpensesParams struct {
	Category string
	DateFrom string
	DateTo string
	SortBy string
	SortOrder string
	Limit int
}

// CreateExpense validates and creates a new expense for a user
func CreateExpense(userID int, req models.CreateExpenseRequest) (*models.Expense, error) {
	if req.Title == "" {
        return nil, fmt.Errorf("title is required")
    }
    if req.Amount <= 0 {
        return nil, fmt.Errorf("amount must be greater than zero")
    }
    if req.ExpenseDate == "" {
        return nil, fmt.Errorf("expense_date is required")
    }
    if !models.IsValidCategory(req.Category) {
        return nil, fmt.Errorf("invalid category")
    }

	expense := &models.Expense{
		UserID: userID,
		Title: req.Title,
		Amount: req.Amount,
		Category: req.Category,
		Note: req.Note,
		ExpenseDate: req.ExpenseDate,
	}

	if err := models.CreateExpense(expense); err != nil {
		return nil, fmt.Errorf("failed to save expense: %w", err)
	}
	return expense, nil
}

func GetExpenses(userID int, params ListExpensesParams) ([]models.Expense, error) {
	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Filter by category
	if params.Category != "" {
		filtered := make([]models.Expense, 0)
		for _, e := range expenses {
			if e.Category == params.Category {
				filtered = append(filtered, e)
			}
		}
		expenses = filtered
	}

	// Filter by date range
	if params.DateFrom != "" {
		filtered := make([]models.Expense, 0)
		for _, e := range expenses {
			if e.ExpenseDate >= params.DateFrom {
				filtered = append(filtered, e)
			}
		}
		expenses = filtered
	}

	if params.DateTo != "" {
		filtered := make([]models.Expense, 0)
		for _, e := range expenses {
			if e.ExpenseDate <= params.DateTo {
				filtered = append(filtered, e)
			}
		}
		expenses = filtered
	}

	// Sort expenses
	if params.SortBy != "" {
		sort.Slice(expenses, func(i, j int) bool {
			asc := params.SortOrder != "desc"
			switch params.SortBy {
			case "amount":
				if asc {
					return expenses[i].Amount < expenses[j].Amount
				}
				return expenses[i].Amount > expenses[j].Amount
			case "expense_date":
				if asc {
					return expenses[i].ExpenseDate < expenses[j].ExpenseDate
				}
				return expenses[i].ExpenseDate > expenses[j].ExpenseDate
			}
			return true
		})
	} else {
		// Default sort by newest first
		sort.Slice(expenses, func(i, j int) bool {
			return expenses[i].ExpenseDate > expenses[j].ExpenseDate
		})
	}

	// Paginate results
	if params.Limit > 0 && len(expenses) > params.Limit {
		expenses = expenses[:params.Limit]
	}

	return expenses, nil
}

// GetExpense returns a single expense by ID for a user
func GetExpense(id int, userID int) (*models.Expense, error) {
	expense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		return nil, err
	} 
	if expense == nil {
		return nil, fmt.Errorf("expense not found")
	}
	return expense, nil
}

// UpdateExpense updates an existing expense
func UpdateExpense(id int, userID int, req models.UpdateExpenseRequest) (*models.Expense, error) {
	expense, err := GetExpense(id, userID)
	if err != nil {
		return nil, err
	}

	if expense == nil {
		return nil, fmt.Errorf("expense not found")
	}

	if req.Title != nil {
		expense.Title = *req.Title
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return nil, fmt.Errorf("amount must be greater than zero")
		}
		expense.Amount = *req.Amount
	}
	if req.Category != nil {
		if !models.IsValidCategory(*req.Category) {
			return nil, fmt.Errorf("invalid category")
		}
		expense.Category = *req.Category
	}
	if req.Note != nil {
		expense.Note = *req.Note
	}
	if req.ExpenseDate != nil {
		expense.ExpenseDate = *req.ExpenseDate
	}

	if err := models.UpdateExpense(expense); err != nil {
		return nil, err
	}
	return expense, nil
}

// DeleteExpense deletes an expense by ID for the user
func DeleteExpense(id int, userID int) error {
	expense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		return err
	}
	if expense == nil {
		return fmt.Errorf("expense not found")
	}
	return models.DeleteExpense(id, userID)
}

// CategorySummary holds total and count for one category
type CategorySummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

// SummaryResult holds the full summary response data
type SummaryResult struct {
	DateFrom    string            `json:"date_from"`
	DateTo      string            `json:"date_to"`
	TotalAmount float64           `json:"total_amount"`
	TotalCount  int               `json:"total_count"`
	ByCategory  []CategorySummary `json:"by_category"`
}

// GetSummary returns a spending summary for a user within a date range.
func GetSummary(userID int, dateFrom, dateTo string) (*SummaryResult, error) {
	if dateFrom == "" {
		return nil, fmt.Errorf("date_from is required")
	}
	if dateTo == "" {
		return nil, fmt.Errorf("date_to is required")
	}
	if dateFrom > dateTo {
		return nil, fmt.Errorf("date_from cannot be after date_to")
	}

	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// filter by date range
	var filtered []models.Expense
	for _, e := range expenses {
		if e.ExpenseDate >= dateFrom && e.ExpenseDate <= dateTo {
			filtered = append(filtered, e)
		}
	}

	// group by category using a map
	totals := make(map[string]float64)
	counts := make(map[string]int)
	var totalAmount float64

	for _, e := range filtered {
		totals[e.Category] += e.Amount
		counts[e.Category]++
		totalAmount += e.Amount
	}

	// build category summary slice
	var byCategory []CategorySummary
	for category, total := range totals {
		byCategory = append(byCategory, CategorySummary{
			Category: category,
			Total:    total,
			Count:    counts[category],
		})
	}

	// sort by total descending so highest spend appears first
	sort.Slice(byCategory, func(i, j int) bool {
		return byCategory[i].Total > byCategory[j].Total
	})

	return &SummaryResult{
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		TotalAmount: totalAmount,
		TotalCount:  len(filtered),
		ByCategory:  byCategory,
	}, nil 
}
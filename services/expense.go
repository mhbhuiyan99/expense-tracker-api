package services

import(
	"expense-tracker-api/models"
	"fmt"
)

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

func GetExpenses(userID int) ([]models.Expense, error) {
	return models.GetExpensesByUserID(userID)
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
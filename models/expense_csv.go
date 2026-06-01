package models

import(
	"encoding/csv"
	beego "github.com/beego/beego/v2/server/web"
	"os"
	"strconv"
	"time"
	"github.com/beego/beego/v2/core/logs"
	"fmt"
)

func getExpenseFilePath() string {
	dir := beego.AppConfig.DefaultString("datadir", "data")
	file := beego.AppConfig.DefaultString("expensefile", "expenses.csv")
	return dir + "/" + file
}

// getAllExpenses reads all expenses from the CSV file and returns them as a slice of Expense
func getAllExpenses()([]Expense, error){
	path := getExpenseFilePath()

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return []Expense{}, nil // No file means no expenses
		}
		return nil, err
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) <= 1 {
		return []Expense{}, nil 
	}

	var expenses []Expense
	for _, row := range rows[1:] {
		expense, err := rowToExpense(row)
		if err != nil {
			logs.Warn("Skipping invalid expense row: %v", err)
			continue
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

// GetExpensesByUserID returns all expenses belonging to a specific user.
func GetExpensesByUserID(userID int) ([]Expense, error) {
	all, err := getAllExpenses()
	if err != nil {
		return nil, err
	}
	var expenses []Expense
	for _, e := range all {
		if e.UserID == userID {
			expenses = append(expenses, e)
		}
	}
	return expenses, nil
}

// GetExpenseByID retrieves a single expense by its ID from the CSV file
// Return nil if not found or does not belong to the user
func GetExpenseByID(id int, userID int) (*Expense, error) {
	all, err := getAllExpenses()
	if err != nil {
		return nil, err
	}

	for _, e := range all {
		if e.ID == id && e.UserID == userID {
			return &e, nil
		}
	}
	return nil, nil
}

// CreateExpense adds a new expense to the CSV file
func CreateExpense(expense *Expense) error {
	path := getExpenseFilePath()

	dir := beego.AppConfig.DefaultString("datadir", "data")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		if err := writer.Write([]string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}); err != nil {
			return err
		}
	}

	nextID, err := GetNextExpenseID()
	if err != nil {
		return err
	}
	expense.ID = nextID
	expense.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := writer.Write(expenseToRow(expense)); err != nil {
		return err
	}

	return writer.Error()
}

// UpdateExpense updates an existing expense in the CSV file
func UpdateExpense(updated *Expense) error {
	all, err := getAllExpenses()
	if err != nil {
		return err
	}

	found := false
	for i, e := range all {
		if e.ID == updated.ID && e.UserID == updated.UserID {
			all[i] = *updated
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("expense not found")
	}
	return writeAllExpenses(all)
}

// DeleteExpense rewrites the CSV without the deleted expense.
func DeleteExpense(id int, userID int) error {
	all, err := getAllExpenses()
	if err != nil {
		return err
	}

	filtered := make([]Expense, 0, len(all))
	for _, e := range all {
		if !(e.ID == id && e.UserID == userID) {
			filtered = append(filtered, e)
		}
	}
	return writeAllExpenses(filtered)
}

// GetNextExpenseID returns the next available expense ID.
func GetNextExpenseID() (int, error) {
	all, err := getAllExpenses()
	if err != nil {
		return 0, err
	}
	maxID := 0
	for _, e := range all {
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	return maxID + 1, nil
}

// writeAllExpenses overwrites the CSV file with the provided expenses.
func writeAllExpenses(expenses []Expense) error {
	path := getExpenseFilePath()

	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	if err := writer.Write([]string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}); err != nil {
		return err
	}

	for _, e := range expenses {
		if err := writer.Write(expenseToRow(&e)); err != nil {
			return err
		}
	}
	return writer.Error()
}

// rowToExpense converts a CSV row to an Expense struct
func rowToExpense(row []string) (Expense, error) {
	if len(row) < 8 {
		return Expense{}, fmt.Errorf("invalid expense row: expected 8 fields, got %d", len(row))
	}
	id, err := strconv.Atoi(row[0])
	if err != nil {
		return Expense{}, fmt.Errorf("invalid expense id %q: %w", row[0], err)
	}
	userID, err := strconv.Atoi(row[1])
	if err != nil {
		return Expense{}, fmt.Errorf("invalid user_id %q: %w", row[1], err)
	}
	amount, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return Expense{}, fmt.Errorf("invalid amount %q: %w", row[3], err)
	}
	return Expense{
		ID:          id,
		UserID:      userID,
		Title:       row[2],
		Amount:      amount,
		Category:    row[4],
		Note:        row[5],
		ExpenseDate: row[6],
		CreatedAt:   row[7],
	}, nil
}

// expenseToRow converts an Expense struct to a CSV row
func expenseToRow(e *Expense) []string {
	return []string{
		strconv.Itoa(e.ID),
		strconv.Itoa(e.UserID),
		e.Title,
		strconv.FormatFloat(e.Amount, 'f', 2, 64),
		e.Category,
		e.Note,
		e.ExpenseDate,
		e.CreatedAt,
	}
}
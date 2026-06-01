package controllers

import(
	"expense-tracker-api/models"
	"expense-tracker-api/services"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
)

type ExpenseController struct {
	BaseController
}

// getUserID reads and converts the X-User-ID header to an integer
func (c *ExpenseController) getUserID() (int, bool) {
	userIDstr := c.Ctx.Input.Header("X-User-ID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.RespondError(401, "Unauthorized")
		return 0, false
	}
	return userID, true
}

// List handles GET /api/v1/expenses
// @Title List Expenses
// @Description Get all expenses for the authenticated user
// @Param X-User-ID header int true "User ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse 
// @router /api/v1/expenses [get]
func (c *ExpenseController) List() {
	userID, ok := c.getUserID()
	if !ok {
		return
	}

	expenses, err := services.GetExpenses(userID)
	if err != nil {
		logs.Error("[List] %v", err)
		c.RespondError(500, "Failed to retrieve expenses")
		return
	}
	c.RespondSuccess(expenses)
}

// Create handles POST /api/v1/expenses
// @Title Create Expense
// @Description Create a new expense for the authenticated user
// @Param X-User-ID header int true "User ID"
// @Param body body models.CreateExpenseRequest true "Expense data"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @router /api/v1/expenses [post]
func (c *ExpenseController) Create() {
    userID, ok := c.getUserID()
    if !ok {
        return
    }

    var req models.CreateExpenseRequest
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
        logs.Error("[Create] Unmarshal failed: %v", err)
        c.RespondError(400, "Invalid JSON")
        return
    }

    // no ValidateInput here — service handles validation
    expense, err := services.CreateExpense(userID, req)
    if err != nil {
        logs.Warn("[Create] %v", err)
        c.RespondError(400, err.Error())
        return
    }
    c.RespondCreated(expense)
}

// GetOne handles GET /api/v1/expenses/:id
// @Title Get Expense
// @Description Get a single expense by ID
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Expense ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @router /api/v1/expenses/:id [get]
func (c *ExpenseController) GetOne() {
	userID, ok := c.getUserID()
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.RespondError(400, "Invalid expense ID")
		return
	}
	expense, err := services.GetExpense(id, userID)
	if err != nil {
		c.RespondError(404, "Expense not found")
		return
	}
	c.RespondSuccess(expense)
}

// Update handles PUT /api/v1/expenses/:id
// @Title Update Expense
// @Description Update an existing expense
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Expense ID"
// @Param body body models.UpdateExpenseRequest true "Updated fields"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @router /api/v1/expenses/:id [put]
func (c *ExpenseController) Update() {
	userID, ok := c.getUserID()
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.RespondError(400, "Invalid expense ID")
		return
	}

	var req models.UpdateExpenseRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.RespondError(400, "Invalid JSON")
		return
	}

	expense, err := services.UpdateExpense(id, userID, req)
	if err != nil {
		if err.Error() == "expense not found" {
			c.RespondError(404, "Expense not found")
			return
		}
		c.RespondError(400, err.Error())
		return
	}
	c.RespondSuccess(expense)
}

// Remove handles DELETE /api/v1/expenses/:id
// @Title Delete Expense
// @Description Delete an expense by ID
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Expense ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 404 {object} models.ErrorResponse
// @router /api/v1/expenses/:id [delete]
func (c *ExpenseController) Remove() {
	userID, ok := c.getUserID()
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.RespondError(400, "Invalid expense ID")
		return
	}

	if err := services.DeleteExpense(id, userID); err != nil {
		c.RespondError(404, "Expense not found")
		return
	}
	c.RespondSuccess(nil)
}

// Summary handles GET /api/v1/expenses/summary
// @Title Expense Summary
// @Description Get spending summary for a date range
// @Param X-User-ID header string true "User ID"
// @Param date_from query string true "Start date (YYYY-MM-DD)"
// @Param date_to query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @router /api/v1/expenses/summary [get]
func (c *ExpenseController) Summary() {
	c.RespondSuccess(nil) // implement later
}
package controllers

type ExpenseController struct {
	BaseController
}

func (c *ExpenseController) List()    { c.RespondSuccess(nil) }
func (c *ExpenseController) Create()  { c.RespondSuccess(nil) }
func (c *ExpenseController) Summary() { c.RespondSuccess(nil) }
func (c *ExpenseController) GetOne()  { c.RespondSuccess(nil) }
func (c *ExpenseController) Update()  { c.RespondSuccess(nil) }
func (c *ExpenseController) Remove()  { c.RespondSuccess(nil) }
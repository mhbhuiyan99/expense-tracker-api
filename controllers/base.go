package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"expense-tracker-api/models"
)


// BaseController: embedded by every controller in the application
// It provides shared response helpers so all endpoints return a consistent shape
type BaseController struct {
	beego.Controller
}

// RespondSuccess: sends a 200 JSON response with data
func (c *BaseController) RespondSuccess (data interface{}) {
	c.Data["json"] = models.SuccessResponse {
		Status: "success",
		Data: data,
	}
	c.ServeJSON()
}

// RespondCreated: sends a 201 JSON response with data
func (c *BaseController) RespondCreated (data interface{}) {
	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = models.SuccessResponse {
		Status: "success",
		Data: data,
	}
	c.ServeJSON()
}

// RespondError: sends an error JSON response with the given status code
func (c *BaseController) RespondError (statusCode int, message string) {
	c.Ctx.Output.SetStatus(statusCode)
	c.Data["json"] = models.ErrorResponse {
		Status: "error",
		Message: message,
	}
	c.ServeJSON()
}
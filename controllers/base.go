package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"expense-tracker-api/models"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/core/logs"
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

// ValidateInput runs Beego validation on a struct and responds with 400 if invalid.
// Returns true if validation passed, false if it failed (response already sent).
func (c *BaseController) ValidateInput(obj interface{}) bool {
    valid := validation.Validation{}
    ok, err := valid.Valid(obj)
    if err != nil {
        logs.Error("[ValidateInput] internal error: %v", err)
        c.RespondError(500, "Validation error")
        return false
    }
    if !ok {
        msg := "Validation failed"
        if len(valid.Errors) > 0 {
            msg = valid.Errors[0].Key + ": " + valid.Errors[0].Message
        }
        c.RespondError(400, msg)
        return false
    }
    return true
}
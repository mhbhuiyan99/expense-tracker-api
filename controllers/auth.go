package controllers

import (
	"encoding/json"
	"expense-tracker-api/models"
	"expense-tracker-api/services"

	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/core/logs"
)

type AuthController struct {
	BaseController
}

// Register handles POST /api/v1/auth/register
// @Title Register User
// @Description Create a new user account
// @Param body body models.RegisterRequest true "Registration details"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @router /api/v1/auth/register [post]
func (c *AuthController) Register() {
	var req models.RegisterRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.RespondError(400, "Invalid JSON")
		return
	}

	// validate using Beego validation
	valid := validation.Validation{}
	ok, _ := valid.Valid(&req)
	if !ok {
		c.RespondError(400, valid.Errors[0].Key+" "+valid.Errors[0].Message)
		return
	}

	if err := services.RegisterUser(req); err != nil {
		if err.Error() == "email already exists" {
			c.RespondError(409, "Email already exists")
			return
		}

		logs.Error("[Register] %v", err)
		c.RespondError(500, "Registration failed")
		return
	}
	c.RespondCreated(nil)
}


// Login handles POST /api/v1/auth/login
// @Title Login
// @Description Authenticate with email and password
// @Param body body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @router /api/v1/auth/login [post]
func (c *AuthController) Login() {
	var req models.LoginRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.RespondError(400, "Invalid JSON")
		return
	}

	valid := validation.Validation{}
	ok, _ := valid.Valid(&req)
	if !ok {
		c.RespondError(400, valid.Errors[0].Key+" "+valid.Errors[0].Message)
		return
	}

	data, err := services.LoginUser(req)
	if err != nil {
		c.RespondError(401, "Invalid email or password")
		return
	}
	c.RespondSuccess(data)
}
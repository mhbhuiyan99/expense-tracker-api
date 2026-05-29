package controllers

type AuthController struct {
	BaseController
}

func (c *AuthController) Register() {
	c.RespondSuccess(nil)
}

func (c *AuthController) Login() {
	c.RespondSuccess(nil)
}
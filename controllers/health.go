package controllers

// HealthController handles the health check endpoint
type HealthController struct {
	BaseController
}

// Get handlers GET /api/v1/health
// @Title Health Check
// @Description Returns server status
// @Success 200 {object} models.SuccessResponse
// @router /api/v1/health [get]

func (c *HealthController) Get() {
	c.RespondSuccess(map[string]string{"status":"running"})
}
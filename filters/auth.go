package filters

import (
	"github.com/beego/beego/v2/server/web/context"
	"expense-tracker-api/models"
	"strconv"
	"encoding/json"
)

// AuthFilter checks that X-User-ID header is present and valid before allowing access to protected routes
// Registered in router.go and run before every /api/v1/expenses route
func AuthFilter(ctx *context.Context) {
	userIDstr := ctx.Input.Header("X-User-ID")
	if userIDstr == "" {
		respondUnauthorized(ctx, "Unauthorized")
		return
	}

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		respondUnauthorized(ctx, "Invalid user ID")
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil || user == nil {
		respondUnauthorized(ctx, "Unauthorized")
		return
	}
}

func respondUnauthorized(ctx *context.Context, message string) {
	ctx.Output.SetStatus(401)
	body, _ := json.Marshal(models.ErrorResponse{
		Status: "error",
		Message: message,
	})
	ctx.Output.Body(body)
}
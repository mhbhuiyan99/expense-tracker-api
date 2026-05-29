// @APIVersion 1.0.0
// @Title Expense Tracker API
// @Description Personal expense tracker REST API
// @Contract mhbhuiyan10023@gmail.com

package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"expense-tracker-api/controllers"
)

func init() {
	ns := beego.NewNamespace("/api/v1", 

		// Health check
		beego.NSRouter("/health", &controllers.HealthController{}),

		// Auth (public)
		beego.NSNamespace("/auth",
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),

		// Expenses (protected - requires X-User-ID header)
		beego.NSNamespace("expenses",
			beego.NSRouter("", &controllers.ExpenseController{}, "get:List;post:Create"),
			beego.NSRouter("/summary", &controllers.ExpenseController{}, "get:Summary"),
			beego.NSRouter("/:id", &controllers.ExpenseController{}, "get:GetOne;put:Update;delete:Remove"),
		),
	)
	beego.AddNamespace(ns)
}
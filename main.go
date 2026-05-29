package main

import (
	_ "expense-tracker-api/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {

	// Serve Swagger UI in dev mode
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

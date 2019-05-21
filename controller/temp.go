package controller

import (
	"github.com/kataras/iris"
	"receipts/auth"
)

func IndexTemp(ctx iris.Context) {

	auth.CheckToken(ctx)

	username := ctx.GetCookie("username")

	if username == "admin" {
		ctx.RemoveCookie("token")
		ctx.RemoveCookie("username")
		ctx.Redirect("/", 302)
	}

	tempQueryResultData, _ := tempQueryResultMap[username]

	ctx.ViewData("TempResult", tempQueryResultData)

	if err := ctx.View("temp.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}
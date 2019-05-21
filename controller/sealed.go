package controller

import (
	"github.com/kataras/iris"
	"receipts/auth"
)

func IndexSealed(ctx iris.Context) {

	auth.CheckToken(ctx)

	username := ctx.GetCookie("username")
	nickname := ctx.GetCookie("nickname")

	if username == "admin" {
		ctx.RemoveCookie("token")
		ctx.RemoveCookie("username")
		ctx.Redirect("/", 302)
	}

	ctx.ViewData("Username", username)
	ctx.ViewData("NickName", nickname)

	if err := ctx.View("sealed.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}

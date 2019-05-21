package controller

import (
	"github.com/kataras/iris"
	"receipts/auth"
	"receipts/dbop"
	"time"
)

func IndexLogin(ctx iris.Context) {

	auth.RemoveToken(ctx.GetCookie("token"))

	ctx.RemoveCookie("token")
	ctx.RemoveCookie("username")
	ctx.RemoveCookie("nickname")

	if err := ctx.View("login.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}

func Login(ctx iris.Context) {

	if userItem, ok := getUserJSON(ctx, "获取登录信息"); ok {

		checked, nickname, err := dbop.CheckUser(userItem)

		if err != nil {
			_, _ = ctx.JSON(iris.Map{
				"status": "failed",
				"message": "检验登录信息失败：" + err.Error(),
			})
			return
		}

		if checked {
			ctx.SetCookieKV("token", auth.NewToken(userItem.Username), iris.CookieExpires(100 * 365 * 24 * time.Hour))
			ctx.SetCookieKV("username", userItem.Username, iris.CookieExpires(100 * 365 * 24 * time.Hour))
			ctx.SetCookieKV("nickname", nickname, iris.CookieExpires(100 * 365 * 24 * time.Hour))

			dbop.WriteLog("user", "登录", userItem.Username)

			_, _ = ctx.JSON(iris.Map{
				"status": "ok",
				"message": "登录成功",
			})
		} else {
			_, _ = ctx.JSON(iris.Map{
				"status": "failed",
				"message": "用户名或密码错误",
			})
		}

	}
}

func Logout(ctx iris.Context) {
	ctx.Redirect("/", 302)
}

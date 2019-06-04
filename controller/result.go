package controller

import (
	"github.com/kataras/iris"
	"receipts/dbop"
	"receipts/model"
	"strconv"
	"strings"
)

func IndexResult(ctx iris.Context) {

	//auth.CheckToken(ctx)

	username := ctx.GetCookie("username")
	nickname := ctx.GetCookie("nickname")

	if username == "admin" {
		ctx.RemoveCookie("token")
		ctx.RemoveCookie("username")
		ctx.Redirect("/", 302)
	}

	ctx.ViewData("Username", username)
	ctx.ViewData("NickName", nickname)

	if err := ctx.View("result.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}

func ProcessResultData(ctx iris.Context) {

	if resultPara, ok := getResultParaJSON(ctx); ok {
		if resultPara.Operation == "getdata" {
			_, _ = ctx.JSON(dbop.GetResultData(resultPara))
		} else if judgeUpdateOperation(resultPara.Operation) {
			_, _ = ctx.JSON(dbop.UpdateResultData(resultPara))
		} else if resultPara.Operation == "removedata" {
			_, _ = ctx.JSON(dbop.RemoveResult(resultPara))
		}
	}
}

func judgeUpdateOperation(operation string) bool {
	if operation == "sealed" || operation == "unsealed" ||
		operation == "ensure" || operation == "unensure" {
		return true
	}
	return false
}

// 解析ResultPara的JSON对象
func getResultParaJSON(ctx iris.Context) (*model.ResultPara, bool) {


	filter := ctx.FormValue("filter")
	filter = strings.Replace(filter, "@@@", "%", -1)

	rows, _ := strconv.Atoi(ctx.FormValue("rows"))
	page, _ := strconv.Atoi(ctx.FormValue("page"))

	resultPara := &model.ResultPara{
		UserName: ctx.FormValue("username"),
		Filter: filter,
		ResultId: ctx.FormValue("resultid"),
		Operation: ctx.FormValue("operation"),
		Rows: rows,
		Page: page,
	}

	return resultPara, true
}
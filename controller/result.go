package controller

import (
	"github.com/kataras/iris"
	"io/ioutil"
	"receipts/auth"
	"receipts/dbop"
	"receipts/model"
	"strconv"
	"strings"
)

func IndexResult(ctx iris.Context) {

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

	if err := ctx.View("result.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}

func ProcessResultData(ctx iris.Context) {

	if resultPara, ok := getResultParaJSON(ctx, "获取归档请求参数"); ok {
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
func getResultParaJSON(ctx iris.Context, info string) (*model.ResultPara, bool) {

	resultPara := &model.ResultPara{}

	b, err := ioutil.ReadAll(ctx.Request().Body)

	if err != nil{
		_, _ = ctx.JSON(iris.Map{
			"status":  "failed",
			"message": info + "失败" + err.Error(),
		})

		return resultPara, false
	}

	paraList := strings.Split(string(b), "&")

	resultPara.UserName = strings.Split(paraList[0], "=")[1]
	resultPara.Filter = paraList[1][7:]
	resultPara.ResultId = strings.Split(paraList[2], "=")[1]
	resultPara.Operation = strings.Split(paraList[3], "=")[1]
	resultPara.Rows, err = strconv.Atoi(strings.Split(paraList[4], "=")[1])
	resultPara.Rows, err = strconv.Atoi(strings.Split(paraList[5], "=")[1])

	if err != nil{
		_, _ = ctx.JSON(iris.Map{
			"status":  "failed",
			"message": info + "失败" + err.Error(),
		})
		return resultPara, false
	}

	return resultPara, true
}

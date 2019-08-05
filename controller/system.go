package controller

import (
	"github.com/kataras/iris"
	"receipts/auth"
	"receipts/dbop"
	"receipts/model"
)

func IndexSystem(ctx iris.Context) {

	auth.CheckToken(ctx)
	auth.CheckAdmin(ctx)

	systemInfo, err := dbop.GetSystemInfo()

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
		return
	}

	//ctx.ViewData("Password", dbop.SystemInfo.Password) // ???
	ctx.ViewData("Password", systemInfo.Password)

	unusedUsage, err := dbop.UCGetUnusedUsage()

	if err != nil {
		ctx.ViewData("UnusedUsage", err.Error())
	} else {
		ctx.ViewData("UnusedUsage", unusedUsage)
	}

	//ctx.ViewData("ApiCode", dbop.SystemInfo.ApiCode) // ???
	ctx.ViewData("ApiCode", systemInfo.ApiCode)

	if err := ctx.View("system.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}

// 更新密码
func UpdateSystemPassword(ctx iris.Context) {

	auth.CheckToken(ctx)
	auth.CheckAdmin(ctx)

	if systemItem, ok := getSystemJSON(ctx, "更新系统密码"); ok {
		_, _ = ctx.JSON(dbop.UpdateSystemPassword(systemItem))
	}
}

// 更新闲置余额
func UpdateSystemUnusedUsage(ctx iris.Context) {

	auth.CheckToken(ctx)
	auth.CheckAdmin(ctx)

	if systemItem, ok := getSystemJSON(ctx, "更新系统闲置余额"); ok {

		if systemItem.UnusedUsage < 0 {
			_, _ = ctx.JSON(iris.Map{
				"status": "failed",
				"message": "更新系统闲置余额不能少于0",
			})
		}

		err := dbop.UCUpdateUnusedUsage(systemItem)

		if err != nil {
			_, _ = ctx.JSON(iris.Map{
				"status": "failed",
				"message": "更新系统闲置余额失败：" + err.Error(),
			})
		} else {
			_, _ = ctx.JSON(iris.Map{
				"status": "ok",
				"message": "更新系统闲置余额成功",
			})
		}
	}
}

// 更新ApiCode
func UpdateSystemApiCode(ctx iris.Context) {

	auth.CheckToken(ctx)
	auth.CheckAdmin(ctx)

	if systemItem, ok := getSystemJSON(ctx, "更新apicode"); ok {
		_, _ = ctx.JSON(dbop.UpdateSystemApiCode(systemItem))
	}
}

// 解析User的JSON对象
func getSystemJSON(ctx iris.Context, info string) (*model.System, bool) {

	systemItem := &model.System{}

	err := ctx.ReadJSON(systemItem)

	if err != nil{

		_, _ = ctx.JSON(iris.Map{
			"status":  "failed",
			"message": info + "失败" + err.Error(),
		})

		return systemItem, false

	}

	return systemItem, true
}

// 导出全部数据
func OutputAllData(ctx iris.Context) {

	auth.CheckToken(ctx)
	auth.CheckAdmin(ctx)

	err := dbop.BuildAllDataExcel()

	if err != nil {
		_, _ = ctx.JSON(iris.Map{
			"status": "failed",
			"message": "导出全量数据失败：" + err.Error(),
		})
		return
	}

	file := "./全量数据.xls"
	err = ctx.SendFile(file, "全量数据.xls")

	if err != nil {
		_, _ = ctx.JSON(iris.Map{
			"status": "failed",
			"message": "下载全量数据失败：" + err.Error(),
		})
		return
	}
}
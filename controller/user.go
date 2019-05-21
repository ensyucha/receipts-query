package controller

import (
	"github.com/kataras/iris"
	"receipts/auth"
	"receipts/dbop"
	"receipts/model"
)

// 新增用户
func AddUser(ctx iris.Context) {

	auth.CheckToken(ctx)

	if userItem, ok := getUserJSON(ctx, "新增用户"); ok {
		_, _ = ctx.JSON(dbop.AddUser(userItem))
	}
}

// 删除用户
func RemoveUser(ctx iris.Context) {

	auth.CheckToken(ctx)

	if userItem, ok := getUserJSON(ctx, "删除用户"); ok {
		_, _ = ctx.JSON(dbop.RemoveUser(userItem))
	}
}

// 更新用户
func UpdateUser(ctx iris.Context) {

	auth.CheckToken(ctx)

	if userItem, ok := getUserJSON(ctx, "更新用户"); ok {
		_, _ = ctx.JSON(dbop.UpdateUser(userItem))
	}
}

// 获取用户列表
func ListUser(ctx iris.Context) {

	auth.CheckToken(ctx)

	_, _ = ctx.JSON(dbop.ListUser())
}

// 解析User的JSON对象
func getUserJSON(ctx iris.Context, info string) (*model.User, bool) {

	userItem := &model.User{}

	err := ctx.ReadJSON(userItem)

	if err != nil{

		_, _ = ctx.JSON(iris.Map{
			"status":  "failed",
			"message": info + "失败" + err.Error(),
		})

		return userItem, false

	}

	return userItem, true
}
package dbop

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"log"
	"receipts/auth"
	"receipts/model"
)

// 获取管理信息
func GetSystemInfo() (*model.System, error) {

	system := &model.System{}

	var result, err = db.Query("SELECT * FROM systems")

	if err != nil {
		return system, err
	}

	defer result.Close()

	if result.Next() {

		err = result.Scan(&system.Password, &system.UnusedUsage, &system.ApiCode)

		if err != nil {
			return system, err
		}
	}

	system.UnusedUsage = -1 // 额度相关获取全部需要经过额度中心

	return system, err
}

func UpdateSystemPassword(system *model.System) context.Map {

	stmt, err := db.Prepare("UPDATE systems SET password = ?;")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新password失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(system.Password)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新password失败" + err.Error(),
		}
	}

	log.Println("管理员密码：", system.Password)

	auth.RemoveTokenByUsername("admin")

	WriteLog("system", "更新管理密码为：" + system.Password, "manager")

	return iris.Map{
		"status": "ok",
		"message": "更新password成功",
	}
}

func UpdateSystemApiCode(system *model.System) context.Map {

	stmt, err := db.Prepare("UPDATE systems SET apicode = ?;")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新ApiCode失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(system.ApiCode)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新ApiCode失败：" + err.Error(),
		}
	}

	WriteLog("system", "更新ApiCode为：" + system.ApiCode, "manager")

	return iris.Map{
		"status": "ok",
		"message": "更新ApiCode成功",
	}
}
package dbop

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"log"
	"receipts/auth"
	"receipts/model"
)

// 新增用户
func AddUser(user *model.User) context.Map {

	if user.Username == "admin" {
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败: 不允许用户名为 admin",
		}
	}

	if user.Usage < 0 {
		return iris.Map{
			"status": "failed",
			"message": "预分配额度不能少于0",
		}
	}

	newUnusedUsage, ok, err := UCEnough(user.Usage)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败: " + err.Error(),
		}
	}

	// 额度判断能否分配
	if !ok { // 新增用户预分配额度大于闲置额度，无法创建用户
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败: 预分配额度大于闲置额度",
		}
	}

	stmt, err := db.Prepare("INSERT INTO users(username, nickname, password, usages) VALUES (?,?,?,?);")

	if err != nil {

		return iris.Map{
			"status":  "failed",
			"message": "新增用户失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.NickName, user.Password, user.Usage)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "新增用户失败" + err.Error(),
		}
	}

	err = UCUpdateUnusedUsage(&model.System{UnusedUsage:newUnusedUsage})

	if err != nil {
		return iris.Map{
			"status": "ok",
			"message": "新增用户失败：" + err.Error(),
		}
	}

	err = createResultTable(user.Username)

	if err != nil {
		err2 := removeResultTable(user.Username)

		if err2 != nil {
			return iris.Map{
				"status":  "failed",
				"message": "新增用户失败: 无法创建结果表：" + err.Error() + err2.Error(),
			}
		} else {
			return iris.Map{
				"status":  "failed",
				"message": "新增用户失败: 无法创建结果表：" + err.Error(),
			}
		}
	}

	WriteLog("system", "新增用户：" + user.Username, "manager")

	return iris.Map{
		"status": "ok",
		"message": "新增用户成功",
	}
}

// 删除用户
func RemoveUser(user *model.User) context.Map {

	remainUsage, err := UCGetUserUsage(user)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	nowUnusedUsage, err := UCGetUnusedUsage()

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	err = UCUpdateUnusedUsage(&model.System{UnusedUsage:remainUsage + nowUnusedUsage})

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	stmt, err := db.Prepare("DELETE FROM users WHERE username=?;")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Username)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除用户失败" + err.Error(),
		}
	}

	auth.RemoveTokenByUsername(user.Username) // 删除对应的token

	err = removeResultTable(user.Username)

	if err != nil {
		return iris.Map{
			"status": "ok",
			"message": "删除用户数据表失败：" + err.Error(),
		}
	}

	WriteLog("system", "删除用户：" + user.Username, "manager")

	return iris.Map{
		"status": "ok",
		"message": "删除用户成功",
	}
}

// 更新用户
func UpdateUser(user *model.User) context.Map {

	if user.Usage < 0 {
		return iris.Map{
			"status": "failed",
			"message": "分配额度不能少于0",
		}
	}

	nowUserUsage, err := UCGetUserUsage(user)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "获取用户原始额度失败:" + err.Error(),
		}
	}

	diffUsage := user.Usage - nowUserUsage // 获取额度差值

	newUnusedUsage, ok, err := UCEnough(diffUsage)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "计算额度差值失败:" + err.Error(),
		}
	}

	if !ok {
		return iris.Map{
			"status": "failed",
			"message": "更新用户失败: 预分配额度大于未分配额度",
		}
	}

	stmt, err := db.Prepare("UPDATE users SET nickname=?, password=?, usages=? WHERE username=?;")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新用户失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	log.Println(user.NickName, user.Password, user.Usage, user.Username)

	_, err = stmt.Exec(user.NickName, user.Password, user.Usage, user.Username)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新用户失败" + err.Error(),
		}
	}

	err = UCUpdateUnusedUsage(&model.System{UnusedUsage:newUnusedUsage})

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "更新未分配额度失败:" + err.Error(),
		}
	}

	auth.RemoveTokenByUsername(user.Username) // 删除对应的token

	WriteLog("system", "更新用户：" + user.Username, "manager")

	return iris.Map{
		"status": "ok",
		"message": "更新用户成功",
	}
}

// 获取用户列表
func ListUser() context.Map {

	var group []model.User

	userResult, err := db.Query("SELECT * FROM users")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "获取用户列表失败: " + err.Error(),
		}
	}

	defer userResult.Close()

	for userResult.Next() {

		var userItem model.User

		err = userResult.Scan(&userItem.UserId, &userItem.Username, &userItem.NickName, &userItem.Password, &userItem.Usage)

		if err != nil {
			return iris.Map{
				"status": "failed",
				"message": "获取用户列表失败: " + err.Error(),
			}
		}

		group = append(group, userItem)
	}

	return iris.Map{
		"status": "ok",
		"message": group,
	}
}

func CheckUser(user *model.User) (bool, string, error) {

	var systemItem model.System
	var userItem model.User
	var querySQL string

	if user.Username == "admin" {
		querySQL = "SELECT password FROM systems;"
	} else {
		querySQL = "SELECT nickname, password FROM users WHERE username='" + user.Username + "';"
	}

	queryResult, err := db.Query(querySQL)

	if err != nil {
		return false, "", err
	}

	defer queryResult.Close()

	if queryResult.Next() {

		if user.Username == "admin" {
			err = queryResult.Scan(&systemItem.Password)
		} else {
			err = queryResult.Scan(&userItem.NickName, &userItem.Password)
		}

		if err != nil {
			return false, "", err
		}
	}

	if (user.Username == "admin" && user.Password != systemItem.Password) ||
		(user.Username != "admin" && user.Password != userItem.Password) {
		return false, "", nil
	}

	return true, userItem.NickName, nil
}

func createResultTable(username string) error {

	sql := `CREATE TABLE IF NOT EXISTS result_` + username + ` (
	resultid INTEGER PRIMARY KEY AUTO_INCREMENT,
	ensured TEXT,
	sealed TEXT,
	respCode TEXT,
	respMsg TEXT,
	qd TEXT,
	fpdm TEXT,
	fphm TEXT,
	kprq TEXT,
	yzmSj TEXT,
	fpzt TEXT,
	fxqy TEXT,
	fplx TEXT,
	jqbm TEXT,
	jym TEXT,
	gfName TEXT,
	gfNsrsbh TEXT,
	gfAddressTel TEXT,
	gfBankZh TEXT,
	jshjL TEXT,
	sfName TEXT,
	sfNsrsbh TEXT,
	sfAddressTel TEXT,
	sfBankZh TEXT,
	bz TEXT,
	jshjU TEXT,
	zpListString TEXT
);`

	_, err := db.Exec(sql)

	return err
}

func removeResultTable(username string) error {

	sql := `DROP TABLE result_` + username + `;`

	_, err := db.Exec(sql)

	return err
}
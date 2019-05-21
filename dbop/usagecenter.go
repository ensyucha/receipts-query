package dbop

import (
	"receipts/model"
	"strconv"
)

// 更新系统闲置额度ok
func UCUpdateUnusedUsage(system *model.System) error {

	stmt, err := db.Prepare("UPDATE systems SET unusedusage=?;")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(system.UnusedUsage)

	if err != nil {
		return err
	}

	str := strconv.Itoa(system.UnusedUsage)

	WriteLog("system", "更新系统闲置额度为：" + str, "manager")

	return nil
}

// 更新用户额度ok
func UCUpdateUserUsage(user *model.User) error {

	stmt, err := db.Prepare("UPDATE users SET usages=? WHERE username=?;")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Usage, user.Username)

	if err != nil {
		return err
	}

	str := strconv.Itoa(user.Usage)

	WriteLog("system", "更新" + user.Username + "额度为：" + str, "manager")

	return nil
}

// 根据用户名查询它的额度ok
func UCGetUserUsage(user *model.User) (int, error) {

	userUsageResult, err := db.Query("SELECT usages FROM users WHERE username='"+ user.Username +"'")

	var userItem model.User

	if err != nil {
		return -1, err
	}

	defer userUsageResult.Close()

	if userUsageResult.Next() {

		err = userUsageResult.Scan(&userItem.Usage)

		if err != nil {
			return -1, err
		}
	}

	return userItem.Usage, nil
}

// 获取系统剩余额度ok
func UCGetUnusedUsage() (int, error) {

	unusedUsageResult, err := db.Query("SELECT unusedusage FROM systems;")

	var systemItem model.System

	if err != nil {
		return -1, err
	}

	defer unusedUsageResult.Close()

	if unusedUsageResult.Next() {

		err = unusedUsageResult.Scan(&systemItem.UnusedUsage)

		if err != nil {
			return -1, err
		}
	}

	return systemItem.UnusedUsage, nil
}

// 判断剩余额度是否足够分配ok，且计算出假设完成分配后的未分配额度
func UCEnough(t int) (int, bool, error) {

	nowUnusedUsage, err := UCGetUnusedUsage()

	if err != nil {
		return -1, false, err
	}

	return nowUnusedUsage - t, nowUnusedUsage >= t, nil
}

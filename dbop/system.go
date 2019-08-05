package dbop

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"log"
	"receipts/auth"
	"receipts/model"
	"strconv"
	"strings"
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

func BuildAllDataExcel() error {

	resultItemList, err := getAllData()

	if err != nil {
		return err
	}

	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "所有者")
	f.SetCellValue("Sheet1", "B1", "确认状态")
	f.SetCellValue("Sheet1", "C1", "是否密封")
	f.SetCellValue("Sheet1", "D1", "验证状态")
	f.SetCellValue("Sheet1", "E1", "发票状态")
	f.SetCellValue("Sheet1", "F1", "发票类型")
	f.SetCellValue("Sheet1", "G1", "发票代码")
	f.SetCellValue("Sheet1", "H1", "发票号码")
	f.SetCellValue("Sheet1", "I1", "开票日期")
	f.SetCellValue("Sheet1", "J1", "验证时间")
	f.SetCellValue("Sheet1", "K1", "购方名称")
	f.SetCellValue("Sheet1", "L1", "购方识别号")
	f.SetCellValue("Sheet1", "M1", "购方联系地址")
	f.SetCellValue("Sheet1", "N1", "购方开户行")
	f.SetCellValue("Sheet1", "O1", "销售方名称")
	f.SetCellValue("Sheet1", "P1", "销售方识别号")
	f.SetCellValue("Sheet1", "Q1", "销售方联系地址")
	f.SetCellValue("Sheet1", "R1", "销售方开户行")
	f.SetCellValue("Sheet1", "S1", "价税合计")
	f.SetCellValue("Sheet1", "T1", "校验码")
	f.SetCellValue("Sheet1", "U1", "商品名称")
	f.SetCellValue("Sheet1", "V1", "规格型号")
	f.SetCellValue("Sheet1", "W1", "数量")
	f.SetCellValue("Sheet1", "X1", "单位")
	f.SetCellValue("Sheet1", "Y1", "单价")
	f.SetCellValue("Sheet1", "Z1", "金额")
	f.SetCellValue("Sheet1", "AA1", "税率")
	f.SetCellValue("Sheet1", "AB1", "税额")
	f.SetCellValue("Sheet1", "AC1", "发票总金额")
	f.SetCellValue("Sheet1", "AD1", "发票总税额")
	f.SetCellValue("Sheet1", "AE1", "查询时间")
	f.SetCellValue("Sheet1", "AF1", "价税合计（大写）")
	f.SetCellValue("Sheet1", "AG1", "清单")
	f.SetCellValue("Sheet1", "AH1", "备注")
	f.SetCellValue("Sheet1", "AI1", "风险企业验证")
	f.SetCellValue("Sheet1", "AJ1", "机器编码")

	for index, item := range resultItemList {

		index += 2
		indexStr := strconv.Itoa(index)

		f.SetCellValue("Sheet1", "A" + indexStr, item.Username)
		if strings.Contains(item.Ensured, "未") {
			f.SetCellValue("Sheet1", "B" + indexStr, "未确认")
		} else {
			f.SetCellValue("Sheet1", "B" + indexStr, "已确认")
		}
		if item.Sealed == "0" {
			f.SetCellValue("Sheet1", "C"+indexStr, "否")
		} else {
			f.SetCellValue("Sheet1", "C"+indexStr, "是")
		}
		f.SetCellValue("Sheet1", "D" + indexStr, item.RespMsg)
		f.SetCellValue("Sheet1", "E" + indexStr, item.Fpzt)
		f.SetCellValue("Sheet1", "F" + indexStr, item.Fplx)
		f.SetCellValue("Sheet1", "G" + indexStr, item.Fpdm)
		f.SetCellValue("Sheet1", "H" + indexStr, item.Fphm)
		f.SetCellValue("Sheet1", "I" + indexStr, item.Kprq)
		f.SetCellValue("Sheet1", "J" + indexStr, item.YzmSj)
		f.SetCellValue("Sheet1", "K" + indexStr, item.GfName)
		f.SetCellValue("Sheet1", "L" + indexStr, item.GfNsrsbh)
		f.SetCellValue("Sheet1", "M" + indexStr, item.GfAddressTel)
		f.SetCellValue("Sheet1", "N" + indexStr, item.GfBankZh)
		f.SetCellValue("Sheet1", "O" + indexStr, item.SfName)
		f.SetCellValue("Sheet1", "P" + indexStr, item.SfNsrsbh)
		f.SetCellValue("Sheet1", "Q" + indexStr, item.SfAddressTel)
		f.SetCellValue("Sheet1", "R" + indexStr, item.SfBankZh)
		f.SetCellValue("Sheet1", "S" + indexStr, item.JshjL)
		f.SetCellValue("Sheet1", "T" + indexStr, item.Jym)
		f.SetCellValue("Sheet1", "U" + indexStr, item.MxName)
		f.SetCellValue("Sheet1", "V" + indexStr, item.Ggxh)
		f.SetCellValue("Sheet1", "W" + indexStr, item.Num)
		f.SetCellValue("Sheet1", "X" + indexStr, item.Unit)
		f.SetCellValue("Sheet1", "Y" + indexStr, item.Price)
		f.SetCellValue("Sheet1", "Z" + indexStr, item.Je)
		f.SetCellValue("Sheet1", "AA" + indexStr, item.Sl)
		f.SetCellValue("Sheet1", "AB" + indexStr, item.Se)
		f.SetCellValue("Sheet1", "AC" + indexStr, item.TotalJe)
		f.SetCellValue("Sheet1", "AD" + indexStr, item.TotalSe)
		f.SetCellValue("Sheet1", "AE" + indexStr, item.QueryTime)
		f.SetCellValue("Sheet1", "AF" + indexStr, item.JshjU)
		f.SetCellValue("Sheet1", "AG" + indexStr, item.Qd)
		f.SetCellValue("Sheet1", "AH" + indexStr, item.Bz)
		f.SetCellValue("Sheet1", "AI" + indexStr, item.Fxqy)
		f.SetCellValue("Sheet1", "AJ" + indexStr, item.Jqbm)
	}



	err = f.SaveAs("./全量数据.xls")
	if err != nil {
		return err
	}

	return nil
}

func getAllData() ([]model.ResultItemWithUsername, error){

	var resultItemList []model.ResultItemWithUsername

	users, err := getAllUser()

	if err != nil {
		return nil, err
	}

	sql := buildOutputAllDataSQL(users)

	results, err := db.Query(sql)

	if err != nil {
		return nil, err
	}

	defer results.Close()

	for results.Next() {

		var rdb model.ResultItemWithUsername

		err = results.Scan(&rdb.ResultId, &rdb.Ensured, &rdb.Sealed, &rdb.RespCode, &rdb.RespMsg, &rdb.Qd, &rdb.Fpdm,
			&rdb.Fphm, &rdb.Kprq, &rdb.YzmSj, &rdb.Fpzt, &rdb.Fxqy, &rdb.Fplx, &rdb.Jqbm, &rdb.Jym, &rdb.GfName,
			&rdb.GfNsrsbh, &rdb.GfAddressTel, &rdb.GfBankZh, &rdb.JshjL, &rdb.SfName, &rdb.SfNsrsbh, &rdb.SfAddressTel,
			&rdb.SfBankZh, &rdb.Bz, &rdb.JshjU, &rdb.MxName, &rdb.Ggxh, &rdb.Unit, &rdb.Price, &rdb.Je, &rdb.Sl, &rdb.Se,
			&rdb.TotalJe, &rdb.TotalSe, &rdb.QueryTime, &rdb.Num, &rdb.Username)

		if err != nil {
			return nil, err
		}

		resultItemList = append(resultItemList, rdb)
	}

	return resultItemList, nil
}

func buildOutputAllDataSQL(users []string) string {

	var sqls []string

	for _, item := range users {
		temp := `SELECT *,'` + item + `' as 'username' FROM result_` + item
		sqls = append(sqls, temp)
	}

	return strings.Join(sqls, " UNION ")
}

func getAllUser() ([]string, error) {

	var users []string

	var result, err = db.Query("SELECT username FROM users")

	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {

		var tempUsername string

		err = result.Scan(&tempUsername)

		if err != nil {
			return nil, err
		}

		users = append(users, tempUsername)
	}

	return users, nil
}
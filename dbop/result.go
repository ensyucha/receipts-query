package dbop

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"receipts/model"
	"strings"
)

func GetResultData(resultPara *model.ResultPara) context.Map {

	var resultItemList []model.ResultItem

	var sql = "SELECT * FROM result_" + resultPara.UserName + " WHERE "
	sql += resultPara.Filter + ";"

	var result, err = db.Query(sql)
	
	if err != nil {
		return iris.Map{
			"total": 0,
			"rows": []model.ResultItem{},
		}
	}

	defer result.Close()

	for result.Next() {

		rdb := model.ResultItem{}

		err := result.Scan(&rdb.ResultId, &rdb.Ensured, &rdb.Sealed, &rdb.RespCode, &rdb.RespMsg, &rdb.Qd, &rdb.Fpdm,
			&rdb.Fphm, &rdb.Kprq, &rdb.YzmSj, &rdb.Fpzt, &rdb.Fxqy, &rdb.Fplx, &rdb.Jqbm, &rdb.Jym, &rdb.GfName,
			&rdb.GfNsrsbh, &rdb.GfAddressTel, &rdb.GfBankZh, &rdb.JshjL, &rdb.SfName, &rdb.SfNsrsbh, &rdb.SfAddressTel,
			&rdb.SfBankZh, &rdb.Bz, &rdb.JshjU, &rdb.MxName, &rdb.Ggxh, &rdb.Unit, &rdb.Price, &rdb.Je, &rdb.Sl, &rdb.Se,
			&rdb.TotalJe, &rdb.TotalSe, &rdb.QueryTime, &rdb.Num)

		if err != nil {
			return iris.Map{
				"total": 0,
				"rows": []model.ResultItem{},
			}
		}

		resultItemList = append(resultItemList, rdb)
	}

	if len(resultItemList) == 0 {
		return iris.Map{
			"total": 1,
			"rows": "",
		}
	}

	return iris.Map{
		"total": len(resultItemList),
		"rows": resultItemList,
	}
}

// 删除归档
func RemoveResult(para *model.ResultPara) context.Map {

	idSet := getRequestIdSet(para)

	stmt, err := db.Prepare("DELETE FROM result_" + para.UserName+ " WHERE resultid IN " + idSet + " ;")

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除归档失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "删除归档失败" + err.Error(),
		}
	}

	WriteLog("user", "删除归档：" + idSet, para.UserName)

	return iris.Map{
		"status": "ok",
		"message": "请求成功，正在删除归档数据...",
	}
}

// 更新结果：封存，解封，确认，取消确认
func UpdateResultData(para *model.ResultPara) context.Map {

	sql := ""
	hint := ""
	idSet := getRequestIdSet(para)

	if para.Operation == "sealed" {
		sql = "UPDATE result_" + para.UserName + " SET sealed='1' WHERE resultid IN " + idSet + " ;"
		hint = "请求封存"
	} else if para.Operation == "unsealed" {
		sql = "UPDATE result_" + para.UserName + " SET sealed='0' WHERE resultid IN " + idSet + " ;"
		hint = "请求取消封存"
	} else if para.Operation == "ensure" {
		sql = "UPDATE result_" + para.UserName + ` SET ensured='<span class="my-success">已确认</span>' WHERE resultid IN ` + idSet + " ;"
		hint = "请求确认"
	} else if para.Operation == "unensure" {
		sql = "UPDATE result_" + para.UserName + ` SET ensured='<span class="my-failed">未确认</span>' WHERE resultid IN ` + idSet + " ;"
		hint = "请求取消确认"
	}

	stmt, err := db.Prepare(sql)

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": hint + "失败: " + err.Error(),
		}
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": hint + "失败: " + err.Error(),
		}
	}

	WriteLog("resultAndSealed", hint + " : " + idSet, para.UserName)

	return iris.Map{
		"status": "ok",
		"message": hint + "成功，正在处理数据...",
	}
}

func getRequestIdSet(para *model.ResultPara) string {

	ids := strings.Split(para.ResultId,"-")

	return "(" + strings.Join(ids, ",") + ")"
}

package controller

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"io/ioutil"
	"log"
	"net/http"
	"receipts/auth"
	"receipts/dbop"
	"receipts/model"
	"strconv"
	"sync"
	"time"
)

var waitGroup sync.WaitGroup
var client *http.Client
var tempQueryResultMap = make(map[string][]*model.QueryResult)

func init() {
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
}

func IndexQuery(ctx iris.Context) {

	auth.CheckToken(ctx)

	username := ctx.GetCookie("username")

	if username == "admin" {
		ctx.RemoveCookie("token")
		ctx.RemoveCookie("username")
		ctx.Redirect("/", 302)
	}

	user := &model.User{Username:username}

	usage, err := dbop.UCGetUserUsage(user)

	if err != nil {
		ctx.ViewData("Usage", err.Error())
	} else {
		ctx.ViewData("Usage", usage)
	}

	nickName := ctx.GetCookie("nickname")

	ctx.ViewData("NickName", nickName)


	if err := ctx.View("query.html"); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.Writef(err.Error())
	}
}

func AcceptQuery(ctx iris.Context) {

	auth.CheckToken(ctx)

	username := ctx.GetCookie("username")

	userUsage, err := dbop.UCGetUserUsage(&model.User{Username:username})

	if err != nil {
		_, _ = ctx.JSON(iris.Map{
			"status":  "failed",
			"message": "查询用户额度失败：" + err.Error(),
		})
		return
	}

	if queryArray, ok := getQueryArrayJSON(ctx, "读取查询组信息"); ok {
		_, _ = ctx.JSON(processQueryArray(queryArray, &model.User{Username:username, Usage:userUsage}))
	}
}

func processQueryArray(queryArray *model.QueryArray, user *model.User) context.Map {

	queryNum := len(queryArray.QueryArray) // 需要查询的发票总数

	if queryNum > user.Usage {
		return iris.Map{
			"status": "failed",
			"message": "用户额度不足（为" + strconv.Itoa(user.Usage) + "），请增加额度或减少查询量",
		}
	}

	systemInfo, err := dbop.GetSystemInfo()

	if err != nil {
		return iris.Map{
			"status": "failed",
			"message": "获取系统ApiCode失败：" + err.Error(),
		}
	}

	queryChan := make(chan *model.Query, queryNum) // 将需要查询的发票放入 in Chan
	for _, query := range queryArray.QueryArray {
		queryChan <- query
	}

	queryResultChan := make(chan *model.QueryResult, queryNum) // 将查询结果放入 out Chan

	goroutineNum := queryNum

	if dbop.GetRate() < 1.0 {
		goroutineNum = int(dbop.GetRate()*float64(queryNum)) + 1 // 开启的 goroutine数量，加1保证至少有一个
	}

	qinfo := fmt.Sprintf("查询发票数量：%d，协程数量：%d \n", queryNum, goroutineNum)
	dbop.WriteLog("user", qinfo, user.Username)

	waitGroup.Add(queryNum)

	for i := 0; i < goroutineNum; i++ {
		go doQuery(systemInfo.ApiCode, queryChan, queryResultChan)
	}

	successQuery := 0
	var tempQueryResultList []*model.QueryResult

	for {
		select {
		case queryResult := <- queryResultChan:
			//log.Println(queryResult)

			tempQueryResultList = append(tempQueryResultList, queryResult)

			if judgeQuerySuccess(queryResult.RespCode) {

				err = dbop.AddResult(user.Username, queryResultToResultDBs(queryResult))

				if err != nil {
					return iris.Map{
						"status": "failed",
						"message": "查询结果新增失败，查询中止：" + err.Error(),
					}
				} else {
					successQuery++
				}
			}
			queryNum--
		default:
			if queryNum == 0 && len(queryResultChan) == 0 {

				err = dbop.UCUpdateUserUsage(&model.User{Username:user.Username, Usage:user.Usage - successQuery})

				if  err != nil{
					dbop.WriteLog("error", "【严重错误！】查询已经成功，但额度更新失败 | 用户名：" +
						user.Username + " | 错误描述：" + err.Error(), "error")
					log.Fatal("【严重错误！】查询已经成功，但额度更新失败 | 用户名：" +
						user.Username + " | 错误描述：" + err.Error())
				} else {
					err = dbop.UCAddUserHistoryTotal(user.Username, successQuery)

					if err != nil {
						panic("更新历史使用额度失败：" + err.Error())
					}

					tempQueryResultMap[user.Username] = tempQueryResultList
					return iris.Map{
						"status": "ok",
						"message": "查询完毕",
					}
				}
			}
		}
	}
}

func doQuery(apiCode string, queryChan chan *model.Query, queryResultChan chan *model.QueryResult) {

	for {
		select {
		case query := <- queryChan:
			queryRequest := makeQueryRequest(query, apiCode)

			resp, err := client.Do(queryRequest) // 发送请求

			queryResult := &model.QueryResult{
				Data:model.QueryResultData{
					Fpdm: query.Fpdm,
					Fphm: query.Fphm,
					Kprq: query.Kprq,
					Jym: query.Jym,
					Bz: query.Je, // ！！临时用 备注 表示金额
				},
			}

			if err != nil {
				queryResult.RespCode = "3000"
				queryResult.RespMsg = "API请求超时，请稍后再试"
				dbop.WriteLog("system", "3000：" + err.Error(), "system")
			} else {
				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					queryResult.RespCode = "3001"
					queryResult.RespMsg = "读取结果二进制失败：" + err.Error()
				} else {
					if len(body) != 0 {
						err = json.Unmarshal(body, queryResult)

						if err != nil {
							queryResult.RespCode = "3002"
							queryResult.RespMsg = `<span style="color:#AD1457;">Client请求超时，请稍后再试</span>`
							dbop.WriteLog("system", "3002："+err.Error(), "system")
						}
					} else {
						queryResult.RespCode = "3003"
						queryResult.RespMsg = `<span style="color:#AD1457;">联系管理员检查系统额度</span>`
						dbop.WriteLog("system", "3003：系统额度为0", "system")
					}
				}
			}

			queryResultChan <- queryResult
		default:
			if len(queryChan) == 0 {
				waitGroup.Done()
				return
			}
		}
	}
}

func makeQueryRequest(query *model.Query, apiCode string) *http.Request {

	targetURL := "http://fpcyapi.market.alicloudapi.com/invoice/query?" +
		"fpdm=" + query.Fpdm + "&" + "fphm=" + query.Fphm + "&" + "kprq=" + query.Kprq + "&" +
		"je="   + query.Je   + "&" + "jym=" + query.Jym

	req, _ := http.NewRequest("GET", targetURL, nil) // 新建请求

	req.Header.Set("Authorization", "APPCODE " + apiCode)

	return req
}

func getQueryArrayJSON(ctx iris.Context, info string) (*model.QueryArray, bool) {

	queryArrayItem := &model.QueryArray{}

	err := ctx.ReadJSON(queryArrayItem)

	if err != nil{

		_, _ = ctx.JSON(iris.Map{
			"status":  "failed",
			"message": info + "失败" + err.Error(),
		})

		return queryArrayItem, false
	}

	return queryArrayItem, true
}

func queryResultToResultDBs(queryResult *model.QueryResult) []*model.ResultItem {

	var ret []*model.ResultItem

	resultDB := &model.ResultItem{}

	resultDB.Ensured = "<span class='my-failed'>未确认</span>"
	resultDB.Sealed = "0"
	resultDB.RespCode = queryResult.RespCode
	resultDB.RespMsg = queryResult.RespMsg

	if queryResult.Data.Qd == "0" {
		resultDB.Qd = "有"
	} else if queryResult.Data.Qd == "1" {
		resultDB.Qd = "没有"
	}

	resultDB.Fpdm = queryResult.Data.Fpdm
	resultDB.Fphm = queryResult.Data.Fphm
	resultDB.Kprq = queryResult.Data.Kprq
	resultDB.YzmSj = queryResult.Data.YzmSj

	if queryResult.Data.Fpzt == "0" {
		resultDB.Fpzt = "正常"
	} else if queryResult.Data.Fpzt == "2" {
		resultDB.Fpzt = "作废"
	}

	if queryResult.Data.Fxqy == "0" {
		resultDB.Fxqy = "正常"
	} else if queryResult.Data.Fxqy == "1" {
		resultDB.Fxqy = "异常"
	}

	if queryResult.Data.Fplx == "01" {
		resultDB.Fplx = "增值税专票"
	} else if queryResult.Data.Fplx == "03" {
		resultDB.Fplx = "机动车发票"
	} else if queryResult.Data.Fplx == "04" {
		resultDB.Fplx = "增值税发票"
	} else if queryResult.Data.Fplx == "10" {
		resultDB.Fplx = "电子发票"
	} else if queryResult.Data.Fplx == "11" {
		resultDB.Fplx = "卷式发票"
	} else if queryResult.Data.Fplx == "14" {
		resultDB.Fplx = "通行费发票"
	} else if queryResult.Data.Fplx == "15" {
		resultDB.Fplx = "二手车发票"
	} else {
		resultDB.Fplx = "未知"
	}

	resultDB.Jqbm = queryResult.Data.Jqbm
	resultDB.Jym = queryResult.Data.Jym
	resultDB.GfName = queryResult.Data.GfName
	resultDB.GfNsrsbh = queryResult.Data.GfNsrsbh
	resultDB.GfAddressTel = queryResult.Data.GfAddressTel
	resultDB.GfBankZh = queryResult.Data.GfBankZh
	resultDB.JshjL = queryResult.Data.JshjL
	resultDB.SfName = queryResult.Data.SfName
	resultDB.SfNsrsbh = queryResult.Data.SfNsrsbh
	resultDB.SfAddressTel = queryResult.Data.SfAddressTel
	resultDB.SfBankZh = queryResult.Data.SfBankZh
	resultDB.Bz = queryResult.Data.Bz
	resultDB.JshjU = queryResult.Data.JshjU
	resultDB.QueryTime = time.Now().Format("2006-1-2")

	var totalJe float64 = 0
	var totalSe float64 = 0

	for _, goodsItem := range queryResult.Data.ZpList {

		resultItem := resultDB.DeepCopy()

		resultItem.MxName = goodsItem.MxName
		resultItem.Ggxh = goodsItem.Ggxh
		resultItem.Unit = goodsItem.Unit
		resultItem.Price = goodsItem.Price
		resultItem.Je = goodsItem.Je
		resultItem.Sl = goodsItem.Sl
		resultItem.Se = goodsItem.Se
		resultItem.Num = goodsItem.Num

		ret = append(ret, resultItem)

		itemJe, _ := strconv.ParseFloat(goodsItem.Je, 64)
		itemSe, _ := strconv.ParseFloat(goodsItem.Se, 64)

		totalJe += itemJe
		totalSe += itemSe
	}

	for _, item := range ret {
		item.TotalJe = totalJe
		item.TotalSe = totalSe
	}

	return ret
}

func judgeQuerySuccess(code string) bool {
	if code == "2210" || code == "2213" || code == "2215" || code == "2206" {
		return true
	}

	return false
}
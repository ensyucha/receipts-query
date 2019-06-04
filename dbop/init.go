package dbop

import (
	"database/sql"
	"github.com/Unknwon/goconfig"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

var db *sql.DB

var rateG float64
var portG string
var enableLogG bool

func init() {

	cfgByte, err := ioutil.ReadFile("./conf/conf.ini")

	if err != nil {
		panic("获取配置文件 conf.ini 错误：" + err.Error())
	}

	cfg, err := goconfig.LoadFromData(cfgByte)

	if err != nil{
		panic("读取配置文件 conf.ini 错误：" + err.Error())
	}

	port, err := cfg.GetValue("setting","port")

	if err != nil {
		panic("读取配置文件 conf.ini 错误：" + err.Error())
	}

	portNum, err := strconv.Atoi(port)

	if err != nil || portNum < 1024 || portNum > 65535{
		panic("请设置conf.ini中port为数字，且大于1024并少于65536")
	}

	portG = ":" + port

	rate, err := cfg.GetValue("setting","rate")

	if err != nil || errorRate(rate) {
		panic("请设置conf.ini中rate，可选值：[0.2, 0.4, 0.6, 0.8, 1.0]")
	}

	rateNum, err := strconv.ParseFloat(rate, 64)

	if err != nil {
		panic(err)
	}

	rateG = rateNum

	enableLog, err := cfg.GetValue("setting","log")

	if err != nil || (enableLog != "enable" && enableLog != "disable") {
		panic("请设置conf.ini中log，可选值：[enable, disable]")
	}

	if enableLog == "enable" {
		enableLogG = true
	} else if enableLog == "disable" {
		enableLogG = false
	}

	//////////////////////////////////////////////////////////////////////

	dbUsername, err := cfg.GetValue("database", "dbusername")

	if err != nil {
		panic("请设置conf.ini中dbusername为Mysql账号")
	}

	dbPassword, err := cfg.GetValue("database","dbpassword")

	if err != nil {
		panic("请设置conf.ini中dbpassword为Mysql密码")
	}

	dbHost, err := cfg.GetValue("database", "dbhost")

	if err != nil {
		panic("请设置conf.ini中dbhost为Mysql服务器地址")
	}

	dbPort, err := cfg.GetValue("database", "dbport")

	if err != nil {
		panic("请设置conf.ini中dbPort为Mysql端口")
	}

	/////////////////////////////////////////////////////

	cmd := "mysql -u" + dbUsername + " -p" + dbPassword + " receipts < ./conf/database.sql"
	_ = exec.Command("CMD", "/C", cmd).Run()

	dbDriver := dbUsername+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/receipts?charset=utf8"

	db, err = sql.Open("mysql", dbDriver)

	if err != nil {
		panic("数据库打开失败：" + err.Error())
	}

	systemInfo, err := GetSystemInfo()

	if err != nil {
		panic("获取系统信息失败：" + err.Error())
	}

	unusedUsage, err := UCGetUnusedUsage()

	if err != nil {
		panic("获取系统未分余额失败：" + err.Error())
	}

	/////////////////////////

	log.Println("管理员账号 : admin")
	log.Println("管理员密码 :", systemInfo.Password)
	log.Println("ApiCode    :", systemInfo.ApiCode)
	log.Println("闲置额度   :", unusedUsage)
	log.Println("查询速率   :", rateG)

	WriteLog("system", "服务器启动", "manager")
}

func errorRate(rate string) bool {
	return rate != "0.2" && rate != "0.4" && rate != "0.6" && rate != "0.8" && rate != "1.0"
}

func GetPort() string {
	return portG
}

func GetIP() string {

	address, err := net.InterfaceAddrs()

	if err != nil {
		panic("获取本机IP地址失败：" + err.Error())
	}

	var ip = ""

	for _, address := range address {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if strings.HasPrefix(ipnet.IP.String(), "192.168") ||
					strings.HasPrefix(ipnet.IP.String(), "110") {
					ip = ipnet.IP.String()
					break
				}
			}

		}
	}

	if ip == "" {
		panic("获取本机IP地址失败")
	}

	return ip
}

func GetRate() float64 {
	return rateG
}
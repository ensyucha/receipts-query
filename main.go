package main

import (
	"github.com/kataras/iris"
	"log"
	"receipts/dbop"
	_ "receipts/dbop"
	"receipts/server"
)

func main() {

	// 新建服务器
	app := server.NewApp()

	log.Println("监听地址   : http://" + dbop.GetIP() + dbop.GetPort() + "\n")

	// 启动服务器，监听端口 33333
	log.Fatal(app.Run(iris.Addr(dbop.GetPort()), iris.WithoutStartupLog))
}
package main

import (
	"DPMQ/config"
	"DPMQ/persistent"
	"DPMQ/router"
	"DPMQ/server"
	_ "fmt"
	_ "github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "net/http"
)

/**
 * main
 */

func main() {

	//加载配置
	config.InitConfig()

	//加载日志模块
	server.InitLog()

	//初始化消息队列
	server.InitMQ()

	//加载文件读写模块
	persistent.InitFileRW()

	persistent.InitRecovery()

	persistent.InitPers()

	//初始化消费者客户端连接模块
	server.InitConsumersConn()

	//设置路由
	r := router.SetupRouters()

	//获取端口号
	port := viper.GetString("server.port")
	_ = r.Run(":" + port)
}

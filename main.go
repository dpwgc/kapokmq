package main

import (
	"DPMQ/config"
	"DPMQ/routers"
	"DPMQ/servers"
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

	//初始化消息队列
	servers.InitMQ()

	//初始化消费者客户端连接模块
	servers.InitConsumersConn()

	//设置路由
	r := routers.SetupRouters()

	//获取端口号
	port := viper.GetString("server.port")
	_ = r.Run(":" + port)
}

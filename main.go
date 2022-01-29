package main

import (
	"fmt"
	_ "github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"kapokmq/cluster"
	"kapokmq/config"
	"kapokmq/persistent"
	"kapokmq/router"
	"kapokmq/server"
	_ "net/http"
)

func main() {

	//加载配置
	config.InitConfig()

	//加载常规日志模块
	server.InitLog()

	//初始化消息队列
	server.InitMQ()

	//初始化Gossip集群连接模块
	cluster.InitCluster()

	//加载文件读写模块
	persistent.InitFileRW()

	//加载数据恢复模块
	persistent.InitRecovery()

	//加载WAL持久化模块
	server.InitWAL()

	//加载持久化模块
	persistent.InitPers()

	//加载消息检查模块
	server.InitCheck()

	//初始化消费者客户端连接模块
	server.InitConsumersConn()

	//初始化路由
	r := router.InitRouters()

	//获取端口号
	port := viper.GetString("server.port")
	err := r.Run(fmt.Sprintf("%s%s", ":", port))
	if err != nil {
		panic(err)
	}
}

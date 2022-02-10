package main

import (
	"fmt"
	_ "github.com/gin-gonic/gin"
	"kapokmq/cluster"
	"kapokmq/config"
	"kapokmq/memory"
	"kapokmq/mqLog"
	"kapokmq/persistent"
	"kapokmq/router"
	"kapokmq/server"
	"kapokmq/syncConn"
	_ "net/http"
)

func main() {

	//加载配置
	config.InitConfig()

	//加载常规日志模块
	mqLog.InitLog()

	//初始化消息队列
	memory.InitMQ()

	//初始化Gossip集群连接模块
	cluster.InitCluster()

	//加载文件读写模块
	persistent.InitFileRW()

	//加载数据恢复模块
	persistent.InitRecovery()

	//加载WAL持久化模块
	mqLog.InitWAL()

	//加载持久化模块
	persistent.InitPers()

	//加载消息检查模块
	server.InitCheck()

	//初始化消费者客户端连接模块
	server.InitConsumersConn()

	//初始化主从同步模块
	syncConn.InitSync()

	//初始化路由
	r := router.InitRouters()

	//获取端口号
	port := config.Get.Server.Port
	err := r.Run(fmt.Sprintf("%s%s", ":", port))
	if err != nil {
		panic(err)
	}
}

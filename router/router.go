package router

import (
	"github.com/gin-gonic/gin"
	"kapokmq/cluster"
	"kapokmq/console"
	"kapokmq/middleware"
	"kapokmq/server"
	"net/http"
)

// InitRouters 初始化路由
func InitRouters() (r *gin.Engine) {

	r = gin.Default()

	//前端页面
	r.Static("/view", "view")
	r.LoadHTMLGlob("view/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	//跨域配置
	r.Use(middleware.Cors())

	//控制台接口（http post请求，用于查看消息队列的基本信息）
	consoleApi := r.Group("/Console")
	consoleApi.Use(middleware.Safe)
	{
		consoleApi.POST("/Ping", console.Ping)
		consoleApi.POST("/GetClients", console.GetClients)
		consoleApi.POST("/GetConfig", console.GetConfig)
		consoleApi.POST("/GetMessageList", console.GetMessageList)
		consoleApi.POST("/GetMessage", console.GetMessage)
		consoleApi.POST("/DelMessage", console.DelMessage)
		consoleApi.POST("/CountMessage", console.CountMessage)
		consoleApi.POST("/GetNodes", cluster.GetNodes)
	}

	//生产者接口（http post请求，用于接收生产者客户端发送的消息）
	producerApi := r.Group("/Producer")
	producerApi.Use(middleware.Safe)
	{
		producerApi.POST("/Send", server.ProducerSend)
	}

	//消费者连接（websocket连接，用于推送消息到消费者客户端）ws://127.0.0.1:8011/Consumers/Conn/test_topic/1
	consumersConn := r.Group("/Consumers")
	{
		consumersConn.GET("/Conn/:topic/:consumerId", server.ConsumersConn)
	}
	//生产者连接（websocket连接，用于发送消息到消息队列）ws://127.0.0.1:8011/Producers/Conn/test_topic/1
	producersConn := r.Group("/Producers")
	{
		producersConn.GET("/Conn/:topic/:producerId", server.ProducersConn)
	}
	return
}

package router

import (
	"github.com/gin-gonic/gin"
	"kapokmq/middleware"
	"kapokmq/server"
	"net/http"
)

/**
 * 路由
 */

func InitRouters() (r *gin.Engine) {

	r = gin.Default()
	r.Static("/view", "view")
	r.LoadHTMLGlob("view/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Use(middleware.Cors())

	//控制台接口（http post请求，用于查看消息队列的基本信息）
	console := r.Group("/Console")
	console.Use(middleware.Safe)
	{
		console.POST("/Ping", server.Ping)
		console.POST("/GetConsumers", server.GetConsumers)
		console.POST("/GetProducers", server.GetProducers)
		console.POST("/GetConfig", server.GetConfig)
		console.POST("/GetMessageList", server.GetMessageList)
		console.POST("/GetMessageEasy", server.GetMessageEasy)
		console.POST("/CountMessage", server.CountMessage)
	}

	//生产者接口（http post请求，用于接收生产者客户端发送的消息）
	producer := r.Group("/Producer")
	producer.Use(middleware.Safe)
	{
		producer.POST("/Send", server.ProducerSend)
	}

	//消费者连接（websocket连接，用于推送消息到消费者客户端）ws://127.0.0.1:8011/Consumers/Conn/test_topic/1
	consumers := r.Group("/Consumers")
	{
		consumers.GET("/Conn/:topic/:consumerId", server.ConsumersConn)
	}
	//生产者连接（websocket连接，用于发送消息到消息队列）ws://127.0.0.1:8011/Producers/Conn/test_topic/1
	producers := r.Group("/Producers")
	{
		producers.GET("/Conn/:topic/:producerId", server.ProducersConn)
	}
	return
}

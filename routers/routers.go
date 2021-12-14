package routers

import (
	"DPMQ/servers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/**
 * 路由
 */

func SetupRouters() (r *gin.Engine) {

	r = gin.Default()
	r.Use(Cors())

	//控制台接口（http post请求，用于查看消息队列的基本信息）
	console := r.Group("/Console")
	{
		console.GET("/GetConsumers", servers.GetConsumers)
		console.GET("/GetConfig", servers.GetConfig)
		console.GET("/GetMessageLog", servers.GetMessageList)
	}

	//生产者接口（http post请求，用于接收生产者客户端发送的消息）
	r.POST("/ProducerSend", servers.ProducerSend)

	//消费者连接（websocket连接，用于推送消息到消费者客户端）
	r.GET("/ConsumersConn/:topic/:consumerId", servers.ConsumersConn)

	return
}

//跨域设置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}

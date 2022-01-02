package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Safe 安全验证中间件
func Safe(c *gin.Context) {

	secretKey := c.GetHeader("secretKey")

	if secretKey != viper.GetString("mq.secretKey") {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "Secret key matching error",
		})
		c.Abort()
	}
	return
}

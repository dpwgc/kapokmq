package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//安全验证中间件
func SafeMiddleWare(c *gin.Context) {

	secretKey := c.GetHeader("secretKey")

	if secretKey != viper.GetString("mq.secretKey") {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "key matching error",
		})
		c.Abort()
	}
	return
}

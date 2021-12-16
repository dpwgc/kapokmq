package middlewares

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
			"msg":  "密钥匹配出错",
		})
		c.Abort()
	}
	return
}

package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

/**
 * 控制台服务
 */

//获取全部消费者客户端集合
func GetConsumers(c *gin.Context) {

	resMap := make(map[string]string, len(consumers)+10)

	//遍历消费者客户端集合
	for key := range consumers {

		topic := strings.Split(key, "|")[0]
		consumerId := strings.Split(key, "|")[1]

		resMap[topic] = consumerId
	}

	c.JSON(0, gin.H{
		"code": 0,
		"data": resMap,
	})
}

//获取消息队列详细配置
func GetConfig(c *gin.Context) {

	resMap := make(map[string]interface{}, 3)
	resMap["messageChanBuffer"] = viper.GetInt("mq.messageChanBuffer")
	resMap["pushMessagesSpeed"] = pushMessagesSpeed
	resMap["sendRetryCount"] = sendRetryCount

	c.JSON(0, gin.H{
		"code": 0,
		"data": resMap,
	})
}

package servers

import (
	"DPMQ/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

/**
 * 控制台服务
 */

//获取全部消费者客户端集合
func GetConsumers(c *gin.Context) {

	resMap := make(map[string]string)

	//遍历消费者客户端集合
	for key := range consumers {

		topic := strings.Split(key, "|")[0]
		consumerId := strings.Split(key, "|")[1]

		resMap[consumerId] = topic
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
	resMap["pushMessagesSpeed"] = viper.GetInt("mq.pushMessagesSpeed")
	resMap["sendRetryCount"] = viper.GetInt("mq.sendRetryCount")

	c.JSON(0, gin.H{
		"code": 0,
		"data": resMap,
	})
}

//获取指定状态的消息记录列表
func GetMessageList(c *gin.Context) {

	status, _ := c.GetPostForm("status")
	intStatus, _ := strconv.Atoi(status)

	var resArr []models.Message

	//遍历消息记录集合
	for _, message := range messageList {

		if message.Status == intStatus {
			resArr = append(resArr, message)
		}
	}

	c.JSON(0, gin.H{
		"code": 0,
		"data": resArr,
	})
}

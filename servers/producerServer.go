package servers

import (
	"DPMQ/models"
	"DPMQ/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

/**
 * 接收消息生产者发送过来的消息
 */

func ProducerSend(c *gin.Context) {

	messageData, _ := c.GetPostForm("messageData")
	topic, _ := c.GetPostForm("topic")
	//判断topic字符串是否含有字符“|”，如果有，则返回错误信息，避免影响后续字符串切割操作
	if strings.Contains(topic, "|") {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "topic不能包含字符“|”",
		})
		return
	}

	message := models.Message{}
	message.MessageCode = utils.CreateCode(messageData)
	message.MessageData = messageData
	message.Topic = topic
	message.Status = 0
	message.CreateTime = utils.GetLocalDateTime()

	//把消息写入消息通道
	messageChan <- message

	//发送成功，返回消息标识码
	c.JSON(0, gin.H{
		"code": 0,
		"msg":  message.MessageCode,
	})
}

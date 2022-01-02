package server

import (
	"DPMQ/model"
	"DPMQ/utils"
	"github.com/gin-gonic/gin"
)

// ProducerSend 接收消息生产者发送过来的消息
func ProducerSend(c *gin.Context) {

	messageData, _ := c.GetPostForm("messageData")
	topic, _ := c.GetPostForm("topic")

	if len(messageData) == 0 || len(topic) == 0 {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "Required data cannot be empty",
		})
		return
	}

	message := model.Message{}
	message.MessageCode = utils.CreateCode(messageData)
	message.MessageData = messageData
	message.Topic = topic
	message.Status = -1
	message.CreateTime = utils.GetLocalDateTimestamp()

	//把消息写入消息通道
	MessageChan <- message
	//将消息记录到消息列表
	MessageList.Store(message.MessageCode, message)

	//发送成功，返回消息标识码
	c.JSON(0, gin.H{
		"code": 0,
		"msg":  message.MessageCode,
	})
}

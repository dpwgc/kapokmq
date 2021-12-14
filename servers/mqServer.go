package servers

import (
	"DPMQ/models"
	"github.com/spf13/viper"
)

func InitMQ() {
	messageChanBuffer = viper.GetInt("mq.messageChanBuffer")
	//消息通道初始化
	messageChan = make(chan models.Message, messageChanBuffer)
}

/**
 * 消息队列服务
 */

//消息通道缓冲区大小
var messageChanBuffer int

//消息通道，用于存放待消费的消息(有缓冲区)
var messageChan chan models.Message

//消息列表，存放所有消息记录
var messageList []models.Message

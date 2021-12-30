package servers

import (
	"DPMQ/models"
	"github.com/spf13/viper"
)

func InitMQ() {
	//消息通道缓冲区大小
	messageChanBuffer := viper.GetInt("mq.messageChanBuffer")
	//消息通道初始化
	MessageChan = make(chan models.Message, messageChanBuffer)
}

/**
 * 消息队列服务
 */

//消息通道，用于存放待消费的消息(有缓冲区)
var MessageChan chan models.Message

//消息列表，存放所有消息记录
var MessageList []models.Message

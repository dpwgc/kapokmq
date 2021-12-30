package server

import (
	"DPMQ/model"
	"github.com/spf13/viper"
)

func InitMQ() {
	//消息通道缓冲区大小
	messageChanBuffer := viper.GetInt("mq.messageChanBuffer")
	//消息通道初始化
	MessageChan = make(chan model.Message, messageChanBuffer)
}

/**
 * 消息队列服务
 */

//消息通道，用于存放待消费的消息(有缓冲区)
var MessageChan chan model.Message

//消息列表，存放所有消息记录
var MessageList []model.Message

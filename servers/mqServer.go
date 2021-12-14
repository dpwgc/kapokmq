package servers

import (
	"DPMQ/models"
	"github.com/spf13/viper"
	"time"
)

func InitMQ() {
	messageChanBuffer = viper.GetInt("mq.messageChanBuffer")
	messageChan = make(chan models.Message, messageChanBuffer)
	//将消息持久化刷盘
	go SetLog()
}

/**
 * 消息队列服务
 */

//消息通道缓冲区大小
var messageChanBuffer int

//消息通道，用于存放待消费的消息(有缓冲区)
var messageChan chan models.Message

//消息列表，存放所有消息，用于持久化
var messageList []models.Message

//消息记录持久化 TODO
func SetLog() {

	for {
		time.Sleep(time.Second * 1)
	}
}

package servers

import (
	"DPMQ/models"
	"github.com/spf13/viper"
)

//消息通道，用于存放消息(有缓冲区)
var messageChan = make(chan models.Message, viper.GetInt("mq.messageChanBuffer"))

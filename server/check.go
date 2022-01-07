package server

import (
	"github.com/spf13/viper"
	"kapokmq/model"
	"time"
)

// InitCheck 消息检查-消息重推与过期消息清除模块
func InitCheck() {

	isRePush = viper.GetInt("mq.isRePush")
	checkSpeed = viper.GetInt("mq.checkSpeed")

	isClean = viper.GetInt("mq.isClean")
	cleanTime = viper.GetInt64("mq.cleanTime")

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(checkSpeed))
			checkMessage()
		}
	}()
}

var isRePush int

var checkSpeed int

var isClean int
var cleanTime int64

//检查消息-用于：清理过期消息-重推消息-延时消息推送
func checkMessage() {

	//获取当前时间戳
	ts := time.Now().Unix()

	MessageList.Range(func(key, message interface{}) bool {

		msg := message.(model.Message)

		//该消息超出记录时间限制，彻底删除该消息
		if ts-msg.CreateTime > cleanTime && isClean == 1 {
			MessageList.Delete(key)
			return true
		}

		//如果该消息未推送
		if msg.Status == 0 && isRePush == 1 {

			//将消息标记为无状态，重新推送该消息
			msg.Status = -1
			MessageChan <- msg
			return true
		}
		return true
	})
}

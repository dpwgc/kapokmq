package server

import (
	"github.com/spf13/viper"
	"kapokmq/model"
	"time"
)

// InitCheck 消息检查
func InitCheck() {

	isRePush = viper.GetInt("mq.isRePush")
	checkSpeed = viper.GetInt("mq.checkSpeed")

	isClean = viper.GetInt("mq.isClean")
	cleanTime = viper.GetInt64("mq.cleanTime")

	pushRetryTime = viper.GetInt64("mq.pushRetryTime")

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

var pushRetryTime int64

//检查消息-用于：清理过期消息-重推消息-延时消息推送
func checkMessage() {

	//获取当前时间戳
	ts := time.Now().Unix()

	MessageList.Range(func(key, msg interface{}) bool {

		message := msg.(model.Message)

		//该消息超出记录时间限制，且mq开启了自动清理功能，则彻底删除该消息
		if ts-message.CreateTime > cleanTime && isClean == 1 {
			MessageList.Delete(key)
			return true
		}

		//如果检查到延时消息
		if message.DelayTime > 0 {
			//如果还未到投送时间
			if message.CreateTime+message.DelayTime > ts {
				//等待推送
				return true
			} else {
				//推送消息
				MessageChan <- message
				return true
			}
		}

		//如果该消息到达重推时间，但仍未被确认消费，且mq开启了重推功能
		if message.Status == -1 && message.DelayTime+pushRetryTime < ts-message.CreateTime && isRePush == 1 {

			//延长消费时间
			message.DelayTime = message.DelayTime + pushRetryTime
			//更新该消息
			MessageList.Store(message.MessageCode, message)
			//重新推送
			MessageChan <- message
		}

		//如果是推送失败的消息，且mq开启了重推功能
		if message.Status == 0 && isRePush == 1 {

			//将消息标记为待消费状态，重新推送该消息
			message.Status = -1
			//更新该消息
			MessageList.Store(message.MessageCode, message)
			//重新推送
			MessageChan <- message
			return true
		}
		return true
	})
}

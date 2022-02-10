package server

import (
	"kapokmq/config"
	"kapokmq/memory"
	"kapokmq/model"
	"time"
)

// InitCheck 消息检查
func InitCheck() {

	isRePush = config.Get.Mq.IsRePush
	checkSpeed = config.Get.Mq.CheckSpeed

	isClean = config.Get.Mq.IsClean
	cleanTime = config.Get.Mq.CleanTime

	pushRetryTime = config.Get.Mq.PushRetryTime

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

	memory.MessageList.Range(func(key, msg interface{}) bool {

		message := msg.(model.Message)

		//该消息超出记录时间限制，且mq开启了自动清理功能，则彻底删除该消息
		if ts-message.CreateTime > cleanTime && isClean == 1 {
			memory.MessageList.Delete(key)
			return true
		}

		//如果检查到延时消息
		if message.Status == 0 {
			//如果还未到投送时间
			if message.CreateTime+message.DelayTime > ts {
				//等待推送
				return true
			} else {
				//将消息状态改为未消费
				message.Status = -1
				//更新该消息
				memory.MessageList.Store(message.MessageCode, message)
				//推送消息
				memory.MessageChan <- message
				return true
			}
		}

		//如果该消息到达重推时间，但仍未被确认消费，且mq开启了重推功能
		if message.Status == -1 && message.DelayTime+pushRetryTime < ts-message.CreateTime && isRePush == 1 {

			//将消息的推送时间设为当前时间，即将消息超时消费时间阈值移到当前时间的后{pushRetryTime}秒
			message.DelayTime = ts - message.CreateTime
			//更新该消息
			memory.MessageList.Store(message.MessageCode, message)
			//重新推送
			memory.MessageChan <- message
			return true
		}

		return true
	})
}

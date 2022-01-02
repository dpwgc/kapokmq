package server

import (
	"DPMQ/model"
	"github.com/spf13/viper"
	"time"
)

// InitCheck 消息检查-消息重推与过期消息清除模块
func InitCheck() {

	isRePush = viper.GetInt("mq.isRePush")
	checkSpeed = viper.GetInt("mq.checkSpeed")
	checkCount = viper.GetInt("mq.checkCount")

	isClean = viper.GetInt("mq.isClean")
	cleanTime = viper.GetInt64("mq.cleanTime")

	//既不开启自动清理功能，也不开启重推功能
	if isClean == 0 && isRePush == 0 {
		return
	}

	go func() {
		cnt := 0
		for {
			if cnt == checkCount {
				//消息重推的时间间隔（每重推{rePushCount}条消息，间隔一段时间）
				time.Sleep(time.Second * time.Duration(checkSpeed))
				cnt = 0
			}
			checkMessage()
			cnt++
		}
	}()
}

var isRePush int

var checkCount int
var checkSpeed int

var isClean int
var cleanTime int64

//检查消息
func checkMessage() {

	//获取当前时间戳
	ts := time.Now().Unix()

	MessageList.Range(func(key, message interface{}) bool {

		//该消息超出记录时间限制，彻底删除该消息
		if ts-message.(model.Message).CreateTime > cleanTime && isClean == 1 {
			MessageList.Delete(key)
			return true
		}

		//如果该消息未推送
		if message.(model.Message).Status == 0 && isRePush == 1 {
			//重新推送该消息
			MessageChan <- message.(model.Message)
			return true
		}
		return true
	})
}

package server

import (
	"DPMQ/model"
	"github.com/spf13/viper"
	"time"
)

/**
 * 消息重推与过期消息清除模块
 */

func InitRePush() {
	rePushSpeed = viper.GetInt("mq.rePushSpeed")
	rePushCount = viper.GetInt("mq.rePushCount")
	clearTime = viper.GetInt64("mq.clearTime")
	go func() {
		cnt := 0
		for {
			if cnt == rePushCount {
				//消息重推的时间间隔（每重推{rePushCount}条消息，间隔一段时间）
				time.Sleep(time.Second * time.Duration(rePushSpeed))
				cnt = 0
			}
			rePushMessage()
			cnt++
		}
	}()
}

var rePushSpeed int
var rePushCount int
var clearTime int64

func rePushMessage() {

	//获取当前时间戳
	ts := time.Now().Unix()

	MessageList.Range(func(key, message interface{}) bool {

		//该消息超出记录时间限制，彻底删除该消息
		if ts-message.(model.Message).CreateTime > clearTime {
			MessageList.Delete(key)
			return true
		}

		//如果该消息未推送
		if message.(model.Message).Status == 0 {
			//重新推送该消息
			MessageChan <- message.(model.Message)
			return true
		}
		return true
	})
}

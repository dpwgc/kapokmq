package model

import (
	"encoding/json"
	"kapokmq/utils"
)

// SendMessage 生产者发送的消息模板
type SendMessage struct {
	MessageData string //消息内容（一般为JSON格式的字符串）
	DelayTime   int64  //延迟推送时间（单位：秒）
}

// Message 消息模板
type Message struct {
	MessageCode  string //消息标识码
	MessageData  string //消息内容（一般为JSON格式的字符串）
	Topic        string //消息所属主题
	CreateTime   int64  //消息创建时间
	ConsumedTime int64  //消息被消费时间
	DelayTime    int64  //延迟推送时间（单位：秒）
	Status       int    //消息状态（-1：待消费。0：未到推送时间的延时消息。1：已消费）
}

// NewMessage 生成消息模板（消息队列接收生产者消息时使用）
func NewMessage(topic string, data []byte) (*Message, error) {

	send := SendMessage{}

	//解析json字符串，获取生产者客户端发送的消息内容和延时推送时间
	err := json.Unmarshal(data, &send)
	if err != nil {
		return nil, err
	}

	//生成消息模板
	message := Message{}
	message.MessageCode = utils.CreateCode(send.MessageData)
	message.MessageData = send.MessageData
	message.Topic = topic

	//如果是延时消息
	if send.DelayTime > 0 {
		message.Status = 0
	} else {
		message.Status = -1
	}

	message.CreateTime = utils.GetLocalDateTimestamp()
	message.DelayTime = send.DelayTime

	return &message, err
}

package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kapokmq/model"
	"kapokmq/utils"
	"strconv"
	"sync"
)

// ProducerSend 以HTTP请求方式接收消息
// ProducerSend 接收消息生产者发送过来的消息
func ProducerSend(c *gin.Context) {

	messageData, _ := c.GetPostForm("messageData")
	topic, _ := c.GetPostForm("topic")
	delayTime, _ := c.GetPostForm("delayTime")

	intDelayTime, err := strconv.ParseInt(delayTime, 10, 64)
	if err != nil {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  err,
		})
		return
	}

	if len(messageData) == 0 || len(topic) == 0 {
		c.JSON(-1, gin.H{
			"code": -1,
			"msg":  "Required data cannot be empty",
		})
		return
	}

	message := model.Message{}
	message.MessageCode = utils.CreateCode(messageData)
	message.MessageData = messageData
	message.Topic = topic
	message.Status = -1
	message.CreateTime = utils.GetLocalDateTimestamp()
	message.DelayTime = intDelayTime

	//把消息写入消息通道
	MessageChan <- message
	//将消息记录到消息列表
	MessageList.Store(message.MessageCode, message)

	//发送成功，返回消息标识码
	c.JSON(0, gin.H{
		"code": 0,
		"msg":  message.MessageCode,
	})
}

// Producers 连接的消费者客户端,把每个消费者都放进来。Key为{topic}|{consumerId}，topic与consumerId两者之间用字符”|“分隔。Value为websocket连接
var Producers = make(map[model.Producer]*websocket.Conn)

// ProducersConn 以WebSocket连接方式接收消息
// ProducersConn 生产者连接，接收生产者发送过来的消息
func ProducersConn(c *gin.Context) {

	//消息所属主题
	topic := c.Param("topic")
	//消费者id
	producerId := c.Param("producerId")
	//消费者ip地址
	producerIp := c.ClientIP()

	//升级get请求为webSocket协议
	ws, err := UpGrader.Upgrade(c.Writer, c.Request, nil)

	//登录验证
	for {
		//连接成功，等待消费者客户端输入访问密钥
		err = ws.WriteMessage(1, []byte("Please enter the secret key"))
		if err != nil {
			Loger.Println(err)
			return
		}

		//读取ws中的数据，获取访问密钥
		_, sk, err := ws.ReadMessage()
		if err != nil {
			Loger.Println(err)
			return
		}
		if string(sk) == secretKey {
			//访问密钥匹配成功
			err = ws.WriteMessage(1, []byte("Secret key matching succeeded"))
			if err != nil {
				Loger.Println(err)
				return
			}
			break
		}

		//访问密钥匹配失败
		err = ws.WriteMessage(1, []byte("Secret key matching error"))
		if err != nil {
			Loger.Println(err)
			return
		}
	}

	key := model.Producer{}
	key.ProducerId = producerId
	key.Topic = topic
	key.ProducerIp = producerIp
	key.JoinTime = utils.GetLocalDateTimestamp()

	lock := sync.RWMutex{}

	//将当前连接的消费者放入map中
	lock.RLock()
	Producers[key] = ws
	lock.RUnlock()

	if err != nil {
		Loger.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		delete(Producers, key) //删除map中的生产者
		err = ws.Close()
		if err != nil {
			Loger.Println(err)
		}
	}(ws)

	for {
		//读取websocket中的数据
		_, data, err := ws.ReadMessage()
		if err != nil {
			Loger.Println(err)
			return
		}
		s := model.SendMessage{}
		//解析json字符串，获取生产者客户端发送的消息内容和延时推送时间
		err = json.Unmarshal(data, &s)
		if err != nil {
			Loger.Println(err)
			return
		}

		message := model.Message{}
		message.MessageCode = utils.CreateCode(s.MessageData)
		message.MessageData = s.MessageData
		message.Topic = topic
		message.Status = -1
		message.CreateTime = utils.GetLocalDateTimestamp()
		message.DelayTime = s.DelayTime

		//把消息写入消息通道
		MessageChan <- message
		//将消息记录到消息列表
		MessageList.Store(message.MessageCode, message)
	}
}

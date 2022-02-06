package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kapokmq/config"
	"kapokmq/model"
	"kapokmq/utils"
	"net/http"
	"sync"
	"time"
)

/**
 * 消费者连接模块
 */

//消费者客户端map锁
var cLock = sync.RWMutex{}

//安全密钥
var secretKey string

//每一批推送的消息数量
var pushCount int

//推送消息的速度(单批次消息推送间隔时间，单位：秒)
var pushMessagesSpeed int

//是否立即清除已确认消费的消息
var isCleanConsumed int

// Consumers 连接的消费者客户端，把每个消费者都放进来。Key为Consumer结构体，Value为websocket连接
var Consumers = make(map[model.Consumer]*websocket.Conn)

// UpGrader websocket跨域配置
var UpGrader = websocket.Upgrader{
	//跨域设置
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// InitConsumersConn 初始化消费者连接模块
func InitConsumersConn() {

	secretKey = config.Get.Mq.SecretKey
	pushCount = config.Get.Mq.PushCount
	isCleanConsumed = config.Get.Mq.IsCleanConsumed
	pushMessagesSpeed = config.Get.Mq.PushMessagesSpeed

	Loger.Println("Start pushServer")
	//启动消息推送协程，推送消息到各个消费者客户端
	go pushServer()
}

// ConsumersConn 消费者连接，监听消息队列内部各个主题消息的更新
func ConsumersConn(c *gin.Context) {

	//消息所属主题
	topic := c.Param("topic")
	//消费者id
	consumerId := c.Param("consumerId")
	//消费者ip地址
	consumerIp := c.ClientIP()

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

	//生成消费者客户端模板
	key := model.Consumer{}
	key.ConsumerId = consumerId
	key.Topic = topic
	key.ConsumerIp = consumerIp
	key.JoinTime = utils.GetLocalDateTimestamp()

	//将当前连接的消费者放入map中
	cLock.RLock()
	Consumers[key] = ws
	cLock.RUnlock()

	if err != nil {
		Loger.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		delete(Consumers, key) //删除map中的消费者
		err = ws.Close()
		if err != nil {
			Loger.Println(err)
		}
	}(ws)

	//ACK消息消费确认
	for {
		//读取ws中的数据，获取消费者客户端发来的确认消费信息，ack为messageCode
		_, ack, err := ws.ReadMessage()
		if err != nil {
			Loger.Println(err)
			return
		}

		//从Map中加载出消息
		msg, isOk := MessageList.Load(string(ack))
		if !isOk {
			continue
		}

		message := msg.(model.Message)

		//更新消息信息
		message.Status = 1
		message.ConsumedTime = utils.GetLocalDateTimestamp()

		//持久化：WAL写前日志
		if config.Get.Mq.IsPersistent == 2 {
			SetWAL(message)
		}

		//是否立即清除已被消费的消息
		if isCleanConsumed == 1 {
			MessageList.Delete(message.MessageCode)
			continue
		}

		//确认消费
		MessageList.Store(string(ack), message)
	}
}

//消息推送服务
func pushServer() {
	//获取推送模式
	pushType := config.Get.Mq.PushType
	cnt := 0
	for {
		if cnt == pushCount {
			//消息推送的时间间隔（每发送{sendCount}条消息，间隔一段时间）
			time.Sleep(time.Second * time.Duration(pushMessagesSpeed))
			cnt = 0
		}
		//选择对应的推送模式
		switch pushType {
		case 1:
			//使用点对点模式推送消息
			pushMessagesToOneConsumer()
			break
		case 2:
			//使用订阅/发布推送模式推送消息
			pushMessagesToConsumers()
			break
		default:
			return
		}
		cnt++
	}
}

//订阅/发布推送模式：并发推送消息到各个消费者客户端
func pushMessagesToConsumers() {

	//读取消息通道中的消息
	message := <-MessageChan

	//如果是延时消息
	if message.Status == 0 {
		//等待重推
		return
	}

	//控制通道
	controlChan := make(chan int, len(Consumers))

	//遍历消费者客户端集合
	for key, consumer := range Consumers {

		//多协程并发推送消息
		go func(key model.Consumer, consumer *websocket.Conn) {

			//字符串分割获取该消息所属主题
			topic := key.Topic

			//找到与该消息主题对应的客户端(相同的topic)
			if message.Topic == topic && len(message.Topic) > 0 && len(message.MessageCode) > 0 {

				//发送消息到消费者客户端
				err := consumer.WriteJSON(message)
				//如果连接异常
				if err != nil {
					//客户端关闭
					err = consumer.Close()
					if err != nil {
						Loger.Println(err)
					}
					//删除map中的客户端
					delete(Consumers, key)
				}
			}
			//向控制通道发送信息，表示该协程处理完毕
			controlChan <- 1
		}(key, consumer)
	}

	//待全部推送协程执行完成后，进入下一条消息的推送
	for range Consumers {
		//收到协程执行完毕的信息
		<-controlChan
	}
}

//点对点模式：随机推送消息到某个消费者客户端
func pushMessagesToOneConsumer() {

	//读取消息通道中的消息
	message := <-MessageChan

	//如果是延时消息
	if message.Status == 0 {
		//等待重推
		return
	}

	//遍历消费者客户端集合
	for key, consumer := range Consumers {

		//字符串分割获取该消息所属主题
		topic := key.Topic

		//找到与该消息主题对应的客户端(相同的topic)
		if message.Topic == topic && len(message.Topic) > 0 && len(message.MessageCode) > 0 {

			//发送消息到消费者客户端
			err := consumer.WriteJSON(message)
			//如果连接异常
			if err != nil {
				//客户端关闭
				err = consumer.Close()
				if err != nil {
					Loger.Println(err)
				}
				//删除map中的客户端
				delete(Consumers, key)
				//尝试投送给下一个消费者客户端
				continue
			}

			//点对点，发送给一个客户端后，直接返回
			return
		}
	}
}

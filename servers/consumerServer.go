package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

/**
 * 使用websocket将消息推送给各个消费者客户端
 */

//websocket跨域配置
var UpGrader = websocket.Upgrader{
	//跨域设置
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitConsumersConn() {
	sendCount = viper.GetInt("mq.sendCount")
	sendRetryCount = viper.GetInt("mq.sendRetryCount")
	pushMessagesSpeed = viper.GetInt("mq.pushMessagesSpeed")
	//启动消息推送协程，推送消息到各个消费者客户端
	go pushServer()
}

//每一批推送的消息数量
var sendCount int

//消息推送失败后的重试次数
var sendRetryCount int

//推送消息的速度(单批次消息推送间隔时间，单位：秒)
var pushMessagesSpeed int

//连接的消费者客户端,把每个消费者都放进来。Key为{topic}|{consumerId}，topic与consumerId两者之间用字符”|“分隔。Value为websocket连接
var consumers = make(map[string]*websocket.Conn)

//消费者连接，监听消息队列内部各个主题消息的更新
func ConsumersConn(c *gin.Context) {

	//消息所属主题
	topic := c.Param("topic")
	//判断topic字符串是否含有字符“|”，如果有，则断开连接，避免影响后续字符串切割操作
	if strings.Contains(topic, "|") {
		return
	}

	//消费者id
	consumerId := c.Param("consumerId")
	//判断consumerId字符串是否含有字符“|”，如果有，则断开连接，避免影响后续字符串切割操作
	if strings.Contains(consumerId, "|") {
		return
	}

	//升级get请求为webSocket协议
	ws, err := UpGrader.Upgrade(c.Writer, c.Request, nil)

	//将当前连接的消费者放入map中
	consumers[topic+"|"+consumerId] = ws

	if err != nil {
		delete(consumers, topic+"|"+consumerId) //删除map中的消费者
		return
	}
	defer ws.Close()

	for {
		//开个死循环，将连接挂起，保证连接不被断开
		time.Sleep(time.Second * 10)
	}
}

//消息推送服务
func pushServer() {
	cnt := 0
	for {
		if cnt == sendCount {
			//消息推送的时间间隔（每发送{sendCount}条消息，间隔一段时间）
			time.Sleep(time.Second * time.Duration(pushMessagesSpeed))
			cnt = 0
		}
		//推送消息
		pushMessagesToConsumers()
		cnt++
	}
}

//并发推送消息到各个消费者客户端
func pushMessagesToConsumers() {

	//如果没有消费者客户端，等待
	if len(consumers) == 0 {
		return
	}

	//读取消息通道中的消息
	message := <-messageChan

	//控制通道
	controlChan := make(chan int)

	//遍历消费者客户端集合
	for key, consumer := range consumers {

		//多协程并发推送消息
		go func(key string, consumer *websocket.Conn) {

			//字符串分割获取该消息所属主题
			topic := strings.Split(key, "|")[0]

			//找到与该消息主题对应的客户端(相同的topic)
			if message.Topic == topic && len(message.Topic) > 0 && len(message.MessageCode) > 0 {

				//重试机制
				for i := 0; i < sendRetryCount; i++ {
					//发送消息到消费者客户端
					err := consumer.WriteJSON(message)
					//如果发送成功
					if err == nil {
						//将消息标记为已确认状态
						message.Status = 1
						//记录到消息列表
						messageList = append(messageList, message)
						//结束循环
						break
					}
					//如果到达重试次数，但仍未发送成功
					if i == sendRetryCount-1 && err != nil {
						//客户端关闭
						consumer.Close()
						//删除map中的客户端
						delete(consumers, key)
					}
				}
			}
			//向控制通道发送信息，表示该协程处理完毕
			controlChan <- 1
		}(key, consumer)
	}

	//待全部推送协程执行完成后，进入下一条消息的推送
	for range consumers {
		//收到协程执行完毕的信息
		<-controlChan
	}
}

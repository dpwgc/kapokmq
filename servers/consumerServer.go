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

func init() {
	//从配置文件中获取推送消息的速度、消息推送失败后的重试次数
	pushMessagesSpeed = viper.GetInt("mq.pushMessagesSpeed")
	sendRetryCount = viper.GetInt("mq.sendRetryCount")
	//启动消息推送协程，推送消息到各个消费者客户端
	go pushMessagesToConsumers()
}

//消息推送失败后的重试次数
var sendRetryCount int

//推送消息的速度
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

//推送消息到各个消费者客户端
func pushMessagesToConsumers() {
	for {
		//推送速度控制，延时执行
		time.Sleep(time.Second * time.Duration(pushMessagesSpeed))

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
				if message.Topic == topic && len(message.Topic) > 0 && len(message.MessageData) > 0 {

					//发送消息到消费者客户端
					err := consumer.WriteJSON(message)

					//如果报错
					if err != nil {
						//重试机制
						for i := 0; i < sendRetryCount; i++ {
							//重新发送消息到消费者客户端
							retryErr := consumer.WriteJSON(message)
							//如果发送成功
							if retryErr == nil {
								//将消息标记为已确认状态
								message.Status = 1

								/*
								 * TODO
								 * 对已确认消费的消息进行其他处理
								 */

								//结束循环
								break
							}
							//如果到达重试次数，但仍未发送成功
							if i == sendRetryCount-1 && retryErr != nil {
								//客户端关闭
								consumer.Close()
								//删除map中的客户端
								delete(consumers, key)
							}
						}
					} else {
						//将消息标记为已确认状态
						message.Status = 1

						/*
						 * TODO
						 * 对已确认消费的消息进行其他处理
						 */
					}
				}
				//向控制通道发送信息，表示该协程处理完毕
				controlChan <- 1
			}(key, consumer)
		}

		//待全部协程执行完成后，进入下一轮
		for range consumers {
			//收到协程执行完毕的信息
			<-controlChan
		}
	}
}

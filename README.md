***

# DPMQ.简易消息队列

***

## 基于Golang Gin整合Gorilla WebSocket实现的简易消息队列

`Golang` `Gin` `Gorilla` `WebSocket` `MQ`

***

### 实现功能

* 订阅推送：将单条消息并发推送到多个客户端。

* 流量削峰：将Golang channel的缓冲区当作队列存储大量消息，按固定时间间隔挨个消费缓冲区中的消息。

* 消息推送失败-重试机制。

* 持久化：tan90
***

### 主要部分

##### 生产者消息接收 `producerServer`

* 生产者客户端通过HTTP请求发送消息到消息队列，消息被写入消息通道。

```
//把消息写入消息通道
messageChan <- message
```

##### 消费者消息推送 `consumerServer`

* 消费者客户端与消息队列之间通过WebSocket连接。

* 消息队列从消息通道中读取出消息后，通过for循环结合go协程并发遍历消费者客户端集合，判断该消费者是否与消息同属于一个主题，如果是，则将消息通过WebSocket推送给该客户端。

```
//读取消息通道中的消息
message := <-messageChan

//遍历消费者客户端集合
for key, consumer := range consumers {

	//多协程并发推送消息
	go func(key string,consumer *websocket.Conn) {
	    ...
	    //发送消息到消费者客户端
	    err := consumer.WriteJSON(message)
	    ...
	}
}
```

##### 消息通道 `mqServer`

* 使用golang的通道充当队列，通道的缓冲空间大小决定了消息队列的容量。

```
 //消息通道，用于存放消息(有缓冲)
 var messageChan = make(chan models.Message,viper.GetInt("mq.messageChanBuffer"))
```

##### 控制台 `consoleServer`

* 用于获取消费者客户端列表及消息队列配置信息。

```
//获取全部消费者客户端集合 GetConsumers
GET http://localhost:port/Console/GetConsumers

//获取消息队列详细配置 GetConfig
GET http://localhost:port/Console/GetConfig
```


***

### 连接方式

#### 路由 `routers.go`

```
//生产者接口（http post请求，用于接收生产者客户端发送的消息）
r.POST("/ProducerSend",servers.ProducerSend)

//消费者连接（WebSocket连接，用于推送消息到消费者客户端）
r.GET("/ConsumersConn/:topic/:consumerId", servers.ConsumersConn)
```

#### 访问路径

###### 生产者客户端发送消息到消息队列

* POST `http://localhost:port/ProducerSend`

```
POST请求参数：
messageData   //消息内容  类型：string
topic         //所属主题  类型：string（不能包含符号“|”）
```

```
消息发送成功，返回数据（msg为messageCode）：
{
    "code": 0,
    "msg": "d9a624ffa17a6c51fcb2381686dd335161b7252d"
}

消息发送失败，返回数据（msg为报错信息）：
{
    "code": -1,
    "msg": "topic不能包含字符“|”"
}
```

###### 消息队列通过WebSocket连接推送消息给消费者客户端

* WebSocket `ws://localhost:port/ConsumersConn/{topic}/{consumersId}`

```
WebSocket链接中的参数：
topic        //主题名称（不能包含符号“|”）
consumersId  //消费者客户端Id（不能重复，不能包含符号“|”）
```

```
推送给消费者客户端的消息格式
{
    "MessageCode":"8c01b728ef82ba754a63e61daa43e83c61b744c7",
    "MessageData":"hello",
    "Topic":"test_topic",
    "CreateTime":"2021-12-13 21:04:07",
    "Status":0
}
```

***

### 项目结构

##### config 配置类

* application.yaml `项目配置文件`

* config.go `项目配置文件加载`

##### models 实体类

* model.go `消息模板`

##### routers 路由

* routers.go `路由配置`

##### servers 服务层

* producerServer `生产者消息接收`

* consumerServer `消费者消息推送`

* mqServer `消息通道`

* consoleServer `控制台`

##### utils 工具类

* createCode.go `消息标识码生成`

* localTime.go `获取本地时间`

* md5Sign.go `md5加密`

##### main.go 主函数

***

### 使用说明

* 填写application.yaml内的配置，运行main.go

* application.yaml 配置说明：

```
server:
  # 运行端口号
  port: 80

mq:
  # 消息通道的缓冲空间大小（消息队列的容量）
  messageChanBuffer: 100
  # 推送消息的速度（{pushMessagesSpeed}秒/一条消息）
  pushMessagesSpeed: 0
  # 消息推送失败后的重试次数
  sendRetryCount: 3
```

***
# KapokMQ

## 基于Go整合Gossip+WebSocket的轻量级分布式消息队列

### KapokMQ与Serena应用整合包下载

* https://github.com/dpwgc/kapokmq-server

* https://gitee.com/dpwgc/kapokmq-server

### Golang客户端：kapokmq-go-client

* https://github.com/dpwgc/kapokmq-go-client

* https://gitee.com/dpwgc/kapokmq-go-client

### 控制台前端源码：kapokmq-console

* https://gitee.com/dpwgc/kapokmq-console

### 注册中心源码：Serena

* https://github.com/dpwgc/serena

* https://gitee.com/dpwgc/serena

`Golang` `Gin` `Gorilla` `WebSocket` `MQ` `Gossip`

***

### 软件架构

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/0.jpg)

***

### 实现功能

##### 订阅/发布推送模式：
* 将单条消息通过WebSocket结合go并发推送到多个与消息topic主题相同的消费者客户端。

##### 点对点推送模式：
* 如果有多个消费者客户端连接到消息队列，将消息随机推送给其中一个与消息topic主题相同的客户端。

##### 延时消息发布：
* 可对单条消息设定延时时间，延时推送消息，投送时间精确度受mq.checkSpeed消息检查速度的影响。

##### 负载均衡集群部署：
* 采用Gossip协议连接与同步集群节点，生产者客户端从注册中心获取所有消息队列节点地址并与它们连接，进行负载均衡投递。

##### 消息推送失败重试机制：
* 推送失败后立即重推机制、定期重推未确认消费消息机制(受mq.checkSpeed消息检查速度的影响)。

##### 数据持久化：
* 采用类似于Redis RDB的持久化方案，定期将内存中的消息全量持久化到二进制文件。

##### 过期消息清除：
* 自定义过期时间，定期清除过期消息，默认清除三天前的消息。

##### 网页端控制台：
* 包含查看消息队列配置、生成近一周消息增长折线图、查看各状态消息数量、查看消费者与生产者客户端列表及搜索消息功能。

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/1.jpg)

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/2.jpg)

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/3.jpg)

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/4.jpg)

***

### 配置说明

* config/application.yaml

```yaml
server:
  # ip地址/域名
  addr: 0.0.0.0
  # Gin服务运行端口号
  port: 8011

mq:
  # 生产者、控制台访问密钥（放在请求头部）
  secretKey: test

  # 消息推送模式（1：点对点模式，一条消息只能随机被一个消费者客户端消费。2：订阅/发布推送模式：将消息推送给所有消费者客户端）
  pushType: 1

  # 消息通道的缓冲空间大小（消息队列的容量）
  messageChanBuffer: 10000000

  # 推送消息的速度（{pushMessagesSpeed}秒/一批消息）
  pushMessagesSpeed: 1
  # 单批次推送的消息数量
  pushCount: 1000
  # 消息推送失败后的立即重试的次数
  pushRetryCount: 3

  # 持久化文件
  persistentFile: MQDATA
  # 是否进行持久化（1：是。0：否）
  isPersistent: 1
  # 数据恢复策略（0：清空本地数据，不进行数据恢复。1：将本地数据恢复到内存）
  recoveryStrategy: 1
  # 两次持久化的间隔时间（单位：秒）
  persistentTime: 3

  #是否立即清除已确认消费的消息（1：是。0：否）
  isCleanConsumed: 0

  # 是否开启自动重推未确认消费消息功能（1：是。0：否）
  isRePush: 1
  # 是否开启自动清理过期消息功能（1：是。0：否）
  isClean: 1

  # 检查消息的速度（每隔{checkSpeed}秒检查一批消息，用于消费失败的消息重推、延时消息推送与过期消息清理）
  checkSpeed: 3

  # 消息过期阈值（当消息存在超过{cleanTime}秒后，删除该消息）
  cleanTime: 259200

# Gossip集群配置
cluster:
  # 是否以集群方式部署（1：是。0：否）
  isCluster: 1
  # 该节点的Gossip服务端口号（使用Gossip协议，通过此端口连接注册中心，不能与上面的Gin http服务端口号{server.port}相同）
  gossipPort: 8021
  # 注册中心的ip地址/域名
  registryAddr: 0.0.0.0
  # Serena注册中心的Gossip服务端口号
  registryGossipPort: 8041
```

***

### 打包方式

* 填写application.yaml内的配置。

* 运行项目：

```
（1）GoLand直接运行main.go(调试)
```

```
（2）打包成exe运行(windows部署)

  GoLand终端cd到项目根目录，执行go build命令，生成exe文件
```

```
（3）打包成二进制文件运行(linux部署)

  cmd终端cd到项目根目录，依次执行下列命令：
  SET CGO_ENABLED=0
  SET GOOS=linux
  SET GOARCH=amd64
  go build
  生成二进制执行文件
```

***

### 部署方法

* 在服务器上部署

```
在Windows上部署

/kapokmq                  # 文件根目录
    kapokmq.exe           # 打包后的exe文件
    /config               # 配置目录
        application.yaml  # 配置文件
    /log                  # 日志目录
    /view                 # 前端-Vue项目打包文件
    MQDATA                # 持久化文件
```

```
在Linux上部署

/kapokmq                  # 文件根目录
    kapokmq               # 打包后的二进制文件(程序后台执行:setsid ./KapokMQ)
    /config               # 配置目录
        application.yaml  # 配置文件
    /log                  # 日志目录
    /view                 # 前端-Vue项目打包文件
    MQDATA                # 持久化文件
```

***

### 主要部分

##### 消息通道 `mq.go`

* 使用golang的通道充当队列，通道的缓冲空间大小决定了消息队列的容量。

```
//消息通道，用于存放待消费的消息(有缓冲区)
var messageChan = make(chan models.Message, messageChanBuffer)
```

##### 生产者消息接收 `producer.go`

* 生产者客户端通过WebSocket连接到消息队列，并发送消息到消息队列，消息被写入消息通道。

* 额外提供生产者HTTP接口，可通过HTTP请求向消息队列发送消息。

```
//把消息写入消息通道
messageChan <- message
```

##### 消费者消息推送 `consumer.go`

* 消费者客户端通过WebSocket连接到消息队列。

* 消息队列从消息通道中读取出消息后，通过for循环结合go协程并发遍历消费者客户端集合，判断该消费者是否与消息同属于一个主题，如果是，则将消息通过WebSocket推送给该客户端。

```
//读取消息通道中的消息
message := <-messageChan

...

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

##### 控制台 `console.go`

* 控制台接口：用于获取消费者客户端列表及消息队列配置信息。

```
//检查消息队列服务是否在运行 Ping
POST http://localhost:port/Console/Ping

//获取全部生产者客户端集合 GetProducers
POST http://localhost:port/Console/GetProducers
		
//获取全部消费者客户端集合 GetConsumers
POST http://localhost:port/Console/GetConsumers

//获取消息队列详细配置 GetConfig
POST http://localhost:port/Console/GetConfig

//获取指定状态的消息记录列表 GetMessageList
POST http://localhost:port/Console/GetMessageList

//获取指定状态的简易消息记录列表(不包含消息主体，用于绘制折线图) GetMessageEasy
POST http://localhost:port/Console/GetMessageEasy

//统计各状态消息的数量 CountMessage
POST http://localhost:port/Console/CountMessage
```

* 控制台网页端

```
//启动消息队列后访问：
http://localhost:port/#/Console
```

***

### 客户端连接

#### 路由 `router.go`

```
//生产者连接（WebSocket连接方式，用于接收生产者客户端发送的消息）
GET("/Producers/Conn/:topic/:producerId", server.ProducersConn)

//生产者接口（http post请求方式，用于接收生产者客户端发送的消息）
r.POST("/Producer/Send",servers.ProducerSend)

//消费者连接（WebSocket连接方式，用于推送消息到消费者客户端）
r.GET("/Consumers/Conn/:topic/:consumerId", servers.ConsumersConn)
```

#### 访问路径

##### 生产者客户端发送消息到消息队列

* WebSocket连接方式 `ws://localhost:port/Producers/Conn/{topic}/{producerId}`

```
WebSocket链接中的参数：
topic        //主题名称
ProducerId   //生产者客户端Id
```

```
消息队列接收的消息格式

接收[]byte类型的消息内容(messageData)，转为string类型。
接收int类型的delayTime延时投送时间，为0时表明该消息不是延时消息

//读取websocket中的数据，获取生产者客户端发送的消息内容和延时推送时间
delayTime, data, err := ws.ReadMessage()
messageData := string(data)

再将消息封装成model.Message类型，插入消息通道。
```

* POST请求方式 `http://localhost:port/Producer/Send`

```
POST请求参数：
Header:
secretKey     //访问密钥

Body:
messageData   //消息内容  类型：json string
topic         //所属主题  类型：string
delayTime     //延迟推送时间（单位：秒）
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
    "msg": "Required data cannot be empty"
}
```

##### 消息队列通过WebSocket连接推送消息给消费者客户端

* WebSocket `ws://localhost:port/Consumers/Conn/{topic}/{consumerId}`

```
WebSocket链接中的参数：
topic        //主题名称
consumerId  //消费者客户端Id
```

```
//通过WriteJSON()函数将model.Message类型的消息转为Json字符串发送
推送给消费者客户端的消息格式
{
    "MessageCode":"8c01b728ef82ba754a63e61daa43e83c61b744c7",
    "MessageData":"hello",
    "Topic":"test_topic",
    "CreateTime":"1640975470",
    "ConsumedTime":"1640975520",
    "Status":-1
}
```

* 生产者、消费者客户端与消息队列进行WebSocket连接后，需输入密钥登录
```
ws://127.0.0.1:8011/Consumers/Conn/test_topic/1
消费者客户端与消息队列建立连接

服务端回应 2022-01-02 15:14:53
"Please enter the secret key"   //提示输入密钥

你发送的信息 2022-01-02 15:15:06
qqq                             //输入错误的密钥

服务端回应 2022-01-02 15:15:06
"Secret key matching error"     //提示密钥出错

服务端回应 2022-01-02 15:15:06
"Please enter the secret key"   //再次提示输入密钥

你发送的信息 2022-01-02 15:15:13
dpmq                            //输出正确的密钥

服务端回应 2022-01-02 15:15:13
"Secret key matching succeeded" //密钥验证成功

服务端回应 2022-01-02 15:15:13   //消费者客户端可以开始接收消息
...
...
...
```

```
ws://127.0.0.1:8011/Producers/Conn/test_topic/1
生产者客户端与消息队列建立连接

服务端回应 2022-01-02 15:14:53
"Please enter the secret key"   //提示输入密钥

你发送的信息 2022-01-02 15:15:06
qqq                             //输入错误的密钥

服务端回应 2022-01-02 15:15:06
"Secret key matching error"     //提示密钥出错

服务端回应 2022-01-02 15:15:06
"Please enter the secret key"   //再次提示输入密钥

你发送的信息 2022-01-02 15:15:13
dpmq                            //输出正确的密钥

服务端回应 2022-01-02 15:15:13
"Secret key matching succeeded" //密钥验证成功

服务端回应 2022-01-02 15:15:13   //生产者客户端可以开始发送消息
...
...
...
```

***

### 项目结构

##### cluster 集群相关

* join.go `加入指定集群`

##### config 配置类

* application.yaml `项目配置文件`

* config.go `项目配置文件加载`

##### middleware 中间件

* cors.go `跨域配置`

* safe.go `安全验证`

##### model 实体类

* model.go `数据模板`

##### persistent 持久化

* fileRW.go `文件读写`

* persData.go `持久化到硬盘`

* recovery.go `数据恢复`

##### router 路由

* router.go `路由配置`

##### server 服务层

* console,go `控制台接口`

* producer.go `生产者消息接收`

* consumer.go `消费者消息推送`

* mq.go `消息队列`

* log.go `日志记录`

* check.go `消息检查-消息重推与过期消息清理`

##### utils 工具类

* createCode.go `消息标识码生成`

* localTime.go `获取本地时间`

* md5Sign.go `md5加密`

* toTimestamp.go `日期字符串转时间戳`

##### view 前端Vue项目打包文件

* css

* js

* index.html

##### main.go 主函数
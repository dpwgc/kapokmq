# KapokMQ

## 基于Go整合Gossip+WebSocket的轻量级分布式消息队列

![Go](https://img.shields.io/static/v1?label=LICENSE&message=Apache-2.0&color=orange)
![Go](https://img.shields.io/static/v1?label=Go&message=v1.17&color=blue)
[![github](https://img.shields.io/static/v1?label=Github&message=kapokmq&color=blue)](https://github.com/dpwgc/kapokmq)
[![star](https://gitee.com/dpwgc/kapokmq/badge/star.svg?theme=dark)](https://gitee.com/dpwgc/kapokmq/stargazers)
[![fork](https://gitee.com/dpwgc/kapokmq/badge/fork.svg?theme=dark)](https://gitee.com/dpwgc/kapokmq/members)

***

#### KapokMQ是一个使用Go语言实现的，单机TPS超10万的轻量级分布式消息队列

#### 特点：

* 部署KapokMQ无需安装任何外部依赖应用/环境，单文件应用，且配置简单，启动方便。

* 支持多种消息推送方式（订阅/发布推送、点对点推送、延时消息发布）。

* 消息队列与客户端采用WebSocket连接，全双工通讯，消息与ACK发送都在一条连接上，简易高效。

* 集群基于Gossip构建，可自动探测集群节点，能轻易做到水平扩展、断线重连。

* 采用KV型内存数据存储方式存储消息列表，哈希表插入与更新数据速度快。

* 提供全量数据持久化方案及WAL预写日志记录，用户可在高性能与高可靠之间自由选择。

* 内置网页端控制台，可直接通过网页查看消息队列的运行情况。

***

#### KapokMQ与Serena应用整合包下载与安装
* https://github.com/dpwgc/kapokmq-server `github`
* https://gitee.com/dpwgc/kapokmq-server `gitee`

#### Golang客户端 ~ kapokmq-go-client
* https://github.com/dpwgc/kapokmq-go-client `github`
* https://gitee.com/dpwgc/kapokmq-go-client `gitee`

#### 注册中心源码 ~ Serena
* https://github.com/dpwgc/serena `github`
* https://gitee.com/dpwgc/serena `gitee`

#### 控制台前端源码 ~ kapokmq-console
* https://gitee.com/dpwgc/kapokmq-console `gitee`

***

### 软件架构
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/deploy.jpg)
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/inner.jpg)

***

### 实现功能

#### 订阅/发布推送模式：
* 将单条消息通过WebSocket结合协程并发推送到多个与消息topic主题相同的消费者客户端。

#### 点对点推送模式：
* 如果有多个消费者客户端连接到消息队列，将消息随机推送给其中一个与消息topic主题相同的客户端。

#### 延时消息发布：
* 可对单条消息设定延时时间，秒级延时推送消息，投送时间精确度受mq.checkSpeed消息检查速度的影响。

#### ACK消息确认机制：
* 消息队列接收到消息后，将向生产者发送确认接收ACK，可确保消息不在消息被持久化之前丢失。
* 消费者接收到消息后，将向消息队列发送确认消费ACK，可确保消息不在消息消费环节丢失。

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/ack2.jpg)

#### 负载均衡集群部署：
* 采用Gossip协议连接与同步集群节点，生产者客户端从注册中心获取所有消息队列节点地址并与它们连接，进行负载均衡投递（将消息随机投送到其中一个消息队列节点）。可做到不停机水平扩展。

#### 消息推送失败自动重试机制：
* 定期重推未确认消费且超时的消息(受mq.checkSpeed消息检查速度的影响)。

#### KV型内存数据存储：
* 采用sync.Map存储所有消息，控制台访问、消息检查、ACK确认接收/消费、全量数据持久化等一系列读写操作都在sync.Map上进行。

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/kv.jpg)

#### 数据持久化：
* 方式一：周期性全量数据持久化。
* 方式二：周期性全量数据持久化结合WAL预写日志（影响性能）。
* （注：周期性全量数据持久化不会写入已经被消费的消息）

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/pers.jpg)

#### 过期消息清除：
* 自定义过期时间，定期清除过期消息，默认清除两天前的消息。

#### 网页端控制台：
* 包含查看消息队列配置、生成近一周消息增长折线图、查看各状态消息数量、查看消费者与生产者客户端列表、搜索消息及查看集群节点列表功能。

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/page1.jpg)
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/page2.jpg)
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/page2-5.jpg)
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/page3.jpg)
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/page4.jpg)
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/page4-5.jpg)
![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/page4.jpg)

***

### 单机性能测试

#### 测试程序和消息队列都运行在本机。

#### 配置：

* 2019年的轻薄本，未充电状态，平衡模式

* 4核8线程处理器：Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz 1.80 GHz

* LPDDR3内存：8.00 GB (7.85 GB 可用)

* 固态硬盘：WDC PC SN720 SDAPNTW-512G

#### 测试程序（部分）

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/test_code.jpg)

#### 选择数据持久化方式一，未开启WAL预写日志的情况下，模拟三十万并发消息插入，最高TPS可达十万。

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/nowal_test.jpg)

#### 选择数据持久化方式二，开启WAL预写日志，同样模拟三十万并发消息插入，受磁盘IO影响，消息插入速度存在大幅波动与衰减。

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/wal_test.jpg)

***

### 配置文件

* ./application.yaml

***

### 打包方法

##### 安装并配置go环境

##### 应用打包：

* 使用build.bat一键打包

```
Windows系统下直接运行./build.bat文件
自动执行打包命令
生成Linux二进制文件与Windows exe文件
```

* 打包成windows exe

```
  GoLand终端cd到项目根目录，执行go build命令，生成exe文件
```

* 打包成linux二进制文件

```
  cmd终端cd到项目根目录，依次执行下列命令：
  SET CGO_ENABLED=0
  SET GOOS=linux
  SET GOARCH=amd64
  go build
  生成二进制执行文件
```

***

### 部署方式

* 在服务器上部署

```
在Windows上部署

/kapokmq                  # 文件根目录
    kapokmq.exe           # 打包后的exe文件
    application.yaml      # 配置文件
    /log                  # 日志目录
    /view                 # 前端-Vue项目打包文件
    MQDATA                # 持久化文件
    WAL.log               # 预写日志  
```

```
在Linux上部署

/kapokmq                  # 文件根目录
    kapokmq               # 打包后的二进制文件(后台运行指令:setsid ./kapokmq)
    application.yaml      # 配置文件
    /log                  # 日志目录
    /view                 # 前端-Vue项目打包文件
    MQDATA                # 持久化文件
    WAL.log               # 预写日志  
```

***

### 消息格式说明

```yaml
MessageCode  string   #消息唯一标识码
MessageData  string   #消息内容（一般为JSON格式的字符串）
Topic        string   #消息所属主题
CreateTime   int64    #消息创建时间（秒级时间戳）
ConsumedTime int64    #消息被消费时间（秒级时间戳）
DelayTime    int64    #延迟推送时间（单位：秒）
Status       int      #消息状态（-1：待消费。0：未到推送时间的延时消息。1：已消费）
```

***

### 客户端连接

#### 生产者客户端连接到消息队列

* WebSocket `ws://localhost:port/Producers/Conn/{topic}/{producerId}`

```
WebSocket链接中的参数：
topic        //主题名称
ProducerId   //生产者客户端Id
```

* 消息队列接收的消息格式

* 生产者客户端推送给消息队列的Json字符串消息格式

```json
{
    "MessageData":"hello",
    "DelayTime":0
}
```

* 消息队列接收到该消息后（写入日志后），通过该websocket连接向生产者客户端发送ACK，ACK内容为字符串"ok"

* 如果生产者客户端选择异步发送消息方式，则可忽略该ACK。

* 如果要追求消息的可靠性，可以利用该ACK机制发送同步消息，即生产者在发送完一条消息后，必须收到消息队列发来的ACK才能继续发送下一条消息。

* 消息队列发送给生产者的ACK字符串样式

```json
"ok"
```

#### 消费者客户端连接到消息队列

* WebSocket `ws://localhost:port/Consumers/Conn/{topic}/{consumerId}`

```
WebSocket链接中的参数：
topic        //主题名称
consumerId   //消费者客户端Id
```

* 通过WriteJSON()函数将model.Message类型的消息转为Json字符串发送

* 消息队列推送给消费者客户端的Json字符串消息格式

```json
{
    "MessageCode":"8c01b728ef82ba754a63e61daa43e83c61b744c7",
    "MessageData":"hello",
    "Topic":"test_topic",
    "CreateTime":1640975470,
    "ConsumedTime":1640975520,
    "DelayTime":0,
    "Status":-1
}
```

* 消费者客户端接收到该消息后（写入日志后），通过该websocket连接向消息队列发送ACK，ACK内容为消息的唯一标识码MessageCode

* 消费者客户端发送给消息队列的ACK字符串样式

```json
"8c01b728ef82ba754a63e61daa43e83c61b744c7"
```

***

### 客户端连接流程演示

生产者、消费者客户端与消息队列建立连接后，需输入密钥登录

#### 生产者客户端与消息队列建立连接

```
ws://127.0.0.1:8011/Producers/Conn/test_topic/1

服务端回应 2022-01-02 15:14:53
"Please enter the secret key"   //提示输入密钥

客户端发送 2022-01-02 15:15:06
"qqq"                           //输入错误的密钥

服务端回应 2022-01-02 15:15:06
"Secret key matching error"     //提示密钥出错

服务端回应 2022-01-02 15:15:06
"Please enter the secret key"   //再次提示输入密钥

客户端发送 2022-01-02 15:15:13
"test"                          //输出正确的密钥

服务端回应 2022-01-02 15:15:13
"Secret key matching succeeded" //提示密钥验证成功

客户端发送 2022-01-02 15:15:15   //生产者客户端可以向消息队列发送消息
"{.. Json SendMessage ..}"

服务端回应 2022-01-02 15:15:15   //消息队列接收到消息后，向生产者发送ACK
"ok"                           //ACK内容为字符串"ok"
```

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/ws_test1.jpg)

#### 消费者客户端与消息队列建立连接

```
ws://127.0.0.1:8011/Consumers/Conn/test_topic/1

服务端回应 2022-01-02 15:14:53
"Please enter the secret key"   //提示输入密钥

客户端发送 2022-01-02 15:15:06
"qqq"                           //输入错误的密钥

服务端回应 2022-01-02 15:15:06
"Secret key matching error"     //提示密钥出错

服务端回应 2022-01-02 15:15:06
"Please enter the secret key"   //再次提示输入密钥

客户端发送 2022-01-02 15:15:13
"test"                          //输出正确的密钥

服务端回应 2022-01-02 15:15:13
"Secret key matching succeeded" //提示密钥验证成功

服务端回应 2022-01-02 15:15:13   //消息队列可以向消费者客户端发送消息
"{.. Json Message ..}"
"{.. Json Message ..}"

客户端发送 2022-01-02 15:15:14   //消费者接收到消息后，向消息队列发送ACK
"8c01b728ef82ba754a63e61daa43e83c61b744c7"  //ACK内容为MessageCode
"sdiw2b7quh82basdsa17sdqdqw81d83c61bqdhhu"  //可异步发送确认消费ACK
```

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/ws_test2.jpg)

***

### 主要模块

##### 消息通道与消息列表 `mq.go`

* 使用golang的通道chan充当队列，通道的缓冲空间大小决定了消息队列的容量。

* 使用sync.Map存储所有消息，用于数据持久化、消息检查、控制台数据获取。

```
//消息通道，用于存放待消费的消息(有缓冲区)
var messageChan = make(chan models.Message, messageChanBuffer)

// MessageList 消息列表，存放所有消息记录
var MessageList sync.Map
```

##### 生产者消息接收 `server/producer.go`

* 生产者客户端通过WebSocket连接到消息队列（github.com/gorilla/websocket），并发送消息到消息队列，消息被写入消息通道。

* 额外提供生产者HTTP接口，可通过HTTP请求向消息队列发送消息。

##### 消费者消息推送 `server/consumer.go`

* 消费者客户端通过WebSocket连接到消息队列（github.com/gorilla/websocket）。

* 包含订阅/发布、点对点两种推送模式。

* 消费者客户端接收消息后，将向消息队列发送一条ACK确认字符（内容为消息标识码messageCode），消息队列再根据此ACK将指定messageCode的消息更改为已消费状态。

##### 数据持久化 `persistent`

* 全量数据写入：定期将MessageList消息列表中的未消费消息和延时消息转换为[]byte类型数据，并写入二进制文件，类似于Redis RDB持久化方式。

* 全量数据写入结合WAL日志：定期将内存中的消息全量持久化到二进制文件（不包括已被消费的消息），在两次全量数据持久化之间，每次接收或更新消息操作都将写入WAL日志，最大程度避免消息丢失。

* 数据恢复：从二进制文件及WAL日志中读取数据，并将数据恢复至MessageList消息列表中，重新推送未消费的消息。

![avatar](https://dpwgc-1302119999.cos.ap-guangzhou.myqcloud.com/kapokmq/recovery.jpg)

##### 消息检查 `server/check.go`

* 每隔一段时间遍历一次MessageList消息列表，检查其中是否有消费失败、延时消费、超时未消费、已过期的消息。可重新推送消息，或清除过期的消息。

##### 加入Gossip集群 `cluster/join.go`

* 使用 github.com/hashicorp/memberlist 构建并链接Gossip集群服务。

* 借助Gossip协议扩散同步的特性，可以随时向集群中添加新的消息队列节点。

##### 控制台接口 `console/api.go`

* 控制台接口：用于获取生产者/消费者客户端列表、消息队列配置信息及集群内消息队列节点列表，对消息进行查询或删除操作。

* 控制台网页端，启动消息队列后访问：

```
http://localhost:port/#/Console
```

***

### 项目结构

##### cluster 集群相关

* join.go `加入指定集群`

##### config 配置类

* config.go `项目配置文件加载`

##### console 控制台

* api.go `控制台接口`

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

##### application.yaml 配置文件

##### main.go 主函数

***

### 后期计划

|实现功能|功能说明|当前进度|
|---|---|---|
|Java客户端|Maven包，websocket连接，Demo：https://gitee.com/dpwgc/kapokmq-java-client|未完成|
|拉模式消费|消费者主动拉取消息队列的消息|计划中|
|主备消息队列|为每个mq主节点绑定一个备用节点，宕机时立即切换到备用节点|计划中|
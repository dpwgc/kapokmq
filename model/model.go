package model

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

// Consumer 消费者客户端模板
type Consumer struct {
	ConsumerId string //消费者Id
	Topic      string //消费者所属主题
	ConsumerIp string //消费者ip地址
	JoinTime   int64  //消费者加入时间
}

// Producer 生产者客户端模板
type Producer struct {
	ProducerId string //生产者Id
	Topic      string //生产者所属主题
	ProducerIp string //生产者ip地址
	JoinTime   int64  //生产者加入时间
}

// Node 消息队列节点结构体
type Node struct {
	Name string
	Addr string
	Port string
}

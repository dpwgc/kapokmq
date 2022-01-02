package model

// Message 消息模板
type Message struct {
	MessageCode  string      //消息标识码
	MessageData  interface{} //消息内容（一般为JSON格式的字符串）
	Topic        string      //消息所属主题
	CreateTime   int64       //消息创建时间
	ConsumedTime int64       //消息被消费时间
	Status       int         //消息状态（-1：刚进入无状态，0：未被消费，1：已被消费）
}

// Consumer 消费者客户端模板
type Consumer struct {
	ConsumerId string //消费者Id
	Topic      string //消费者所属主题
	ConsumerIp string //消费者ip地址
	JoinTime   int64  //消费者加入时间
}

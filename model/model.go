package model

//消息模板
type Message struct {
	MessageCode string      //消息标识码
	MessageData interface{} //消息内容（一般为JSON格式的字符串）
	Topic       string      //消息所属主题
	CreateTime  string      //消息创建时间
	Status      int         //消息状态（0：未被消费，1：已被消费）
}

//客户端模板
type Consumer struct {
	ConsumerId string //客户端Id
	Topic      string //客户端所属主题
}

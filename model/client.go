package model

import "kapokmq/utils"

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

// NewConsumer 生成消费者客户端模板
func NewConsumer(topic string, consumerId string, consumerIp string) *Consumer {

	//生成消费者客户端模板
	consumer := Consumer{}
	consumer.ConsumerId = consumerId
	consumer.Topic = topic
	consumer.ConsumerIp = consumerIp
	consumer.JoinTime = utils.GetLocalDateTimestamp()

	return &consumer
}

// NewProducer 生成生产者客户端模板
func NewProducer(topic string, producerId string, producerIp string) *Producer {

	//生成生产者客户端模板
	producer := Producer{}
	producer.ProducerId = producerId
	producer.Topic = topic
	producer.ProducerIp = producerIp
	producer.JoinTime = utils.GetLocalDateTimestamp()

	return &producer
}

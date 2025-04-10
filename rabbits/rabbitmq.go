package rabbitmq

import (
	"errors"
	"log"
	"strings"

	"github.com/streadway/amqp"
)

// 定义全局mqmap
var RabbitMqMap = make(map[string]*rabbitMQ)

// RabbitMQ 用于管理和维护rabbitmq的对象
type rabbitMQ struct {
	//wg           sync.WaitGroup
	channel *amqp.Channel
	mqConn  *amqp.Connection
}

//连接mq，config中可配置channel连接的容量和心跳时长
//默认为:
/**
maxChannelMax = (2 << 15) - 1
defaultHeartbeat         = 10 * time.Second
defaultConnectionTimeout = 30 * time.Second
defaultProduct           = "https://github.com/streadway/amqp"
defaultVersion           = "β"
// Safer default that makes channel leaks a lot easier to spot
// before they create operational headaches. See https://github.com/rabbitmq/rabbitmq-server/issues/1593.
defaultChannelMax = (2 << 10) - 1
*/
func (mq *rabbitMQ) connToMq(url string, config *amqp.Config) (rabbitMq *rabbitMQ, err error) {
	mq.mqConn, err = amqp.DialConfig(url, *config)
	if err != nil {
		return
	}
	mq.channel, err = mq.mqConn.Channel()
	mq.mqConn.Channel()
	if err != nil {
		return
	}
	return mq, nil
}

// 直接初始化队列
func (mq *rabbitMQ) PrepareQueue(queueName string) (queue amqp.Queue, err error) {
	if queueName == "" {
		return queue, errors.New("queueName为空")
	}
	queue, err = mq.channel.QueueDeclare(
		queueName, //name
		true,      //durable，是否持久化，默认持久需要根据情况选择
		false,     //delete when unused
		false,     //exclusive
		false,     //no-wait
		nil,       //arguments
	)
	return
}

// prepareExchange 准备rabbitmq的Exchange
func (mq *rabbitMQ) PrepareExchange(exchangeName, exchangeType string) error {
	if exchangeName == "" {
		return errors.New("exchangeName为空")
	}
	err := mq.channel.ExchangeDeclare(
		exchangeName, // exchange
		exchangeType, // type
		true,         // durable 是否持久化，默认持久需要根据情况选择
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	)

	if nil != err {
		return err
	}

	return nil
}

// 通过exchange发送消息
func (mq *rabbitMQ) ExchangeSend(exchangeName, routingKey string, publishing amqp.Publishing) error {

	return mq.channel.Publish(
		exchangeName, //exchangeName
		routingKey,   //routing key
		true,         //mandatory
		false,        //immediate
		publishing,
	)
}

// 通过队列发送消息
func (mq *rabbitMQ) QueueSend(queueName string, publishing amqp.Publishing) error {

	return mq.channel.Publish(
		"",        //exchangeName
		queueName, //queue name
		false,     //mandatory
		false,     //immediate
		publishing,
	)

}

// 消费队列,内部方法会阻塞,使用时需要单独启用一个线程处理，常驻后台执行
func (mq *rabbitMQ) QueueConsume(queueName, consumer string) (delivery <-chan amqp.Delivery, err error) {
	err = mq.channel.Qos(1, 0, true)
	if err != nil {
		log.Fatal("Queue Consume: ", err.Error())
		return nil, err
	}
	//后期可调整参数
	delivery, err = mq.channel.Consume(
		queueName, // queue
		consumer,  // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatal("Queue Consume: ", err.Error())
		return nil, err
	}
	return delivery, nil
}

// 队列绑定exchange
func (mq *rabbitMQ) QueueBindExchange(queueName, routingKey, exchangeName string) error {
	return mq.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
}

// 关闭连接
func Close() {
	for k := range RabbitMqMap {
		RabbitMqMap[k].channel.Close()
		RabbitMqMap[k].mqConn.Close()
	}
}

func InitRabbitMq() {
	var (
		mq     rabbitMQ
		config amqp.Config
		err    error
	)
	//此处可定义多个配置，可调整
	dsn := "amqp://guest:guest@localhost:5672/"
	config.Vhost = "/"
	RabbitMqMap["mq"], err = mq.connToMq(dsn, &config)
	if err != nil {
		log.Fatal("[rabbit-mq] connect to rabbit-mq error:" + err.Error())
	} else {
		log.Println("[rabbit-mq] connect success")
	}

}

// 获取连接
func GetRabbitConn(name ...string) *rabbitMQ {
	rabbitName := ""
	if len(name) > 0 {
		rabbitName = strings.ToLower(name[0])
	}
	if rabbitName == "" {
		return RabbitMqMap["mq"]
	}
	return RabbitMqMap[rabbitName]
}

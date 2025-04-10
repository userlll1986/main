package main

import (
	"log"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://lafba13j4134:llhafaif99973@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %s", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"main_exchange", // exchange name
		"direct",        // type
		true,            // durable
		false,           // auto-delete
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %s", err)
	}

	for i := 0; i < 10; i++ {
		// 消息Routing Key
		routingKey := "wodekey.log.info"
		if i%2 == 0 {
			routingKey = "wodekey.log.debug"
		} else if i%3 == 0 {
			routingKey = "wodekey.log.error"
		}

		err = ch.Publish(
			"main_exchange", // exchange
			routingKey,      // routing key
			true,            // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte("这是我的消息:" + strconv.Itoa(i)),
			})
		if err != nil {
			log.Fatalf("Failed to publish message: %s", err)
		}
		log.Printf("Sent message to %s: %s", routingKey, "这是我的消息:"+strconv.Itoa(i))
		time.Sleep(1 * time.Second)
	}
}

package main

import (
	"log"

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

	queues := []string{"test-qq", "test-qq2", "test-qq3"}
	bindings := map[string]string{
		"test-qq":  "wodekey.log.info",
		"test-qq2": "wodekey.log.debug",
		"test-qq3": "wodekey.log.error",
	}

	for _, queueName := range queues {
		_, err = ch.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare queue: %s", err)
		}

		err = ch.QueueBind(
			queueName,           // queue name
			bindings[queueName], // routing key
			"main_exchange",     // exchange
			false,               // no-wait
			nil,                 // arguments
		)
		if err != nil {
			log.Fatalf("Failed to bind queue %s: %s", queueName, err)
		}
	}

	msgs := make(chan amqp.Delivery, 1)

	for _, queueName := range queues {
		// 消费消息
		deliveries, err := ch.Consume(
			queueName,           // queue
			bindings[queueName], // consumer
			true,                // auto-ack
			false,               // exclusive
			false,               // no-local
			false,               // no-wait
			nil,                 // args
		)
		if err != nil {
			log.Fatalf("Failed to register a consumer: %s", err)
		}

		go func() {
			for d := range deliveries {
				msgs <- d
			}
		}()
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}

package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func RabbitConnect() (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. Error: %s", err)
	}

	return ch, conn
}
func CreateConsumer(ch *amqp.Channel) <-chan amqp.Delivery {
	messages, err := ch.Consume( //Поставщик
		InputQueueName, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Can't reg consumer: %s", err)
	}
	return messages
}
func RabCtxDefine() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return ctx
}
func SendMessage(ch *amqp.Channel, ctx context.Context, body string) {
	err := ch.PublishWithContext(ctx,
		"",     // exchange
		"Done", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("failed to publish a message. Error: %s", err)
	}
	fmt.Println("Done!")
}

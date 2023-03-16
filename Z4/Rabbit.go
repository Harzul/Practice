package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func (r *Rabbit) RabbitConnect() {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}
	(*r).conn = conn

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. Error: %s", err)
	}
	(*r).ch = ch

}

func (r *Rabbit) CreateConsumer() <-chan amqp.Delivery {
	messages, err := r.ch.Consume( //Поставщик
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
func (r *Rabbit) RabCtxDefine() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	(*r).ctx = ctx
}
func (r *Rabbit) SendMessage(data []byte) {
	err := (*r).ch.PublishWithContext((*r).ctx,
		"",     // exchange
		"Done", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	if err != nil {
		log.Fatalf("failed to publish a message. Error: %s", err)
	}
	fmt.Println("Done!")
}

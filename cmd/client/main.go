package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"

	"github.com/streadway/amqp"

	"bloxroute/internal/message"
)

var (
	action      = flag.String("action", "", "action, possible values: add, remove, get, list")
	key         = flag.String("key", "", "key parameter used for add, remove, get actions")
	value       = flag.String("value", "", "value parameter used for add action")
	rabbitmqURL = flag.String("rabbitmq_url", "amqp://guest:guest@localhost:5672/", "RabbitMQ connection URL")
)

func init() {
	flag.Parse()
}

func main() {
	conn, err := amqp.Dial(*rabbitmqURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	messageType := message.Type(*action)
	switch t := messageType; t {
	case message.AddItemType:
		if *key == "" || *value == "" {
			println("this action requires both key and value parameter")
			return
		}

		log.Println("adding item")
		if err := encoder.Encode(&message.AddItem{Key: *key, Value: *value}); err != nil {
			log.Fatal(err)
		}
	case message.GetItemType:
		if *key == "" {
			println("this action requires key parameter")
			return
		}

		log.Println("getting item")
		if err := encoder.Encode(&message.GetItem{Key: *key}); err != nil {
			log.Fatal(err)
		}
	case message.RemoveItemType:
		if *key == "" {
			println("this action requires key parameter")
			return
		}

		log.Println("remove item")
		if err := encoder.Encode(&message.RemoveItem{Key: *key}); err != nil {
			log.Fatal(err)
		}
	case message.GetAllItemsType:
		log.Println("list items")

		// Let's assume that there could be some filters
		if err := encoder.Encode(&message.GetAllItems{}); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("wrong action")
	}

	var buf2 bytes.Buffer
	if err := gob.NewEncoder(&buf2).Encode(message.Message{Type: messageType, Body: buf.Bytes()}); err != nil {
		panic(err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        buf2.Bytes(),
		})
	if err != nil {
		log.Fatal(err)
	}
}

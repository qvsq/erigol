package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/streadway/amqp"

	"bloxroute/internal/message"
	"bloxroute/internal/orderedmap"
)

type item struct {
	Key   string
	Value string
}

var (
	output       = flag.String("output", "output.txt", "output file")
	debug        = flag.Bool("debug", false, "whether to show debug information")
	rabbitmqURL  = flag.String("rabbitmq_url", "amqp://guest:guest@localhost:5672/", "RabbitMQ connection URL")
	workersCount = flag.Int("workers", 4, "workers count")

	// OutputLogger outputs the result of commands
	OutputLogger *log.Logger
	// DebugLogger outputs debug information
	DebugLogger *log.Logger
	// ErrorLogger outputs errors
	ErrorLogger *log.Logger
)

func init() {
	flag.Parse()

	file, err := os.OpenFile(*output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	OutputLogger = log.New(file, "", 0)
	var debugOutput io.Writer
	if *debug {
		debugOutput = os.Stdout
	} else {
		debugOutput = ioutil.Discard
	}

	DebugLogger = log.New(debugOutput, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	dataMap := orderedmap.New()

	for i := 0; i < *workersCount; i++ {
		go func() {
			// It is something bad happens when using one buffer for different messages
			// Probably it is because buffer reset retains underlying storage
			var msgBuffer bytes.Buffer
			var bodyBuffer bytes.Buffer
			msgDecoder := gob.NewDecoder(&msgBuffer)
			bodyDecoder := gob.NewDecoder(&bodyBuffer)

			for msg := range msgs {
				msgBuffer.Reset()
				bodyBuffer.Reset()

				if _, err := msgBuffer.Write(msg.Body); err != nil {
					ErrorLogger.Print(err)
					continue
				}

				var m message.Message
				if err := msgDecoder.Decode(&m); err != nil {
					ErrorLogger.Print(err)
					continue
				}
				if _, err := bodyBuffer.Write(m.Body); err != nil {
					ErrorLogger.Print(err)
					continue
				}

				switch m.Type {
				case message.AddItemType:
					var addItem message.AddItem
					if err := bodyDecoder.Decode(&addItem); err != nil {
						ErrorLogger.Print(err)
						continue
					}

					DebugLogger.Printf("Adding item, key: %s value: %s", addItem.Key, addItem.Value)

					// We are storing entire item to be able easily write key value pair into file
					item := item{Key: addItem.Key, Value: addItem.Value}
					dataMap.Store(addItem.Key, item)
				case message.GetItemType:
					var getItem message.GetItem
					if err := bodyDecoder.Decode(&getItem); err != nil {
						ErrorLogger.Print(err)
						continue
					}

					value, ok := dataMap.Load(getItem.Key)
					if ok {
						item := value.(item)
						OutputLogger.Printf("%+v", item)
						DebugLogger.Printf("Getting item, key: %s value: %s", item.Key, item.Value)
					} else {
						DebugLogger.Printf("Item with key: %s not found", getItem.Key)
					}
				case message.RemoveItemType:
					var removeItem message.RemoveItem
					if err := bodyDecoder.Decode(&removeItem); err != nil {
						ErrorLogger.Print(err)
						continue
					}
					DebugLogger.Printf("Deleting item with key: %s", removeItem.Key)

					dataMap.Delete(removeItem.Key)
				case message.GetAllItemsType:
					// No reason to decode GetAllItems

					for el := dataMap.Front(); el != nil; el = el.Next() {
						item := el.Value.(item)

						OutputLogger.Printf("%+v", item)
						DebugLogger.Printf("Getting item, key: %s value: %s", item.Key, item.Value)
					}
				}
			}
		}()
	}

	<-forever
}

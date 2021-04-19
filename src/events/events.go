package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gnoirzox/faceit-users/utils"

	"github.com/streadway/amqp"
)

func OpenRabbitMQConnection() (*amqp.Connection, error) {
	var (
		host     = utils.GetEnv("RABBITMQ_HOST", "localhost")
		port     = utils.GetEnv("RABBITMQ_PORT", "5672")
		user     = utils.GetEnv("RABBITMQ_USER", "guest")
		password = utils.GetEnv("RABBITMQ_PASS", "guest")
	)

	rabbitmqInfo := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(rabbitmqInfo)

	if err != nil {
		log.Println("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	return conn, err
}

func jsonSerialize(msg map[string]string) ([]byte, error) {
	var b bytes.Buffer

	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)

	return b.Bytes(), err
}

func PublishMessage(queueName string, message map[string]string) error {
	conn, err := OpenRabbitMQConnection()

	defer conn.Close()

	if err != nil {
		log.Println(err.Error())

		return err
	}

	channel, err := conn.Channel()

	defer channel.Close()

	if err != nil {
		log.Println("%s: %s", "Failed to open a channel", err)

		return err
	}

	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Println("%s: %s", "Failed to define a queue", err)

		return err
	}

	jsonMessage, err := jsonSerialize(message)

	if err != nil {
		log.Println(err.Error())

		return err
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonMessage),
		})

	if err != nil {
		log.Println("%s: %s", "Failed to publish a message", err)

		return err
	}

	return nil
}

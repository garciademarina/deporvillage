package broker

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Broker ...
type Broker struct {
	con       *amqp.Connection
	ch        *amqp.Channel
	queueName string
}

// NewBroker creates an adding service with the necessary dependencies
func NewBroker(url, queueName string) (*Broker, error) {

	fmt.Printf("%s\n", url)
	// url := "amqp://guest:guest@localhost"
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName, //name
		true,      //durable
		false,     //delete when unused
		false,     //exclusive
		false,     //no-wait
		nil,       //arguements
	)
	if err != nil {
		return nil, err
	}

	return &Broker{conn, ch, queueName}, nil
}

// Publish ...
func (b *Broker) Publish(body string) error {
	err := b.ch.Publish(
		"",          //exchange
		b.queueName, //routing key
		false,       //mandatory
		false,       //immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// Close ...
func (b *Broker) Close() error {
	err := b.con.Close()
	if err != nil {
		return err
	}
	err = b.ch.Close()
	if err != nil {
		return err
	}
	return nil
}

package rabbitmq

import (
	"bitbucket.org/smaug-hosting/services/queue/message"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

type QueueArgs struct {
	url string
}

type Queue struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	consumerId string
}

func (q *Queue) Setup(args map[string]interface{}) {
	options := args
	if options["url"] == "" {
		logrus.Fatalf("Could not connect to RabbitMQ instance: please provide a URL")
	}

	conn, err := amqp.Dial(options["url"].(string))
	if err != nil {
		logrus.Fatalf("Could not connect to RabbitMQ instance: %s", err)
	}

	q.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Could not establish RabbitMQ channel: %s", err)
	}

	q.ch = ch

	q.consumerId = uuid.New().String()
	logrus.Infof("Started consumer with id %s", q.consumerId)

	SetupTopology(*ch)

}

func (q *Queue) Subscribe(queueName string) <-chan message.Message {
	msgChan := make(chan message.Message)

	go func() {

		messageChan, err := q.ch.Consume(queueName, q.consumerId, false, false, false, false, nil)
		if err != nil {
			logrus.Errorf("Caught error while consuming from rabbitmq: %s", err)
		} else {
			for m := range messageChan {
				msg := message.Message{}
				err = json.Unmarshal(m.Body, msg)
				if err != nil {
					logrus.Errorf("Could not unmarshal JSON: %s", err)
				} else {
					msgChan <- msg
				}
			}
		}

	}()

	return msgChan
}

func (q *Queue) Publish(topic string, message message.Message) {
	marshalledMessage, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("Could not marshal JSON: %s", err)
	} else {
		err = q.ch.Publish("input", topic, false, false, amqp.Publishing{
			Body:      marshalledMessage,
			Timestamp: time.Now(),
		})
		if err != nil {
			logrus.Errorf("Could not publish message to queue: %s", err)
		}
	}
}

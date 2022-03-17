package rabbitmq

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"io/ioutil"
	"os"
)

type ExchangeBinding struct {
	Destination string
	Key         string
	Source      string
	NoWait      bool
	Args        amqp.Table
}

type QueueBinding struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp.Table
}

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
	Bindings   []ExchangeBinding
	Queues     []RMQQueue
}
//name string, durable, autoDelete, exclusive, noWait bool, args Table
type RMQQueue struct {
	Name string
	Durable bool
	AutoDelete bool
	Exclusive bool
	NoWait bool
	Args amqp.Table
	Bindings []QueueBinding
}

func SetupTopology(ch amqp.Channel) {
	topologyFile := os.Getenv("TOPOLOGY_FILE")

	if topologyFile == "" {
		logrus.Fatalf("Please provide a topology file in the environment variable TOPOLOGY_FILE")
	}

	topologyBytes, err := ioutil.ReadFile(topologyFile)
	if err != nil {
		logrus.Fatalf("Could not read topology file: %s", err)
	}

	exchanges := make([]Exchange, 0)

	err = json.Unmarshal(topologyBytes, &exchanges)
	if err != nil {
		logrus.Fatalf("Could not load topology json: %s", err)
	}

	for _, exc := range exchanges {
		err := ch.ExchangeDeclare(
			exc.Name,
			exc.Kind,
			exc.Durable,
			exc.AutoDelete,
			exc.Internal,
			exc.NoWait,
			exc.Args,
		)
		if err != nil {
			logrus.Fatalf("Could not declare exchange: %s", err)
		}
		for _, binding := range exc.Bindings {
			err := ch.ExchangeBind(
				binding.Destination,
				binding.Key,
				binding.Source,
				binding.NoWait,
				binding.Args,
			)
			if err != nil {
				logrus.Fatalf("Could not bind exchange: %s", err)
			}
		}
		for _, queue := range exc.Queues {
			//name string, durable, autoDelete, exclusive, noWait bool, args Table
			_, err = ch.QueueDeclare(
				queue.Name,
				queue.Durable,
				queue.AutoDelete,
				queue.Exclusive,
				queue.NoWait,
				queue.Args,
				)
			if err != nil {
				logrus.Fatalf("Could not declare queue: %s", err)
			}
			//name, key, exchange string, noWait bool, args Table
			for _, binding := range queue.Bindings {
				err = ch.QueueBind(
					queue.Name,
					binding.Key,
					binding.Exchange,
					binding.NoWait,
					binding.Args,
					)
				if err != nil {
					logrus.Fatalf("Could not bind queue: %s", err)
				}
			}
		}
	}
}

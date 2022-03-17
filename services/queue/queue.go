package queue

import (
	"bitbucket.org/smaug-hosting/services/queue/internal/backends/rabbitmq"
	"bitbucket.org/smaug-hosting/services/queue/message"
	"errors"
)

type Queue interface {
	Subscribe(topic string) <-chan message.Message
	Publish(topic string, message message.Message)
}

type BackendType int

const (
	BackendRabbitMQ BackendType = iota
	BackendMock
)

var ErrQueueBackendUnavailable = errors.New("queue backend not found")

func GetQueueInstance(backend BackendType, args map[string]interface{}) (Queue, error) {
	switch backend {
	case BackendRabbitMQ:
		queue := rabbitmq.Queue{}
		queue.Setup(args)
		return &queue, nil
	default:
		return nil, ErrQueueBackendUnavailable
	}
}

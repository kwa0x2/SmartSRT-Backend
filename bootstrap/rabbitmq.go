package bootstrap

import (
	"log"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ() (*domain.RabbitMQ, error) {
	rabbitMQ := &domain.RabbitMQ{
		Done:    make(chan bool),
		Workers: make([]*domain.Worker, 0),
	}

	maxRetries := 30
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {

		if err := rabbitmq.Connect(rabbitMQ); err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v", err)

			if i == maxRetries-1 {
				return nil, err
			}

			time.Sleep(retryDelay)
			continue
		}

		log.Printf("uccessfully connected to RabbitMQ!")
		break
	}

	pool := &domain.ChannelPool{
		Channels: make(chan *amqp.Channel, domain.ChannelPoolSize),
	}

	for i := 0; i < domain.ChannelPoolSize; i++ {
		ch, err := rabbitMQ.Connection.Channel()
		if err != nil {
			return nil, err
		}

		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)
		if err != nil {
			return nil, err
		}

		pool.Channels <- ch
	}

	rabbitMQ.ChannelPool = pool
	go rabbitmq.HandleReconnect(rabbitMQ)

	return rabbitMQ, nil
}

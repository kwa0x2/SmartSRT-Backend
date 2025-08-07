package bootstrap

import (
	"log/slog"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ() (*domain.RabbitMQ, error) {
	logger := slog.Default()

	rabbitMQ := &domain.RabbitMQ{
		Done:    make(chan bool),
		Workers: make([]*domain.Worker, 0),
	}

	maxRetries := 30
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		if err := rabbitmq.Connect(rabbitMQ); err != nil {
			logger.Warn("RabbitMQ connection attempt failed",
				slog.Int("attempt", i+1),
				slog.Int("max_retries", maxRetries),
				slog.String("error", err.Error()),
			)

			if i == maxRetries-1 {
				logger.Error("RabbitMQ connection completely failed",
					slog.Int("total_attempts", maxRetries),
					slog.String("final_error", err.Error()),
				)
				return nil, err
			}

			time.Sleep(retryDelay)
			continue
		}

		logger.Info("RabbitMQ connection successful",
			slog.Int("attempt", i+1),
			slog.String("status", "connected"),
		)
		break
	}

	pool := &domain.ChannelPool{
		Channels: make(chan *amqp.Channel, domain.ChannelPoolSize),
	}

	for i := 0; i < domain.ChannelPoolSize; i++ {
		ch, err := rabbitMQ.Connection.Channel()
		if err != nil {
			logger.Error("RabbitMQ channel creation failed",
				slog.Int("channel_index", i),
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		err = ch.Qos(
			1,     // prefetch count
			0,     // prefetch size
			false, // global
		)
		if err != nil {
			logger.Error("RabbitMQ channel QoS configuration failed",
				slog.Int("channel_index", i),
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		pool.Channels <- ch
	}

	rabbitMQ.ChannelPool = pool
	go rabbitmq.HandleReconnect(rabbitMQ)

	logger.Info("RabbitMQ setup completed",
		slog.Int("channel_pool_size", domain.ChannelPoolSize),
		slog.String("status", "ready"),
	)

	return rabbitMQ, nil
}

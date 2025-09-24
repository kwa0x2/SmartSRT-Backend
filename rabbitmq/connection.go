package rabbitmq

import (
	"time"

	"github.com/kwa0x2/SmartSRT-Backend/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(r *domain.RabbitMQ, rabbitMQURI string) error {
	conn, err := amqp.Dial(rabbitMQURI)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	_, err = ch.QueueDeclare(
		domain.QueueConversions,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	r.Connection = conn
	r.Channel = ch
	r.IsConnected = true
	r.NotifyClose = make(chan *amqp.Error)
	r.Channel.NotifyClose(r.NotifyClose)

	return nil
}

func HandleReconnect(r *domain.RabbitMQ) {
	for {
		select {
		case <-r.Done:
			return
		case err := <-r.NotifyClose:
			if err != nil {
				r.IsConnected = false
				if r.ChannelPool != nil {
					r.ChannelPool.Close()
				}

				for {
					time.Sleep(domain.ReconnectDelay)
					if err := Connect(r, r.URI); err != nil {
						continue
					}
					ReinitializeWorkers(r)
					break
				}
			}
		}
	}
}

func Close(r *domain.RabbitMQ) {
	if !r.IsConnected {
		return
	}
	r.IsConnected = false

	if r.ChannelPool != nil {
		r.ChannelPool.Close()
	}

	close(r.Done)
	for _, worker := range r.Workers {
		close(worker.Done)
	}

	r.WorkerWg.Wait()

	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}

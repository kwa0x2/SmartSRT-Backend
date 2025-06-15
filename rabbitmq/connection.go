package rabbitmq

import (
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect(r *domain.RabbitMQ) error {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
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
					if err := Connect(r); err != nil {
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

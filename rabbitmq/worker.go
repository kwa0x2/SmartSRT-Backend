package rabbitmq

import (
	"encoding/json"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ReinitializeWorkers(r *domain.RabbitMQ) {
	r.Mu.RLock()
	workers := make([]*domain.Worker, len(r.Workers))
	copy(workers, r.Workers)
	r.Mu.RUnlock()

	for _, w := range workers {
		go StartWorker(w)
	}
}

func StartWorker(w *domain.Worker) {
	defer w.RabbitMQ.WorkerWg.Done()

	for {
		select {
		case <-w.Done:
			return
		default:
			if err := consume(w); err != nil {
				time.Sleep(domain.ReInitDelay)
			}
		}
	}
}

func consume(w *domain.Worker) error {
	ch, err := w.RabbitMQ.ChannelPool.Get()
	if err != nil {
		return err
	}
	defer w.RabbitMQ.ChannelPool.Put(ch)

	msgs, err := ch.Consume(
		w.Queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var convMsg domain.ConversionMessage
		if err = json.Unmarshal(msg.Body, &convMsg); err != nil {
			msg.Reject(false)
			continue
		}

		response, resErr := w.Handler(convMsg)
		if resErr != nil {
			msg.Reject(true)
			response = &domain.LambdaResponse{
				StatusCode: 500,
				Body: domain.LambdaBodyResponse{
					Message: resErr.Error(),
				},
			}
		} else {
			msg.Ack(false)
		}

		if msg.ReplyTo != "" {
			body, jsonErr := json.Marshal(response)
			if jsonErr != nil {
				continue
			}

			ch.Publish(
				"",          // exchange
				msg.ReplyTo, // routing key
				false,       // mandatory
				false,       // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					Body:          body,
					CorrelationId: msg.CorrelationId,
				},
			)
		}
	}

	return nil
}

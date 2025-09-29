package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kwa0x2/SmartSRT-Backend/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func StartWorkerPool(r *domain.RabbitMQ, numWorkers int, handler func(domain.ConversionMessage) (*domain.LambdaResponse, error)) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	for i := 0; i < numWorkers; i++ {
		worker := &domain.Worker{
			ID:       i + 1,
			Queue:    domain.QueueConversions,
			Handler:  handler,
			Done:     make(chan bool),
			RabbitMQ: r,
		}

		r.Workers = append(r.Workers, worker)
		r.WorkerWg.Add(1)
		go StartWorker(worker)
	}

	return nil
}

func PublishConversionMessage(r *domain.RabbitMQ, ctx context.Context, msg domain.ConversionMessage) (*domain.LambdaResponse, error) {
	ch, err := r.Connection.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	replyQueue, err := ch.QueueDeclare(
		"",    // random name
		false, // not durable
		true,  // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	replies, err := ch.Consume(
		replyQueue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(msg)
	if err != nil {
		ch.QueueDelete(replyQueue.Name, false, false, false)
		return nil, err
	}

	err = ch.Publish(
		"",                      // exchange
		domain.QueueConversions, // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			ReplyTo:       replyQueue.Name,
			CorrelationId: msg.FileID,
		},
	)
	if err != nil {
		ch.QueueDelete(replyQueue.Name, false, false, false)
		return nil, err
	}

	responseChan := make(chan *domain.LambdaResponse, 1)
	errorChan := make(chan error, 1)
	done := make(chan struct{})

	go func() {
		defer close(done)
		for reply := range replies {
			if reply.CorrelationId == msg.FileID {
				var response domain.LambdaResponse
				if err = json.Unmarshal(reply.Body, &response); err != nil {
					errorChan <- err
					return
				}
				responseChan <- &response
				return
			}
		}
	}()

	timeoutTimer := time.NewTimer(domain.MessageTimeout)
	defer timeoutTimer.Stop()

	select {
	case <-ctx.Done():
		ch.Cancel("", false)
		<-done
		ch.QueueDelete(replyQueue.Name, false, false, false)
		return nil, ctx.Err()
	case err := <-errorChan:
		ch.Cancel("", false)
		<-done
		ch.QueueDelete(replyQueue.Name, false, false, false)
		return nil, err
	case response := <-responseChan:
		ch.Cancel("", false)
		<-done
		ch.QueueDelete(replyQueue.Name, false, false, false)
		if response.StatusCode != 200 {
			response.Body.Message = "An error occurred. Please try again later or contact support."
		}
		return response, nil
	case <-timeoutTimer.C:
		return &domain.LambdaResponse{
			StatusCode: 202,
			Body: domain.LambdaBodyResponse{
				Message: "Your request is being processed. You will be notified by email when it's completed.",
			},
		}, nil
	}
}

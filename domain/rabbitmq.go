package domain

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	QueueConversions = "srt_conversions"

	ReconnectDelay  = 5 * time.Second
	ReInitDelay     = 2 * time.Second
	ResendDelay     = 5 * time.Second
	MessageTimeout  = 10 * time.Second
	ChannelPoolSize = 10
)

type ChannelPool struct {
	Channels chan *amqp.Channel
	Mu       sync.Mutex
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	select {
	case ch := <-p.Channels:
		return ch, nil
	default:
		return nil, fmt.Errorf("channel pool is empty")
	}
}

func (p *ChannelPool) Put(ch *amqp.Channel) {
	select {
	case p.Channels <- ch:
	default:
		ch.Close()
	}
}

func (p *ChannelPool) Close() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	close(p.Channels)
	for ch := range p.Channels {
		ch.Close()
	}
}

type ConversionMessage struct {
	UserID              bson.ObjectID `json:"user_id"`
	WordsPerLine        int           `json:"words_per_line"`
	Punctuation         bool          `json:"punctuation"`
	ConsiderPunctuation bool          `json:"consider_punctuation"`
	FileName            string        `json:"file_name"`
	FileID              string        `json:"file_id"`
	FileContent         []byte        `json:"file_content"`
	FileSize            int64         `json:"file_size"`
	FileDuration        float64       `json:"file_duration"`
	Email               string        `json:"email"`
}

type RabbitMQ struct {
	Connection  *amqp.Connection
	Channel     *amqp.Channel
	ChannelPool *ChannelPool
	Done        chan bool
	NotifyClose chan *amqp.Error
	IsConnected bool
	Workers     []*Worker
	WorkerWg    sync.WaitGroup
	Mu          sync.RWMutex
	URI         string
}

type Worker struct {
	ID       int
	Channel  *amqp.Channel
	Queue    string
	Handler  func(ConversionMessage) (*LambdaResponse, error)
	Done     chan bool
	RabbitMQ *RabbitMQ
}

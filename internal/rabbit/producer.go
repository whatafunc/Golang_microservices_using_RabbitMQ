package rabbit

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	key      string
}

func NewProducer(uri, exchange, key string) (*Producer, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Producer{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		key:      key,
	}, nil
}

func (p *Producer) Publish(body []byte) error {
	return p.channel.Publish(
		p.exchange,
		p.key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

func (p *Producer) Start(quit <-chan struct{}) {
	i := 0
	for {
		select {
		case <-quit:
			log.Println("producer shutting down...")
			return
		default:
			msg := []byte("message " + time.Now().Format(time.RFC3339))
			if err := p.Publish(msg); err != nil {
				log.Printf("failed to publish: %v", err)
				return
			}
			log.Printf("sent: %s", msg)
			i++
			time.Sleep(time.Second)
		}
	}
}

func (p *Producer) Shutdown() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}

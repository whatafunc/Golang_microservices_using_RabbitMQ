package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/app"
	// "github.com/whatafunc/Golang_Otus_Labs/hw12_13_14_15_16_calendar/internal/logger"
)

// Event represents the structure of an event.
type Event struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type Producer struct {
	app      *app.App
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	key      string
}

// NewProducer initializes the RabbitMQ producer and associates it with the app instance.
func NewProducer(a *app.App, uri, exchange, key string) (*Producer, error) {
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
		app:      a,
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		key:      key,
	}, nil
}

// Publish sends a raw message (JSON) to RabbitMQ.
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
			// msg := []byte("message " + time.Now().Format(time.RFC3339))
			// if err := p.Publish(msg); err != nil {
			// 	log.Printf("failed to publish: %v", err)
			// 	return
			// }
			// log.Printf("sent: %s", msg)
			if err := p.ListEventsDay(i); err != nil {
				log.Printf("failed to publish: %v", err)
				return
			}
			i++
			//logg.Info(fmt.Sprintf("✅ Producer sent day events batch %d", i))
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

// ListEventsDay generates and publishes event data for the day.
func (p *Producer) ListEventsDay(i int) error {
	// Example: Generate a list of events for the day
	// events := []Event{
	// 	{ID: 1, Name: "Event 1", Timestamp: time.Now()},
	// 	{ID: 2, Name: "Event 2", Timestamp: time.Now().Add(1 * time.Hour)},
	// }
	ctx := context.TODO()
	events, err := p.app.ListEvents(ctx, app.PeriodDay)
	if err != nil {
		// p.app.logger.Errorf("failed to list day events: %v", err)
		log.Printf("failed to list day events: %v", err)
		return fmt.Errorf("failed to list day events: %w", err)
	}
	log.Printf("checked list day events: %v", events)
	for _, event := range events {
		// Serialize the event to JSON
		// if err != nil {
		// 	return err
		// }
		msg, err := json.Marshal(event)
		if err != nil {
			log.Printf("failed to serialize event: %v", err)
			continue
		}

		// Publish the message
		if err := p.Publish(msg); err != nil {
			log.Printf("failed to publish event(%d): %v", i, err)
			continue
		}

		log.Printf("sent event(%d): %s", i, msg)
		time.Sleep(time.Second) // Simulate delay between messages
	}
	return nil
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var configFile string

// Message represents a generic message, independent of RabbitMQ
type Message struct {
	Body       []byte
	RoutingKey string
	Headers    map[string]interface{}
}

// ConsumerChannel is a generic message consumer
type ConsumerChannel interface {
	Consume() (<-chan Message, error)
}

type RabbitConsumer struct {
	channel  *amqp.Channel
	queue    string
	consumer string
}

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := LoadConfig(configFile)

	// Build AMQP URI from config
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	// connect to RabbitMQ
	conn, _ := amqp.Dial(amqpURI)
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()

	rc := &RabbitConsumer{
		channel:  ch,
		queue:    cfg.Queue,
		consumer: cfg.ConsumerTag,
	}

	msgs, err := rc.Consume()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan error)
	go handle(msgs, done) // run handler

	log.Println("Consumer running...")

	if cfg.Lifetime > 0 {
		log.Printf("running for %s", cfg.Lifetime)
		time.Sleep(time.Duration(cfg.Lifetime) * time.Second)
		log.Println("lifetime expired, shutting down...")
	} else {
		log.Printf("running forever")

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-done:
			log.Println("handler finished, shutting down...")
		case <-sigs:
			log.Println("received termination signal, shutting down...")
		}
	}

	// Always shutdown transport
	if err := rc.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func (rc *RabbitConsumer) Consume() (<-chan Message, error) {
	deliveries, err := rc.channel.Consume(
		rc.queue,    // queue name
		rc.consumer, // consumer tag
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return nil, err
	}

	// Convert amqp.Delivery → Message
	out := make(chan Message)
	go func() {
		for d := range deliveries {
			out <- Message{
				Body:       d.Body,
				RoutingKey: d.RoutingKey,
				Headers:    map[string]interface{}(d.Headers),
			}
		}
		close(out)
	}()

	return out, nil
}

func handle(msgs <-chan Message, done chan error) {
	for m := range msgs {
		log.Printf("got message: %s (by routing key=%s)", m.Body, m.RoutingKey)
	}
	done <- nil
}

func (rc *RabbitConsumer) Shutdown() error {
	if rc.channel != nil {
		if err := rc.channel.Close(); err != nil {
			return err
		}
	}
	return nil
}

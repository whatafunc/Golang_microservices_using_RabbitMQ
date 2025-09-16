package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/whatafunc/Golang_microservices_using_RabbitMQ/internal/rabbit"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	cfg, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Build AMQP URI from config
	amqpURI := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	producer, err := rabbit.NewProducer(amqpURI, "", cfg.Queue)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer producer.Shutdown()

	quit := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		close(quit)
	}()

	producer.Start(quit)
}

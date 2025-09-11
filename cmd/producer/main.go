package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	routingKey = flag.String("key", "test-key", "AMQP routing key")
	body       = flag.String("body", "foobar", "Body of message")
	reliable   = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
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
	if cfg.Sync {
		err := publish_sync_confirm(amqpURI, cfg.Exchange, cfg.ExchangeType, *routingKey, *body, *reliable)
		if err != nil {
			log.Fatalf("%s", err)
		}
	} else {
		err := publish(amqpURI, cfg.Exchange, cfg.ExchangeType, *routingKey, *body, *reliable)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}
	log.Printf("published %dB OK", len(*body))
}

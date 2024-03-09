package main

import (
	"log"
	"publisher/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Panic("read config FAILED")
	}

	brokerOptions := make([]*mqtt.ClientOptions, len(config.ServerConfig.BROKERS))
	for idx, broker := range config.ServerConfig.BROKERS {
		brokerOptions[idx] = init_mqtt(broker)
	}

	go PublishRoutine(brokerOptions)

	select {}
}

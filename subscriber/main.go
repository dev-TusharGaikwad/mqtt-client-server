package main

import (
	"log"
	"subscriber/config"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Fatal("read config FAILED")
		return
	}

	for _, broker := range config.ClientConfig.BROKERS {
		opts := init_mqtt(broker)
		subscribe(opts)
	}
	go DecompressSubscribedData()

	select {}
}

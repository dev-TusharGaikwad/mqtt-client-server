package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const CONFIG_FILE_PATH = "config/client-config.yaml"

type Options struct {
	BROKERS         []string `yaml:"BROKERS,omitempty"`
	CERTS_REQUIRED  bool     `yaml:"CERTS_REQUIRED,omitempty"`
	CERT_FILE       string   `yaml:"CERT_FILE,omitempty"`
	KEY_FILE        string   `yaml:"KEY_FILE,omitempty"`
	BROKER_CA       string   `yaml:"BROKER_CA,omitempty"`
	CLIENT_ID       string   `yaml:"CLIENT_ID,omitempty"`
	QOS             uint8    `yaml:"QOS,omitempty"`
	MAX_PACKET_SIZE int64    `yaml:"MAX_PACKET_SIZE,omitempty"`
	TOPICS          []string `yaml:"TOPICS,omitempty"`
}

var ClientConfig Options

func ReadConfig() error {
	log.Println()
	bytes, err := os.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		log.Fatalln("unable to read file...:", err)
		return err
	}

	err = yaml.Unmarshal(bytes, &ClientConfig)
	if err != nil {
		log.Fatalln("unable to unmarshal: ", err)
		return err
	}
	return nil
}

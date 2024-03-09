package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"subscriber/config"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var SubscribedDatChan = make(chan mqtt.Message, 10)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	select {
	case SubscribedDatChan <- msg:
		log.Printf("Recieved packet size(%d)\n", len(msg.Payload()))
	default:
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v\n", err)
}

func init_mqtt(broker string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	log.Println("Connecting to : ", broker)
	opts.AddBroker(broker)

	if config.ClientConfig.CERTS_REQUIRED {
		tlsConfig := NewTlsConfig()
		opts.SetTLSConfig(tlsConfig)
	}

	opts.SetClientID(config.ClientConfig.CLIENT_ID)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetKeepAlive(time.Second * 60)
	opts.ResumeSubs = true
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	return opts

}

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()

	ca, err := os.ReadFile(config.ClientConfig.BROKER_CA)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)
	clientKeyPair, err := tls.LoadX509KeyPair(config.ClientConfig.CERT_FILE, config.ClientConfig.KEY_FILE)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{clientKeyPair},
	}
}

func subscribe(opts *mqtt.ClientOptions) {
	var client mqtt.Client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, topic := range config.ClientConfig.TOPICS {
		token := client.Subscribe(topic, 1, messagePubHandler)
		token.Wait()
		log.Printf("Subscribed to topic: %s\n", topic)
	}
}

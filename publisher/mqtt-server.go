package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	random "math/rand"
	"os"
	"publisher/config"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

func init_mqtt(broker string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	log.Println("Connecting to : ", broker)
	opts.AddBroker(broker)

	if config.ServerConfig.CERTS_REQUIRED {
		tlsConfig := NewTlsConfig()
		opts.SetTLSConfig(tlsConfig)
	}
	// tlsConfig := NewTlsConfig()
	// opts.SetTLSConfig(tlsConfig)

	opts.SetClientID(config.ServerConfig.CLIENT_ID)
	opts.SetKeepAlive(time.Second * 60)
	opts.ResumeSubs = true
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	return opts

}

func NewTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := os.ReadFile(config.ServerConfig.BROKER_CA)
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)

	// Import client certificate/key pair
	clientKeyPair, err := tls.LoadX509KeyPair(config.ServerConfig.CERT_FILE, config.ServerConfig.KEY_FILE)
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

func PublishRoutine(opts []*mqtt.ClientOptions) {
	clients := make([]mqtt.Client, len(opts))

	for idx, ServerOpts := range opts {
		clients[idx] = mqtt.NewClient(ServerOpts)
		if token := clients[idx].Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	log.Println("Publishing to topic: ", config.ServerConfig.TOPICS)
	for {
		var BATCH_SIZE = random.Intn(int(config.ServerConfig.MAX_PACKET_SIZE))
		pktBatch, err := createBatch(BATCH_SIZE)

		if err != nil {
			log.Fatalln("invalid batch")
			continue
		}

		buf, err := CompressData(pktBatch)
		if err != nil {
			continue
		}
		PublishToAll(clients, buf)
		time.Sleep(time.Second * 1)
	}
}

var pubCnt int64 = 0

func PublishToAll(clients []mqtt.Client, buf []byte) {
	pubErrFlag := true
	timeTaken := time.Now()
	for _, client := range clients {
		for _, topic := range config.ServerConfig.TOPICS {
			if token := client.Publish(topic, byte(config.ServerConfig.QOS), false, buf); token.Wait() && token.Error() != nil {
				log.Printf("packet publish failed for %s\n", topic)
				pubErrFlag = false
			} else {
				pubErrFlag = true
				continue
			}
		}
		if pubErrFlag {
			pubCnt++
			log.Printf("published packets: %d, publish time: %s\n\n", pubCnt, time.Since(timeTaken))
		}
	}
}

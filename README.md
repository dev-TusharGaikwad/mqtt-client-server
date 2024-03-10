# mqtt-publish-subscribe
 
###### MIT License Copyright (c) 2024 Tushar Gaikwad

---

This repo contains golang implemetation of mqtt publisher and subscriber.


> `this implemetation uses protobuf to serialize and deserialize data.`

> `it uses lz4 compression algorithm (C library) to compress and decompress data`

> `available brokers are, ` 
> * `mqtt://test.mosquitto.org:1883`
> * `tcp://broker.emqx.io:1883` 

> **Note:** `you can add your own broker`.
> * [create your own broker](https://test.mosquitto.org/)

## HOW to generate certs:
* go to `certs-gen` directory and run script to generate certificates for server, client and broker
* > $ `./create_certs.sh`
* server and client certs will be automatically copied to certs directory of respective module.
* You can find broker certs in `etc/certs/broker/` 

## HOW to serialize your data:

* go to `proto_gen` dir and add proto message in .proto file.

* to generate `*.pb.go` run following command
```bash
$ ./build_proto.sh
```
* generated files will be copied to `ipc` dir of respective module.


## HOW to run publisher/subscriber:

* ### To run mqtt-publisher
    * goto `publisher` module and run following command.
        > `$ make build` <br> 
        `$ make run`
* ### To run mqtt-subscriber
    * got `subscriber` module and run following command.
        > `$ make build` <br> 
        `$ make run`

## HOW to configure your client/server
* config file for each module is present in the `config` dir of respective module
* You add `clients`/`servers` or you can enable certificate authentication.
* You can add multiple topics to publish/subscribe.

### Credits:
**Inintial Development**: Tushar Gaikwad <br>
**Maintainer**: Tushar Gaikwad
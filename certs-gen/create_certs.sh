#!/bin/bash

IP="192.168.1.22"
SUBJECT_CA="/C=SE/ST=Stockholm/L=Stockholm/O=himinds/OU=CA/CN=$IP"
SUBJECT_SERVER="/C=SE/ST=Stockholm/L=Stockholm/O=himinds/OU=Server/CN=$IP"
SUBJECT_CLIENT="/C=SE/ST=Stockholm/L=Stockholm/O=himinds/OU=Client/CN=$IP"

function generate_CA () {
   echo "$SUBJECT_CA"
   openssl req -x509 -nodes -sha256 -newkey rsa:2048 -subj "$SUBJECT_CA"  -days 365 -keyout ca.key -out ca.crt
}

function generate_server () {
   echo "$SUBJECT_SERVER"
   openssl req -nodes -sha256 -new -subj "$SUBJECT_SERVER" -keyout server.key -out server.csr
   openssl x509 -req -sha256 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365
}

function generate_client () {
   echo "$SUBJECT_CLIENT"
   openssl req -new -nodes -sha256 -subj "$SUBJECT_CLIENT" -out client.csr -keyout client.key 
   openssl x509 -req -sha256 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365
}

function copy_keys_to_broker () {
   sudo cp ca.crt ./etc/certs/
   sudo cp server.crt ./etc/certs/
   sudo cp server.key ./etc/certs/

   sudo chmod 777 ./etc/certs/*

   sudo cp -r ./etc/certs/ ../publisher/publisher-certs
   sudo chmod 777 ../publisher/publisher-certs/certs/*

   sudo cp -r ./etc/certs/ ../subscriber/subscriber-certs
   sudo chmod 777 ../subscriber/subscriber-certs/certs/*

}

function clean_certs() {
   rm *.crt *.key *.csr *.srl
}

generate_CA
generate_server
generate_client
copy_keys_to_broker
clean_certs
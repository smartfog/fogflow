#!/bin/bash

BOLD=$(tput bold)
CLEAR=$(tput sgr0)


# for root CA

echo -e "${BOLD}Generating RSA AES-256 Private Key for Root Certificate Authority${CLEAR}"
openssl genrsa  -out root_ca.key 4096

echo -e "${BOLD}Generating Certificate for Root Certificate Authority${CLEAR}"
openssl req -x509 -new -nodes -key root_ca.key -days 1825 -out root_ca.pem -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$1"


# for server

echo -e "${BOLD}Generating RSA Private Key for Server Certificate${CLEAR}"
openssl genrsa -out cloud_node.key 4096

echo -e "${BOLD}Generating Certificate Signing Request for Server Certificate${CLEAR}"
openssl req -new -key cloud_node.key -out cloud_node.csr -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$1"  

echo -e "${BOLD}Generating Certificate for Server Certificate${CLEAR}"
openssl x509 -req -in cloud_node.csr -CA root_ca.pem -CAkey root_ca.key -CAcreateserial -out cloud_node.pem -outform PEM -days 1825 


# for nginx server

echo -e "${BOLD}Generating RSA Private Key for Client Certificate${CLEAR}"
openssl genrsa -out client.key 4096

echo -e "${BOLD}Generating Certificate Signing Request for Client Certificate${CLEAR}"
openssl req -new -key client.key -out client.csr -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$2"

echo "Done!"


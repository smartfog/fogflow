#!/bin/bash

BOLD=$(tput bold)
CLEAR=$(tput sgr0)


# for root CA

echo -e "${BOLD}Generating RSA AES-256 Private Key for Root Certificate Authority${CLEAR}"
openssl genrsa  -out root_ca.key 4096

echo -e "${BOLD}Generating Certificate for Root Certificate Authority${CLEAR}"
openssl req -x509 -new -nodes -key root_ca.key -days 1825 -out root_ca.pem -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$1"


# for cloud-node

echo -e "${BOLD}Generating RSA Private Key for cloud-node Certificate${CLEAR}"
openssl genrsa -out cloud_node.key 4096

echo -e "${BOLD}Generating Certificate Signing Request for cloud-node Certificate${CLEAR}"
openssl req -new -key cloud_node.key -out cloud_node.csr -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$1"  

echo -e "${BOLD}Generating Certificate for cloud-node Certificate${CLEAR}"
openssl x509 -req -in cloud_node.csr -CA root_ca.pem -CAkey root_ca.key -CAcreateserial -out cloud_node.pem -outform PEM -days 1825 


# for edge-node

echo -e "${BOLD}Generating RSA Private Key for edge-node Certificate${CLEAR}"
openssl genrsa -out edge_node.key 4096

echo -e "${BOLD}Generating Certificate Signing Request for edge-node Certificate${CLEAR}"
openssl req -new -key edge_node.key -out edge_node.csr -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$1"

echo -e "${BOLD}Generating Certificate for edge-node Certificate${CLEAR}"
openssl x509 -req -in edge_node.csr -CA root_ca.pem -CAkey root_ca.key -CAcreateserial -out edge_node.pem -outform PEM -days 1825 

# for designer

echo -e "${BOLD}Generating RSA Private Key for Designer Certificate${CLEAR}"
openssl genrsa -out designer.key 4096

echo -e "${BOLD}Generating Certificate Signing Request for Designer Certificate${CLEAR}"
openssl req -new -key designer.key -out designer.csr -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=www.fogflow.org"

echo -e "${BOLD}Generating Certificate for Designer Certificate${CLEAR}"
openssl x509 -req -in designer.csr -CA root_ca.pem -CAkey root_ca.key -CAcreateserial -out designer.pem -outform PEM -days 1825 

echo "Done!"


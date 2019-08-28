#!/bin/bash

BOLD=$(tput bold)
CLEAR=$(tput sgr0)


# for cloud-node

echo -e "${BOLD}Generating RSA Private Key for cloud-node Certificate${CLEAR}"
openssl genrsa -out cloud_node.key 4096

echo -e "${BOLD}Generating Certificate Signing Request for cloud-node Certificate${CLEAR}"
openssl req -new -key cloud_node.key -out cloud_node.csr -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$1"  

echo -e "${BOLD}Generating Certificate for cloud-node Certificate${CLEAR}"
openssl x509 -req -in cloud_node.csr -CA root_ca.pem -CAkey root_ca.key -CAcreateserial -out cloud_node.pem -outform PEM -days 1825 


echo "Done!"


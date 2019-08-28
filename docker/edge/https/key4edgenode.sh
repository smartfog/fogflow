#!/bin/bash

BOLD=$(tput bold)
CLEAR=$(tput sgr0)

# for edge-node

echo -e "${BOLD}Generating RSA Private Key for edge-node Certificate${CLEAR}"
openssl genrsa -out edge_node.key 4096

echo -e "${BOLD}Generating Certificate Signing Request for edge-node Certificate${CLEAR}"
openssl req -new -key edge_node.key -out edge_node.csr -subj "/C=DE/ST=BW/L=Heidelberg/O=NEC/OU=NLE/CN=$1"

echo -e "${BOLD}Generating Certificate for edge-node Certificate${CLEAR}"
openssl x509 -req -in edge_node.csr -CA root_ca.pem -CAkey root_ca.key -CAcreateserial -out edge_node.pem -outform PEM -days 1825 

echo "Done!"


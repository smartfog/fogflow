#!/bin/bash

go get

#build the linux version (amd64) of master
env GOOS=linux GOARCH=amd64 go build  -a  -o change_config
docker build -t "rahafrouz/fogflow-prom" .


#!/bin/bash

go get

#build the linux version (amd64) of prometheus for fogflow

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -a -installsuffix cgo  -o PrometheusConfigUpdaterAPI
docker build -t "fogflow/prometheus" .

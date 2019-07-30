#!/bin/sh

cid=$(docker run -d -p 4545:4545 -p 9090:9090 -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus)

echo $cid
docker cp tgroups/ $cid:/etc/prometheus/
docker cp change_config $cid:/etc/prometheus/
#docker exec -u 0 -d $cid /etc/prometheus/change_config
docker exec -u 0 -d $cid chmod 777 /etc/prometheus/tgroups 
docker exec -u 0 -d $cid chmod 777 /etc/prometheus/tgroups/*
docker exec  -d $cid /etc/prometheus/change_config

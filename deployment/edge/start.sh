dir=$(pwd)

docker run -d --name=broker -v $(pwd)/config.json:/config.json -p 8080  fogflow/broker:latest
docker run -d --name=worker -v $(pwd)/config.json:/config.json -p 8080  fogflow/worker:latest


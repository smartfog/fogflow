wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/https/docker-compose.yml
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/https/config.json
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/https/nginx.conf

wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/https/key4cloudnode.sh
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/https/key4edgenode.sh

mkdir bind
cd bind
wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/https/bind/docker-compose.yml
mkdir data

cd ..





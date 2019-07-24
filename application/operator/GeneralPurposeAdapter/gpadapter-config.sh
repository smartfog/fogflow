echo "Configuring General Purpose Adapter..."

echo $1
echo $2

gunicorn --env IOT_BROKER_IP=$1 --env IOT_BROKER_PORT=$2 -b 0.0.0.0:1026 main

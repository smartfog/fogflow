#This script is used to stop and remove the docker containers for  edge components.

docker stop edgebroker && docker rm $_
docker stop edgeworker && docker rm $_

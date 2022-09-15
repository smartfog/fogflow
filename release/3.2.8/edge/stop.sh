#SCript to bring down the edge and remove corresponding docker containers
docker stop edgebroker && docker rm $_
docker stop edgeworker && docker rm $_

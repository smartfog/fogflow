#SCript to bring down the edge and remove corresponding docker containers
docker stop metricbeat && docker rm $_
docker stop edgebroker && docker rm $_
docker stop edgeworker && docker rm $_
docker stop pepEdge && docker rm $_

sh $(pwd)/delete.sh
if [ $? -eq 0 ]; then
    docker stop metricbeat && docker rm $_
    docker stop edgebroker && docker rm $_
    docker stop edgeworker && docker rm $_
else
     echo failed in delete Application
fi



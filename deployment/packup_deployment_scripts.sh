rm  ./core/fogflow-core.tar.gz 
tar -cf ./core/fogflow-core.tar.gz  ./core/*.json ./core/*.yml

rm ./edge/fogflow-edge.tar.gz 
tar -cf ./edge/fogflow-edge.tar.gz  ./edge/*.json ./edge/*.yml


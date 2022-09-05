kubectl create namespace fogflow-edge

kubectl  -n fogflow-edge create -f edge-configmap.yaml

kubectl  -n fogflow-edge create -f edge-broker.yaml
kubectl  -n fogflow-edge create -f edge-worker.yaml


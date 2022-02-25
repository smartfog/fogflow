kubectl  -n fogflow-edge delete -f edge-broker.yaml
kubectl  -n fogflow-edge delete -f edge-worker.yaml

kubectl  -n fogflow-edge delete -f edge-configmap.yaml

kubectl  delete namespace fogflow-edge

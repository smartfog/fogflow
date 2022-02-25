kfg="--kubeconfig=edge.yaml"

kubectl $kfg create namespace fogflow

kubectl  $kfg -n fogflow create -f edge-configmap.yaml

kubectl  $kfg -n fogflow create -f edge-broker.yaml
kubectl  $kfg -n fogflow create -f edge-worker.yaml


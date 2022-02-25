kfg="--kubeconfig=edge.yaml"

kubectl  $kfg -n fogflow delete -f edge-broker.yaml
kubectl  $kfg -n fogflow delete -f edge-worker.yaml

kubectl  $kfg -n fogflow delete -f edge-configmap.yaml

kubectl   $kfg delete namespace fogflow

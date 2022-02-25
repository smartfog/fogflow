kfg="--kubeconfig=cloud.yaml"

kubectl $kfg create namespace fogflow

kubectl $kfg apply -f designer-pv.yaml

kubectl $kfg -n fogflow create -f designer-pvc.yaml


kfg="--kubeconfig=cloud.yaml"

kubectl  $kfg -n fogflow create -f configmap.yaml

kubectl  $kfg -n fogflow create -f rabbitmq.yaml
kubectl  $kfg -n fogflow create -f nginx.yaml

kubectl  $kfg -n fogflow create -f discovery.yaml
kubectl  $kfg -n fogflow create -f broker.yaml
kubectl  $kfg -n fogflow create -f master.yaml
kubectl  $kfg -n fogflow create -f worker.yaml
kubectl  $kfg -n fogflow create -f designer.yaml


kfg="--kubeconfig=cloud.yaml"

kubectl  $kfg -n fogflow delete -f nginx.yaml

kubectl  $kfg -n fogflow delete -f discovery.yaml
kubectl  $kfg -n fogflow delete -f broker.yaml
kubectl  $kfg -n fogflow delete -f master.yaml
kubectl  $kfg -n fogflow delete -f worker.yaml
kubectl  $kfg -n fogflow delete -f designer.yaml

kubectl  $kfg -n fogflow delete -f rabbitmq.yaml

kubectl  $kfg -n fogflow delete -f configmap.yaml

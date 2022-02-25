kubectl -n fogflow delete -f nginx.yaml

kubectl -n fogflow delete -f discovery.yaml
kubectl -n fogflow delete -f broker.yaml
kubectl -n fogflow delete -f master.yaml
kubectl -n fogflow delete -f worker.yaml
kubectl -n fogflow delete -f designer.yaml

kubectl -n fogflow delete -f rabbitmq.yaml
kubectl -n fogflow delete -f configmap.yaml

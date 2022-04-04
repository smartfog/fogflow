kubectl -n fogflow create -f configmap.yaml

kubectl -n fogflow create -f rabbitmq.yaml

kubectl -n fogflow create -f discovery.yaml
kubectl -n fogflow create -f broker.yaml
kubectl -n fogflow create -f master.yaml
kubectl -n fogflow create -f worker.yaml
kubectl -n fogflow create -f designer.yaml

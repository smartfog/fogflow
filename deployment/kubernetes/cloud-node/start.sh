kubectl -n fogflow-cloud create -f configmap.yaml

kubectl -n fogflow-cloud create -f rabbitmq.yaml

kubectl -n fogflow-cloud create -f discovery.yaml
kubectl -n fogflow-cloud create -f broker.yaml
kubectl -n fogflow-cloud create -f master.yaml
kubectl -n fogflow-cloud create -f worker.yaml
kubectl -n fogflow-cloud create -f designer.yaml

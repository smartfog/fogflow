kubectl -n fogflow-cloud delete -f discovery.yaml
kubectl -n fogflow-cloud delete -f broker.yaml
kubectl -n fogflow-cloud delete -f master.yaml
kubectl -n fogflow-cloud delete -f worker.yaml
kubectl -n fogflow-cloud delete -f designer.yaml

kubectl -n fogflow-cloud delete -f rabbitmq.yaml
kubectl -n fogflow-cloud delete -f configmap.yaml

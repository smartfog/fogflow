kubectl delete -f nginx.yaml

kubectl delete -f configmap.yaml

kubectl delete -f discovery.yaml
kubectl delete -f broker.yaml
kubectl delete -f dgraph.yaml
kubectl delete -f rabbitmq.yaml
kubectl delete -f master.yaml
kubectl delete -f worker.yaml
kubectl delete -f designer.yaml


kubectl delete namespace fogflow



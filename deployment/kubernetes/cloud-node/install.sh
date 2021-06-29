mkdir -p dgraph

kubectl  create namespace fogflow

kubectl create -f serviceaccount.yaml

kubectl create -f configmap.yaml

kubectl create -f discovery.yaml
kubectl create -f broker.yaml
kubectl create -f dgraph.yaml
kubectl create -f rabbitmq.yaml
kubectl create -f master.yaml
kubectl create -f worker.yaml
kubectl create -f designer.yaml

kubectl create -f nginx.yaml




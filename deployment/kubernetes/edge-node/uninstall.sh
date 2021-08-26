microk8s.kubectl delete serviceaccount edge -n fogflow

microk8s.kubectl delete -f edge-configmap.yaml

microk8s.kubectl delete -f edge-serviceaccount.yaml

microk8s.kubectl delete -f edge-broker.yaml
microk8s.kubectl delete -f edge-worker.yaml

microk8s.kubectl delete ns fogflow


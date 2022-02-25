kfg="--kubeconfig=cloud.yaml"

kubectl   $kfg -n fogflow delete pvc designer-pvc
kubectl   $kfg delete pv  designer-pv

#kubectl $kfg delete -f serviceaccount.yaml

kubectl   $kfg delete namespace fogflow



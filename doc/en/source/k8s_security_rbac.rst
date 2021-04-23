***********************
Kubernetes Security 
***********************

Kubernetes provides many controls that can greatly improve an application's security. In order to use any of those methods provided by kubernetes, one need to properly configure the apiserver. **Role Based Access Control (RBAC)** is one such security implementation. RBAC is a method of regulating access to computer or network resources based on the roles of individual users within an organization. RBAC authorization uses the **rbac.authorization.k8s.io** API group to drive authorization decisions, allowing user to dynamically configure policies through the Kubernetes API.

RBAC Impementation over Cloud Node Kubernetes Cluster
--------------------------------------------------------
________________________________________________________________

It is assumed that kubernetes cluter is setup and running at cloud node. Inorder to setup RBAC in cloud node download and extract cloud-chart, configure **config.json** file as shown `above`_. 

.. _`above`: https://github.com/smartfog/fogflow/development/doc/en/source/k8sIntegration.rst#deploy-fogflow-cloud-components-on-k8s-environment


values.yaml Configurations over Cloud Node Kubernetes Cluster
-------------------------------------------------------------------

values.yaml can be accessed from fogflow repository using **"fogflow/helm/cloud-chart/values.yaml"** path.

- Configure the namespace and service account name in values.yaml file as shown below:

.. code-block:: console

   
   #Kubernetes namespace of FogFlow components 
   namespace: fogflow   //CAN BE CHANGED AS PER USER'S NEED

   #replicas will make sure that no. of replicaCount mention in values.yaml
   # are running all the time for the deployment
   replicaCount: 1

   serviceAccount:
   # Specifies whether a service account should be created
   create: true
   # Annotations to add to the service account
   annotations: {}
   # The name of the service account to use.
   # If not set and create is true, a name is generated using the fullname template
   name: "fogflow-dns"   //CAN BE CHANGED AS PER USER'S NEED

        
- On deploying this chart using helm, the **namespace** is created with name **fogflow**  and inside that a **sericeaccount** is created with name **fogflow-dns**. Once these namespace and serviceaccount is created, next roles and their rolebindings are created. The table lists the created roles and rolebinding. 

+--------------------+----------------+----------------------+
|     Roles          |  RoleBindings  |    Scope             |
+--------------------+----------------+----------------------+
| fogflow-root-role  |   RootUser     |  Cluster             |
+--------------------+----------------+----------------------+
| fogflow-admin-role |   Admin        |  fogflow - namespace |
+--------------------+----------------+----------------------+
| fogflow-user-role  |   EndUser      |  fogflow - namespace |
+--------------------+----------------+----------------------+

- To verify the creation of above resources, use following commands:

.. code-block:: console

   $kubectl get ns 

.. figure:: figures/ns.png

.. code-block:: console

   $kubectl get rolebindings --namespace=fogflow

.. figure:: figures/rbaccloud.png

Steps To Add Users in Cloud Node Kubernetes Cluster
-------------------------------------------------------

To add users in kubernetes cluster at cloud node, follow below steps:

1. Certificate Generation And Root User Addition
--------------------------------------------------

**Step 1**: Generate User's private key, using below command.

.. code-block:: console

   $openssl genrsa -out RootUser1.key 2048

**Step 2**: Generate User's certificate signing request using below commands.

.. code-block:: console

   $openssl req -new -key RootUser1.key -out RootUser1.csr -subj "/CN=RootUser1/O=RootUser"

   #the tag "/O=RootUser" defines the rolebinding, so enter carefully

**Step 3**: Generate User's certificate using below command.

.. code-block:: console

   $openssl x509 -req -in RootUser1.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key  -CAcreateserial -out RootUser1.crt -days 365

   #The "-day" tag justifies the no of days for which user's certificate will be valid. so it can be changed accordingly.

**Step 4**: To add user to kubernetes cluster, use following command.

.. code-block:: console

   $kubectl config set-credentials RootUser1 --client-certificate /root/RootUser/RootUser1.crt --client-key /root/RootUser/RootUser1.key

Note: The tags **--client-certificate** is followed by the path where user's private key is kept and **--client-key** is followed by path where user's certificate is kept. To verify added user, use below command.

.. code-block:: console

   $kubectl config view

.. figure:: figures/addedrootuser.png

**Step 5**: Set the context in kubeconfig to recently added user using following command.

.. code-block:: console

   $kubectl config set-context RootUser-context1 --cluster=kubernetes --namespace=fogflow --user=RootUser1

Note: set the value of namespace according to the value mentioned in values.yaml. Here **RootUser-context1** is the new context set for RootUser1.

**Step 6**: Now verify the permissions RootUser1 has by using various kubectl commands with above context as shown below.

.. code-block:: console

   $kubectl get node --context=RootUser-context1

   $kubectl delete pods "any pod name" --context=RootUser-context1

   $kubectl get pods --context=RootUser-context1

   $kubectl get pods --namespace=fogflow --context=RootUser-context1


.. figure:: figures/addedrootuseroutput.png


2. Certificate Generation And Admin User Addition
--------------------------------------------------

**Step 1**: Generate User's private key, using below command.

.. code-block:: console

   $openssl genrsa -out AdminUser1.key 2048

**Step 2**: Generate User's certificate signing request using below commands.

.. code-block:: console

   $openssl req -new -key AdminUser1.key -out AdminUser1.csr -subj "/CN=AdminUser1/O=Admin"

   #the tag "/O=Admin" defines the rolebinding, so enter carefully

**Step 3**: Generate User's certificate using below command.

.. code-block:: console

   $openssl x509 -req -in AdminUser1.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key  -CAcreateserial -out AdminUser1.crt -days 365

   #The "-day" tag justifies the no of days for which user's certificate will be valid. so it can be changed accordingly.

**Step 4**: To add user to kubernetes cluster, use following command.

.. code-block:: console

   $kubectl config set-credentials AdminUser1 --client-certificate /root/AdminUser/AdminUser1.crt --client-key /root/AdminUser/AdminUser1.key

Note: The tags **--client-certificate** is followed by the path where user's private key is kept and **--client-key** is followed by path where user's certificate is kept. To verify added user, use below command.

.. code-block:: console

   $kubectl config view

.. figure:: figures/addedadminuser.png

**Step 5**: Set the context in kubeconfig to recently added user using following command.

.. code-block:: console

   $kubectl config set-context AdminUser-context1 --cluster=kubernetes --namespace=fogflow --user=AdminUser1

Note: set the value of namespace according to the value mentioned in values.yaml. Here **AdminUser-context1** is the new context set for RootUser1.

**Step 6**: Now verify the permissions RootUser1 has by using various kubectl commands with above context as shown below.

.. code-block:: console

   $kubectl get node --context=AdminUser-context1

   $kubectl delete pods "any pod name" --context=AdminUser-context1

   $kubectl get pods --context=AdminUser-context1

   $kubectl get pods --namespace=fogflow --context=AdminUser-context1

.. figure:: figures/addedadminuseroutput.png


3. Certificate Generation And End User Addition
--------------------------------------------------

**Step 1**: Generate User's private key, using below command.

.. code-block:: console

   $openssl genrsa -out EndUser1.key 2048

**Step 2**: Generate User's certificate signing request using below commands.

.. code-block:: console

   $openssl req -new -key EndUser1.key -out EndUser1.csr -subj "/CN=EndUser1/O=EndUser"

   #the tag "/O=EndUser" defines the rolebinding, so enter carefully

**Step 3**: Generate User's certificate using below command.

.. code-block:: console

   $openssl x509 -req -in EndUser1.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key  -CAcreateserial -out EndUser1.crt -days 365

   #The "-day" tag justifies the no of days for which user's certificate will be valid. so it can be changed accordingly.

**Step 4**: To add user to kubernetes cluster, use following command.

.. code-block:: console

   $kubectl config set-credentials EndUser1 --client-certificate /root/EndUser/EndUser1.crt --client-key /root/EndUser/EndUser1.key

Note: The tags **--client-certificate** is followed by the path where user's private key is kept and **--client-key** is followed by path where user's certificate is kept. To verify added user, use below command.

.. code-block:: console

   $kubectl config view

.. figure:: figures/addedenduser.png

**Step 5**: Set the context in kubeconfig to recently added user using following command.

.. code-block:: console

   $kubectl config set-context EndUser-context1 --cluster=kubernetes --namespace=fogflow --user=EndUser1

Note: set the value of namespace according to the value mentioned in values.yaml. Here **EndUser-context1** is the new context set for RootUser1.

**Step 6**: Now verify the permissions RootUser1 has by using various kubectl commands with above context as shown below.

.. code-block:: console

   $kubectl get node --context=EndUser-context1

   $kubectl delete pods "any pod name" --context=EndUser-context1

   $kubectl get pods --context=EndUser-context1

   $kubectl get pods --namespace=fogflow --context=EndUser-context1


.. figure:: figures/addedenduseroutput.png


RBAC Implementation over Edge Node Microk8s Kubernetes Cluster
----------------------------------------------------------------
______________________________________________________________________

It is assumed that kubernetes cluter is setup and running at cloud node. Inorder to setup RBAC in cloud node download and extract edge-chart,configure **config.json** file as shown `above`_.

.. _`above`: https://github.com/smartfog/fogflow/development/doc/en/source/k8sIntegration.rst#deploying-edge-chart-with-microk8s-and-helm


values.yaml Configurations over Edge Node Kubernetes Cluster
----------------------------------------------------------------------

values.yaml can be accessed from fogflow repository using **"fogflow/helm/edge-chart/values.yaml"** path.

- Configure the namespace and service account name in values.yaml file as shown below:

.. code-block:: console

   
   #Kubernetes namespace of FogFlow components 
   namespace: fogflow   //CAN BE CHANGED AS PER USER'S NEED

   #replicas will make sure that no. of replicaCount mention in values.yaml
   # are running all the time for the deployment
   replicaCount: 1

   serviceAccount:
   # Specifies whether a service account should be created
   create: true
   # Annotations to add to the service account
   annotations: {}
   # The name of the service account to use.
   # If not set and create is true, a name is generated using the fullname template
   name: "fogflow-dns"   //CAN BE CHANGED AS PER USER'S NEED

        
- On deploying this chart using helm, the **namespace** with name **fogflow** is created and inside that a **sericeaccount** with name **fogflow-dns** is created. Once these namespace and serviceaccount is created, next roles and their rolebindings are created. The table lists the created roles and rolebinding. 

+--------------------+----------------+----------------------+
|     Roles          |  RoleBindings  |    Scope             |
+--------------------+----------------+----------------------+
| fogflow-root-role  |   RootUser     |  Cluster             |
+--------------------+----------------+----------------------+
| fogflow-admin-role |   Admin        |  fogflow - namespace |
+--------------------+----------------+----------------------+
| fogflow-user-role  |   EndUser      |  fogflow - namespace |
+--------------------+----------------+----------------------+

- To verify the creation of above resources, use following commands:

.. code-block:: console

   $mirok8s.kubectl get ns 

.. figure:: figures/nsedge.png

.. code-block:: console

   $microk8s.kubectl get rolebindings --namespace=fogflow

.. figure:: figures/rbacedge.png

Steps to Add Users in Edge Node Kubernetes Cluster
-----------------------------------------------------------

To add users in kubernetes cluster at edge node, follow below steps:

1. Certificate Generation And Root User Addition
--------------------------------------------------

**Step 1**: Generate User's private key, using below command.

.. code-block:: console

   $openssl genrsa -out RootUser1.key 2048

**Step 2**: Generate User's certificate signing request using below commands.

.. code-block:: console

   $openssl req -new -key RootUser1.key -out RootUser1.csr -subj "/CN=RootUser1/O=RootUser"

   #the tag "/O=RootUser" defines the rolebinding, so enter carefully

**Step 3**: Generate User's certificate using below command.

.. code-block:: console

   $openssl x509 -req -in RootUser1.csr -CA /var/snap/microk8s/current/certs/ca.crt -CAkey /var/snap/microk8s/current/certs/ca.key  -CAcreateserial -out RootUser1.crt -days 365

   #The "-day" tag justifies the no of days for which user's certificate will be valid. so it can be changed accordingly.

**Step 4**: To add user to kubernetes cluster, use following command.

.. code-block:: console

   $microk8s.kubectl config set-credentials RootUser1 --client-certificate /root/RootUser/RootUser1.crt --client-key /root/RootUser/RootUser1.key

Note: The tags **--client-certificate** is followed by the path where user's private key is kept and **--client-key** is followed by path where user's certificate is kept. To verify added user, use below command.

.. code-block:: console

   $microk8s.kubectl config view

.. figure:: figures/addedrootuseredge.png

**Step 5**: Set the context in kubeconfig to recently added user using following command.

.. code-block:: console

   $microk8s.kubectl config set-context RootUser1-context --cluster=microk8s-cluster --namespace=fogflow --user=RootUser1

Note: set the value of namespace according to the value mentioned in values.yaml. Here **RootUser-context1** is the new context set for RootUser1.

**Step 6**: Now verify the permissions RootUser1 has by using various kubectl commands with above context as shown below.

.. code-block:: console

   $microk8.kubectl get node --context=RootUser1-context

   $microk8.kubectl delete pods "any pod name" --context=RootUser1-context

   $microk8s.kubectl get pods --context=RootUser1-context

   $microk8s.kubectl get pods --namespace=fogflow --context=RootUser1-context


.. figure:: figures/addedrootuseredgeoutput.png   

2. Certificate Generation And Admin User Addition
--------------------------------------------------

**Step 1**: Generate User's private key, using below command.

.. code-block:: console

   $openssl genrsa -out AdminUser1.key 2048

**Step 2**: Generate User's certificate signing request using below commands.

.. code-block:: console

   $openssl req -new -key AdminUser1.key -out AdminUser1.csr -subj "/CN=AdminUser1/O=Admin"

   #the tag "/O=Admin" defines the rolebinding, so enter carefully

**Step 3**: Generate User's certificate using below command.

.. code-block:: console

   $openssl x509 -req -in AdminUser1.csr -CA /var/snap/microk8s/current/certs/ca.crt -CAkey /var/snap/microk8s/current/certs/ca.key  -CAcreateserial -out AdminUser1.crt -days 365

   #The "-day" tag justifies the no of days for which user's certificate will be valid. so it can be changed accordingly.

**Step 4**: To add user to kubernetes cluster, use following command.

.. code-block:: console

   $microk8s.kubectl config set-credentials AdminUser1 --client-certificate /root/AdminUser/AdminUser1.crt --client-key /root/AdminUser/AdminUser1.key

Note: The tags **--client-certificate** is followed by the path where user's private key is kept and **--client-key** is followed by path where user's certificate is kept. To verify added user, use below command.

.. code-block:: console

   $microk8s.kubectl config view

.. figure:: figures/addedadminuseredge.png

**Step 5**: Set the context in kubeconfig to recently added user using following command.

.. code-block:: console

   $microk8s.kubectl config set-context AdminUser-context1 --cluster=microk8s-cluster --namespace=fogflow --user=AdminUser1

Note: set the value of namespace according to the value mentioned in values.yaml. Here **AdminUser-context1** is the new context set for RootUser1.

**Step 6**: Now verify the permissions RootUser1 has by using various kubectl commands with above context as shown below.

.. code-block:: console

   $microk8s.kubectl get node --context=AdminUser-context1

   $microk8s.kubectl delete pods "any pod name" --context=AdminUser-context1

   $microk8s.kubectl get pods --context=AdminUser-context1

   $microk8s.kubectl get pods --namespace=fogflow --context=AdminUser-context1


.. figure:: figures/addedadminuseredgeoutput.png


3. Certificate Generation And End User Addition
--------------------------------------------------

**Step 1**: Generate User's private key, using below command.

.. code-block:: console

   $openssl genrsa -out EndUser1.key 2048

**Step 2**: Generate User's certificate signing request using below commands.

.. code-block:: console

   $openssl req -new -key EndUser1.key -out EndUser1.csr -subj "/CN=EndUser1/O=EndUser"

   #the tag "/O=EndUser" defines the rolebinding, so enter carefully

**Step 3**: Generate User's certificate using below command.

.. code-block:: console

   $openssl x509 -req -in EndUser1.csr -CA /var/snap/microk8s/current/certs/ca.crt -CAkey /var/snap/microk8s/current/certs/ca.key  -CAcreateserial -out EndUser1.crt -days 365

   #The "-day" tag justifies the no of days for which user's certificate will be valid. so it can be changed accordingly.

**Step 4**: To add user to kubernetes cluster, use following command.

.. code-block:: console

   $microk8s.kubectl config set-credentials EndUser1 --client-certificate /root/EndUser/EndUser1.crt --client-key /root/EndUser/EndUser1.key

Note: The tags **--client-certificate** is followed by the path where user's private key is kept and **--client-key** is followed by path where user's certificate is kept. To verify added user, use below command.

.. code-block:: console

   $microk8s.kubectl config view

.. figure:: figures/addedenduseredge.png

**Step 5**: Set the context in kubeconfig to recently added user using following command.

.. code-block:: console

   $microk8s.kubectl config set-context EndUser1-context --cluster=microk8s-cluster --namespace=fogflow --user=EndUser1

Note: set the value of namespace according to the value mentioned in values.yaml. Here **EndUser-context1** is the new context set for RootUser1.

**Step 6**: Now verify the permissions RootUser1 has by using various kubectl commands with above context as shown below.

.. code-block:: console

   $microk8s.kubectl get node --context=EndUser1-context

   $micr0k8s.kubectl delete pods "any pod name" --context=EndUser1-context

   $microk8s.kubectl get pods --context=EndUser1-context

   $microk8s.kubectl get pods --namespace=fogflow --context=EndUser1-context

.. figure:: figures/addedenduseredgeoutput.png

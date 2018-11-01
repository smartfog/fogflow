Set up a docker registry for image management
==============================================


Create a Self Signed Certificate
-------------------------------------
You need to create a self signed certificate on your server to use it for the private Docker Registry.

.. code-block:: console    

	mkdir registry_certs
	openssl req -newkey rsa:4096 -nodes -sha256 \
                -keyout registry_certs/domain.key -x509 -days 356 \
                -out registry_certs/domain.cert
	ls registry_certs/

Finally you have two files:

- domain.cert – this file can be handled to the client using the private registry
- domain.key – this is the private key which is necessary to run the private registry with TLS



Run the Private Docker Registry with TLS
-----------------------------------------
Now we can start the registry with the local domain certificate and key file:

.. code-block:: console    

	docker run -d -p 5000:5000 \
		 	-v $(pwd)/registry_certs:/certs \
 		 	-e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.cert \
 		 	-e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
 			--restart=always --name registry registry:2

Here we map the folder /registry_certs as an volume into the docker registry container. 
We use environment variables pointing to the certificate and key file.

Now you can push your local image into the new registry:

.. code-block:: console    

	docker push localhost:5000/proxy:1.0.0


Access the Remote Registry form a fog node
------------------------------------------------

Now as the private registry is started with TLS Support you can access the registry from any client which has the domain certificate.
Therefor the certificate file “domain.cert” must be located on the client in a file

.. code-block:: console    

	/etc/docker/certs.d/<registry_address>/ca.cert

Where <registry_address> is the server host name. After the certificate was updated you need to restart the local docker daemon:


.. code-block:: console    

	mkdir -p /etc/docker/certs.d/dock01:5000 
	cp domain.cert /etc/docker/certs.d/dock01:5000/ca.crt
	service docker restart
	
	
Now finally you can push you images into the new private registry:

.. code-block:: bash

	docker tag imixs/proxy dock01:5000/proxy:dock01
	docker push dock01:5000/proxy:dock01
	

Start Docker Registry Frontend
------------------------------------------------

The project konradkleine/docker-registry-frontend provides a cool web front-end which can be used to simplify access to the registry through a web browser.
The docker-registry-frontend can be started as a docker container. Assuming your registry runs on

https://yourserver.com:5000

use the following docker run command to start the frontend container:

.. code-block:: console    

	docker run \
 		-d \
 		-e ENV_DOCKER_REGISTRY_HOST=yourserver.com \
 		-e ENV_DOCKER_REGISTRY_PORT=5000 \
 		-e ENV_DOCKER_REGISTRY_USE_SSL=1 \
 		-p 0.0.0.0:80:80 \
 		konradkleine/docker-registry-frontend:v2
		
You can now access your registry via web browser url:

http://localhost:80/
Build everything from the source code
=========================================

FogFlow can be build and installed on Linux for both ARM and X86 processors (32bits and 64bits)
Install dependencies
--------------------------

#. To build FogFlow, first install the following dependencies.

        - install nodejs version 16 and above: please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-16-04

        - install Go and set up PATH environment variable: please follow the instruction at https://go.dev/doc/install
                
        - install git client: please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-git-on-ubuntu-16-04

        - install Docker CE: please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04

                .. note:: all the scripts are prepared under the assumption that docker can be run without sudo.

#. To check out the code repository

        .. code-block:: bash

                cd /home/smartfog/go/src/
                git clone https://github.com/smartfog/fogflow.git


#. To build all components from the source code

.. note:: Before building the components make sure that you are in the same directory in which you have check out the code from github repository as explained in the Step. 2 above. 

* Discovery: To build discovery use the following method: 

        .. code-block:: bash

                cd fogflow/discovery
                go get; go build
                ./discovery

                .. note:: make sure to configure **my_hostip** config.json.

* Broker: To build broker use the following method: 

        .. code-block:: bash

                cd fogflow/broker
                go get; go build
                ./broker

                .. note:: make sure to configure **my_hostip** config.json.

* Master: To build master use the following method: 

        .. code-block:: bash

                cd fogflow/master
                go get; go build
                ./master

                .. note:: make sure to configure **my_hostip** config.json.

* Worker: To build worker use the following method: 

        .. code-block:: bash

                cd fogflow/worker
                go get; go build
                ./worker

                .. note:: make sure to configure **my_hostip** and as **container_management** in config.json.

* Designer: To build designer use the following method: 

        .. code-block:: bash

                cd fogflow/discovery
                npm install
                node main.js

                .. note:: make sure to configure **my_hostip** config.json.


* Rabbitmq: To run rabbitmq use the command: 

        .. code-block:: bash

                docker run -it -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 --env RABBITMQ_DEFAULT_USER=admin --env RABBITMQ_DEFAULT_PASS=mypass  rabbitmq:3.8-management

**To check if all the Fogflow is build successfully:** Open Fogflow dashboard using this address "http://<externalIPs>:80".

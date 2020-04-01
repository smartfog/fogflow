Build FogFlow components from the source code
=========================================

FogFlow can be build and installed on Linux for both ARM and X86 processors (32bits and 64bits). 

Install dependencies
--------------------

#. To build FogFlow, first install the following dependencies.

	- install git client: please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-git-on-ubuntu-16-04
	
	- install Docker CE: please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04
	
		.. note:: all the scripts are prepared under the assumption that docker can be run without sudo.
	

	- install the latest version of golang(>v1.9): please download and install it according to the information at https://golang.org/doc/install

	- install nodejs(>6.11) and npm (>3.10): please download and install them according to the information at https://nodejs.org/en/download/


#. To check the installed version


	.. code-block:: bash

		go version   #output  go version go1.9 linux/amd64 
  		nodejs -v    #output 	v6.10.2
  		npm -v       #output  3.10.10


#. To set the environment variable GOPATH


	.. note:: GOPATH defines the workspace for go-based projects. Please note that the go workspace folder must have a "src" folder and the fogflow code repository must be cloned into this "src" folder. 
		For example, assume that your home folder is "/home/smartfog" and then you create a new folder "go" as your workspace. 
		In this case, you must create a "src" (must be exactly this name) under "/home/smartfog/go" first 
		and then check out the FogFlow code repository within the "/home/smartfog/go/src" folder.

	.. code-block:: bash	

		export GOPATH="/home/smartfog/go"


#. To check out the code repository

	.. code-block:: bash	
		
		cd /home/smartfog/go/src/	
		git clone https://github.com/smartfog/fogflow.git
		
		
#. To build all components from the source code as below


Build IoT Discovery
------------------------

	- to build the native executable program
	
		.. code-block:: bash	
			
			# go the discovery folder
			cd /home/smartfog/go/src/fogflow/discovery
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- to build the docker image, 

		.. code-block:: bash			
		
			# Simply ./build  can be run to perform the following commands
		
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/discovery" .										
		
			
Build IoT Broker
--------------------------

	- to build the native executable program
	
		.. code-block:: bash	
			
			# go the broker folder
			cd /home/smartfog/go/src/fogflow/broker
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- to build the docker image
		
		.. code-block:: bash			
		
			# simply ./build can be run to perform the following commands		
				
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/broker" .			



Build Topology Master
--------------------------

	- to build the native executable program
	
		.. code-block:: bash	
			
			# go the master folder
			cd /home/smartfog/go/src/fogflow/master
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- to build the docker image
		
		.. code-block:: bash							
		
			# simply ./build can be run to perform the following commands		
					
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/master" .			



Build Worker
--------------------------

	- to build the native executable program
	
		.. code-block:: bash	
			
			# go the worker folder
			cd /home/smartfog/go/src/fogflow/worker
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- to build the docker image
		
		.. code-block:: bash	
					
			# simply ./build  can be run to perform the following commands									
			
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/worker" .			


Build Task Designer
--------------------------

	- to install third-party library dependencies
	
		.. code-block:: bash	
			
			# go the designer folder
			cd /home/smartfog/go/src/fogflow/designer
			
			# install all required libraries
			npm install
	
	- to build the docker image
		
		.. code-block:: bash	
		
			# simply ./build can be run to perform the following commands					

			# install all required libraries
			npm install
			
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/designer"  .






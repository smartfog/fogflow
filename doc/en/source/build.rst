Build everything from the source code
=========================================

FogFlow can be build and installed on Linux for both ARM and X86 processors (32bits and 64bits). 

Install dependencies
--------------------------

#. To build FogFlow, first install the following dependencies.

	- install git client: please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-git-on-ubuntu-16-04
	
	- install Docker CE: please follow the instruction at https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04
	
		.. note:: all the scripts are prepared under the assumption that docker can be run without sudo.
	


#. To check out the code repository

	.. code-block:: bash	
		
		cd /home/smartfog/go/src/	
		git clone https://github.com/smartfog/fogflow.git
		
		
#. To build all components from the source code with multistage building

	.. code-block:: bash	
		
		./build multistage

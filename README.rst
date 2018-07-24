FogFlow
==========================

FogFlow is a standard-based data processing framework for service providers to easily programming and managing IoT services over cloud and edges. Below are the motivation, functionalities, and benefits of FogFlow. 

* Why do we need FogFlow?
	- the cost of a cloud-only solution is too high to run a large scale IoT system with >1000 geo-distributed devices
	- many IoT services require fast response time, such as <10ms end-to-end latency
	- service providers are facing huge complexity and cost to fast design and deploy their IoT services in a cloud-edge environment	- business demands are changing fast over time and service providers need to try out and release any new services over their shared cloud-edge infrastructure at a fast speed
	- lack of programming model to fast design and deploy IoT services over geo-distributed ICT infrastructure
	- lack of interoperability and openness to share and reuse data and dervied results across various applications	 

* What does FogFlow provide?
	- efficient programming model: programming a service is like building lego blocks 
	- dynamic service orchestration: launching necessary data processing only when it is required
	- optimized task deployment: assigning tasks between cloud and edges based on the locality of producers and consumers
	- scalable context management: allowing flexible information exchanging (both topic-based and scope-based) between producers and consumers

* How can customers benefit from FogFlow? 
	- fast time-to-market when realizing and releasing new services over the shared, geo-distributed ICT infrastructure
	- reduced operation cost and management complexity when operating variou services
	- being able to provide services that require low latency and fast response time


More Information
----------------

- `Tutorial`_
- `IoT-J paper`_

.. _`Tutorial`: http://fogflow.readthedocs.io/en/latest/index.html
.. _`IoT-J paper`: http://ieeexplore.ieee.org/document/8022859/

License
----------------
FogFlow is licensed under BSD-4-Clause.

Overview
=================

In This part of the document a conceptual overview of FogFlow and how to use FogFlow in developing any instance is being covered.
FogFlow is a cloud and edge environment to orchestrate dynamic `NGSI`_-based (Next Generation Service Interface - based) data processing
flows on-demand between producers and consumers for providing timely results to make fast actions. A context producer will be a sensor 
based device whereas a consumer is an Actuator device that will recieve command to perform some action.

.. _`NGSI`: https://knowage.readthedocs.io/en/6.1.1/user/NGSI/README/index.html

FogFlow can carry out IoT service orchestration decisions in a decentralized and autonomous manner. This means each FogFlow edge node can make its own decisions only based on a local context view. This way the majority of workloads can be directly handled at edges without always relying on the central cloud. With this "cloudless" approach, FogFlow can not only provide fast response time, but also achieve high scalability and reliability.

To define and trigger FogFlow based instances follow `Intent based programming model`_.

.. _`Intent based programming model`: https://fogflow.readthedocs.io/en/latest/intent_based_program.html

Core concepts
======================

Operator
----------------------

In FogFlow an operator presents a type of data processing unit, 
which receives certain input streams as NGSI10 notify messages via a listening port,
processes the received data, generates certain results, and publishes the generated results as NGSI10 updates.   

The implementation of an operator is associated with at least one docker images. 
To support various hardware architectures (e.g., X86 and ARM for 64bits or 32 bits), 
the same operator can be associated with multiple docker images.  

Task
------------------

A task is a data structure to represent a logic data processing unit within a service topology. 
Each task is associated with an operator. 
A task is defined with the following properties:

- name: a unique name to present this task
- operator: the name of its associated operator
- groupBy: the granularity to control the unit of its task instances, which is used by service orchestrator to determine how many task instances must be created
- input_streams: the list of its selected input streams, each of which is identified by an entity type
- output_streams: the list of its generated output streams, each of which is identified by an entity type

In FogFlow, each input/output stream is represented as a type of NGSI context entities, 
which are usually generated and updated by either an endpoint device or a data processing task. 

During the runtime, multiple task instances can be created for the same task, 
according to its granularity defined by the *groupBy* property. 
In order to determine which input stream goes to which task instances, 
the following two properties are introduced to specify the input streams of tasks: 

- Shuffling: associated with each type of input stream for a task; its value can be either *broadcast* or *unicast*. 	

	- broadcast: the selected input streams should be repeatedly assigned to every task instance of this operator
	- unicast: each of the selected input streams should be assigned to a specific task instance only once
	
- Scoped: determines whether the geo-scope in the requirement should be applied to select the input streams; its value can be either *true* or *false*.


Task Instance
-----------------------

During the runtime, a task is configured by FogFlow with its input data and specified output type 
and then the configured task will be launched as a task instance, running in a docker container. 
Currently, each task instance is deployed in a dedicated docker container, either in the cloud or at an edge node. 


Service Template
-------------------------

Each IoT service is described by a service template, which can be a service topology with a set of linked operators
or a fog function with a single operator. For example, when a service topology is used to specify the service template, 
the following information will be included. 

- topology name: the unique name of the topology
- service description: some text to describe what this service is about
- priority: define the priority level of all tasks in the topology, which will be utilized by edge nodes to decide how resource should be assigned to tasks 
- resource usage: define if the tasks in this topology can use the resources on edge nodes in an exclusive way, meaning that not sharing resources with any task from the other topologies

Currently, FogFlow provides a graphical editor to allow developers to easily define and annotate their service topology or fog function during the design phrase.


Dynamic data flow 
-----------------------

On receiving a requirement, Topology Master creates a dataflow execution graph and then deploys them over the cloud and edges. 
The main procedure is illustrated by the following figure, including two major steps. 

.. figure:: figures/service-topology.png

- from *service topology* to *execution plan*: done by the task generation algorithm of Topology Master. 
	The generated execution plan includes:
	1) which part of service topology is triggered; 
	2) how many instances need to be created for each triggered task;
	3) and how each task instance should be configured with its input streams and output streams. 

- from *execution plan* to *deployment plan*: done by the task assignment algorithm of Topology Master.
	The generated deployment plan determines which task instance should be assigned to which worker (in the cloud or at edges),  
	according to certain optimization objectives. Currently, the task assignment in FogFlow is optimized to reduce across-node data traffic
	without overloading any edge node. 



Intent Model
==================

An intent object basically consists of these properties as illustrated in below figure.

- Serivce topology specifies the computation logic for which the intent object is triggered. 

- Geoscope is defined a geographical location where input streams should be selected. Geoscope cab be setect as global value as well as can be set custom geoscopes. 

- Service Level Object (SLO) is the objective of maximum throughput, minimum latency and minimum cost can be set for task assignment at workers. However, this feature is not fully supported yet, so User can ignore this. It can be set as “None” for now. 

- Resource Usage defines how a topology can use resources on edge nodes. It can either exlusive or inclusive. In an exclusive way means the topology will not share the resources with any task from other topologies. Whereas in an inclusive one topology will share the resources with any task from other topologies.


.. figure:: figures/intent_model.png




Three Key Elements
=====================
The three key elements to program an IoT service is illustrated via below figure.




.. figure:: figures/key_elements.png




In FogFlow several operators form a graph which is defined as a service topology. Each operator in the service topology is annotated with its inputs and outputs, which indicate their dependency to the other tasks in the same topology. Service topology can easily compose different operators to form their service logic in just a few minutes. After that, during the runtime data processing flows can be automatically triggerred based on the high level data usage intent defined by service users. Service users can be either data producers or result consumers.




Service topology example of "Hello-World"
=========================================

Register task operators
----------------------------------------------

**Register an Operator** to define what would be the name of Operator and what input parameters it would need. Here in this context, an operator is nothing but a named element having some parameters.
The following picture shows the list of all registered operators and their parameter count.


.. figure:: figures/operator_dashboard.png



After clicking the "register" button, a design area can be seen below and an operator can be created and parameters can be added to it. Add name of operator "Hello-World" and description for operator is optional, for parameter name can be anything here it is given as hello and value is set as default.



.. figure:: figures/register_operator3.png




**Register a Docker Image and choose Operator** to define the docker image and associate an already registered Operator with it. 

The following picture shows the list of all registered docker images and the key information of each image. 


.. figure:: figures/docker_image1.png


After clicking the "register" button, a form can be seen as below. 
Please fill out the required information and click the "register" button to finish the registration. 
The form is explained as the following. 

* Image: the name of your operator docker image
* Tag: the tag is used to publish the operator docker image; by default it is "latest"
* Hardware Type: the hardware type that the docker image supports, including X86 or ARM (e.g. Raspberry Pi)
* OS Type: the operating system type that the docker image supports; currently this is only limited to Linux
* Operator: the operator name, which must be unique and will be used when defining a service topology
* Prefetched: if this is checked, that means all edge nodes will start to fetch this docker image in advance; otherwise, the operator docker image is fetched on demand, only when edge nodes need to run a scheduled task associated with this operator. 


.. important::
    
    Please notice that the name of the docker image must be consistent with the one that is published to `Docker Hub`_.
    By default, FogFlow will fetch the required docker images from Docker Hub using the name that is registered here for an operator. 
    
.. _`Docker Hub`: https://github.com/smartfog/fogflow/tree/master/application/operator/anomaly
    
    

.. figure:: figures/docker_image.png




Define and trigger a service topology for "Hello-World"
------------------------------------------------------------

Define service topology using FogFlow Topology Editor.



.. figure:: figures/topology_register.png


As seen in the picture, the following important information must be provided. 

#. define topology profile, including
    * topology name: the unique name of your topology
    * service description: some text to describe what this service is about
    
#. draw the graph of data processing flows within the service topology
    With a right click at some place of the design board, you will see a menu pops up 
    and then you can start to choose either task or input streams or shuffle
    to define the data processing flow.
   
   
The following information is required to specify a task profile.
	
    * name: the name of the task 
    * operator: the name of the operator that implements the data processing logic of this task; an operator shall be registered beforehand so that it can be shown from the list
    * entity type of output streams: to specify the entity type of the produced output stream.
    
    The following information is required to specify an EntityStream Profile.

    * SelectedType: is used to define what Entity Type will be chosen by the task as its Input Stream
    * SelectedAttributes: is used to define what attribute (or attributes) of the Selected Entity Type will be considered for changing the state of a task.
    * Groupby: to determine how many instances of this task should be created on the fly; currently including the following cases
	
        *  if there is only one instance to be created for this task, please use "groupby" = "all"
        *  if it is needed to create one instance for each entity ID of the input streams, please use "groupby" = "entityID"
        *  if it is neded to create one instance for each unique value of some specific context metadata, please use the name of this registered context metadata
    
    
* Scoped: tells if the Entity data are location-specific or not. True indicates that location-specific data are recorded in the Entity and False is used in case of broadcasted data, for example, some rule or threshold data that holds true for all locations, not for a specific location.

    Shuffling element serves as a connector between two tasks such that output of a task is the input for the shuffle element and same is forwarded by Shuffle to another task (or tasks) as input.
    
    
    
Trigger the service topology by sending an Intent
------------------------------------------------------------------------------

Once developers submit a specified service topology and the implemented operators, 
the service data processing logic can be triggered by following two steps:

* Sending a high level intent object which breaks the service topology into separate tasks
* Providing Input Streams to the tasks of that service topology.

The intent object is sent using the fogflow dashboard with the following properties: 

* Topology: specifies which topology the intent object is meant for.
* Priority: defines the priority level of all tasks in your topology, which will be utilized by edge nodes to decide how resources should be assigned to the tasks.
* Resource Usage: defines how a topology can use resources on edge nodes. Sharing in an exclusive way means the topology will not share the resources with any task from other topologies. The other way is inclusive one.
* Objective: of maximum throughput, minimum latency and minimum cost can be set for task assignment at workers. However, this feature is not fully supported yet, so it can be set as "None" for now.
* Geoscope: is a defined geographical area where input streams should be selected. Global as well as custom geoscopes can be set.



.. figure:: figures/intent.png




Fogflow topology master will now be waiting for input streams for the tasks contained in the service topology. As soon as context data are received, which fall within the scope of the intent object, tasks are launched on the nearest workers.









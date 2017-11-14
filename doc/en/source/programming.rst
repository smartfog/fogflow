Core concepts and programming model
======================================

The FogFlow programming model defines the way of how to specify a *service topology* 
using *declarative hints* and how to implement tasks based on NGSI. 
First, developers decompose an IoT service into multiple tasks and then define its service topology 
as a DAG in JSON format to express the data dependencies between different tasks. 
On the other hand, the FogFlow programming model provides declarative hints for developers 
to guide service orchestration without introducing much complexity.

The following concepts are required to understand the FogFlow programming model. 


Operator
------------------

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

* name: a unique name to present this task
* operator: the name of its associated operator
* groupBy: the granularity to control the unit of its task instances, which is used by service orchestrator to determine how many task instances must be created
* input_streams: the list of its selected input streams, each of which is identified by an entity type
* output_streams: the list of its generated output streams, each of which is identified by an entity type

In FogFlow, each input/output stream is represented as a type of NGSI context entities, 
which are usually generated and updated by either an endpoint device or a data processing task. 

During the runtime, multiple task instances can be created for the same task, 
according to its granularity defined by the *groupBy" property. 
In order to determine which input stream goes to which task instances, 
the following two properties are introduced to specify the input streams of tasks: 

* Shuffling: associated with each type of input stream for a task; its value can be either *broadcast* or *unicast*. 

	- broadcast: the selected input streams should be repeatedly assigned to every task instance of this operator
	- unicast: each of the selected input streams should be assigned to a specific task instance only once
	
* Scoped: determines whether the geo-scope in the requirement should be applied to select the input streams; its value can be either *true* or *false*.


Topology
------------------

A service topology presents the decomposited data processing logic of a service application. 
Decomposting the entire service application logic into small data processing tasks
allows FogFlow to dynamically and seamlessly migrate data processing tasks between the cloud and edges. 
It also enables better sharing of intermediate results. 

Each service topology is represented as a directed acyclic graph (DAG) of multiple linked tasks, 
with the following properties: 

* description: a text to describe what this service topology is about
* name: a unique name of this service topology
* priority: the task priority of this service topology 
* tasks: the list of its tasks


Requirement
------------------

Once developers submit a specified service topology and the implemented operators, 
the service data processing logic can be triggered on demand by a high level processing requirement. 
The processing requirement is sent as NGSI10 update, with the following properties: 

* topology: which topology to trigger
* expected output: the output stream type expected by external subscribers
* scope: a defined geoscope for the area where input streams should be selected
* scheduler: which type of scheduling method should be chosen by Topology Master for task assignment


Dynamic dataflow execution over cloud and edges
------------------------------------------------

On receiving a requirement, Topology Master creates a dataflow execution graph and then deploys them over the cloud and edges. 
The main procedure is illutrated by the following figure, including two major steps. 

.. figure:: figures/service-topology.png
   :scale: 100 %
   :alt: map to buried treasure

* from *service topology* to *execution plan*: done by the task generation algorithm of Topology Master. 
The generated execution plan includes:1) which part of service topology is triggered; 
2) how many instances need to be created for each triggered task;
3) and how each task instance should be configured with its intput streams and output streams. 

* from *execution plan* to *deployment plan*: done by the task assignment algorithm of Topology Master.
The generated deployment plan determines which task instance should be assgined to which worker (in the cloud or at edges),  
according to certain optimization objectives. Currently, the task assignment in FogFlow is optimized to reduce across-node data traffic
without overloading any edge node. 



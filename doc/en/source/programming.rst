Edge programming models
======================================

The FogFlow programming model defines the way of how to specify a *service topology* 
using *declarative hints* and how to implement tasks based on NGSI. 
First, developers decompose an IoT service into multiple tasks and then define its service topology 
as a DAG in JSON format to express the data dependencies between different tasks. 
On the other hand, the FogFlow programming model provides declarative hints for developers 
to guide service orchestration without introducing much complexity.

The following concepts are required to understand the FogFlow programming model. 


Service Topology
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

Once developers submit a specified service topology and the implemented operators, 
the service data processing logic can be triggered on demand by a high level processing requirement. 
The processing requirement is sent as NGSI10 update, with the following properties: 

* topology: which topology to trigger
* expected output: the output stream type expected by external subscribers
* scope: a defined geoscope for the area where input streams should be selected
* scheduler: which type of scheduling method should be chosen by Topology Master for task assignment









Fog Function
------------------




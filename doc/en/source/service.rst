Dynamic service orchestration
==============================

For a given processing requirement, Topology Master performs the following steps to dynamically orchestrate tasks over cloud and edges. 

.. figure:: figures/service-orchestration.jpg
   :scale: 100 %
   :alt: map to buried treasure


Topology Lookup
---------------- 

Iterating over the requested service topology to find out the processing tree in order to produce the expected output. This extracted processing tree represents the requested processing topology which is further used for task generation. 

Task generation
----------------

First querying IoT Discovery to discover all available input streams and then deriving an execution plan based on this discovery and the declarative hints in the service topology. 
The execution plan includes all generated tasks that are properly configured with right input and output streams and also the parameters for the workers to instantiate the tasks.

Task Deployment
----------------

Performing the specified scheduling method to assign the generated tasks to geo-distributed workers according to their available computation capabilities. The derived assignment result represents the deployment plan. To carry out the deployment plan, TM sends each task to the task's assigned worker and then monitors the status of the task. Each worker receives its assigned tasks and then instantiates them in docker containers. Meanwhile, worker communicates with the nearby IoT Broker to assist the launched task instances for establishing their input and output streams. 



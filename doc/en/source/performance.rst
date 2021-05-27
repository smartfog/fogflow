*****************************************
Performance Evaluation
*****************************************


Performance of NGSI-LD based thinBroker in FogFlow 
================================================================


Experiement Setup
-------------------

to explain the system environment for the performance evaluation, such as how FogFlow is deployed and where is the test client,
which tool to generate the test workload


Throughput and latency to create new entities
--------------------------------------------------

to evaluate how fast a new NGSI-LD entity can be created

to compare the performance with the other NGSI-LD brokers

to test how the performance can be scaled up with more FogFlow edge nodes



Throughput and latency to query entities
--------------------------------------------------

to prepare different types of queries: by entity ID, by entity type, by the filter of entity attribute

to compare the performance with the other NGSI-LD brokers

to test how the performance can be scaled up with more FogFlow edge nodes


Update Propagation from Context Producers to Context Consumer
------------------------------------------------------------------

to measure the delay of context update from the moment sent by a context producer to the time received by a subscriber

to measure how many updates can flow from the context producer to the subscriber per second

to compare the performance with the other NGSI-LD brokers

to test how the performance can be scaled up with more subscribers



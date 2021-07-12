*****************************************
Performance Evaluation
*****************************************


Performance of NGSI-LD based thinBroker in FogFlow 
================================================================

This tutorial introduces the performace evaluation of NGSILD based context management systems of the FogFlow framework. Our analyses include the performance comparision of FogFlow broker with other NGSILD broker(orion broker, stelio Broker, Scorpio Broker) in terms of  throughput (number of messages per second) and response time/message propagation latency and  efficiency of context availability discoveries and context transfers in the smart city scale.Moreover, we analyze the scalability of FogFlow using multiple IoT Brokers.


Experiement Setup
-------------------

**FogFlow system:** To test the performance, I have deployed one cloud node(2vCPUs, 4 GB RAM) and two edge node(2vCPUs, 4 GB RAM) in doker container.

**Listener client:** To measure the delay of context update from the moment sent by a context producer to the time received by a subscriber we are using listener client. Listener client is deployed on localhost(5cpu, 8GB RAM)

**Testing tool:** TO produce the test data for fogflow broker we are using Apache JMeter testing tool. JMeter is deployed on localhost(5cpu, 8GB RAM)



Throughput and latency to create new entities
--------------------------------------------------

.. figure:: figures/1.1.png

.. figure:: figures/1.1NewData.png

The above graph depicts the variation of latency over number of threads (the real time users). The Y axis represent latency and X axis represent number of threads. On analysing above data, it becomes evident that with increasing number of thread the total number of requests increases. The throughput in contrast to the increasing number of request indicates the good performance of fogflow for Upsert requests.

The above graph depicts the variation of latency over number of threads (the real time users) in case of subscription. The Y axis represent latency and X axis represent number of threads. On analysing above data, it becomes evident that with increasing number of thread the total number of requests increases. The throughput in contrast to the increasing number of request indicates the average performance of fogflow for Subscription requests. Given the fact that fogflow subscriptions are interacting with fogflow component like fogflow discovery making it reliable but adding an extra tint of time in generating response.

.. figure:: figures/1.2NewUpsert.png

.. figure:: figures/1.2upsertdata.png

.. figure:: figures/1.2Subscription.png

.. figure:: figures/1.2SubscriptionData.png

.. figure:: figures/1.3upsert.png

.. figure:: figures/1.3upsertdata.png


Throughput and latency to query entities
--------------------------------------------------

.. figure:: figures/2.1Id.png

.. figure:: figures/2.1IDData.png

.. figure:: figures/2.1SubID.png

.. figure:: figures/2.1SubBYIDData.png


Update Propagation from Context Producers to Context Consumer
------------------------------------------------------------------

to measure the delay of context update from the moment sent by a context producer to the time received by a subscriber

to measure how many updates can flow from the context producer to the subscriber per second

to compare the performance with the other NGSI-LD brokers

to test how the performance can be scaled up with more subscribers



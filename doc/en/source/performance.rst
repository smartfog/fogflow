*****************************************
Performance Evaluation
*****************************************


Performance of NGSI-LD based thinBroker in FogFlow 
================================================================

This tutorial introduces the performace evaluation of NGSILD based context management systems of the FogFlow framework. Our analyses include the performance comparision of FogFlow broker with other NGSILD broker(orion broker, stelio Broker, Scorpio Broker) in terms of  throughput (number of messages per second) and response time/message propagation latency and  efficiency of context availability discoveries and context transfers in the smart city scale.Moreover, we analyze the scalability of FogFlow using multiple IoT Brokers.


Experiement Setup
-------------------

**FogFlow system:** To test the performance, I have deployed one cloud node(4vCPUs, 16 GB RAM) and two edge node(4vCPUs, 8 GB RAM) in doker container.

**Listener client:** To measure the delay of context update from the moment sent by a context producer to the time received by a subscriber we are using listener client. Listener client is deployed on a VM(4cpu, 8GB RAM)

**Testing tool:** TO produce the test data for fogflow broker we are using Apache JMeter testing tool. JMeter is deployed on a VM(8cpu, 8GB RAM)



Throughput and latency to create new entities
--------------------------------------------------

**To test the performance of Subscription and Upsert API in Fogflow** 

.. figure:: figures/1.1.png

.. figure:: figures/1.1NewData.png

**Fogflow Upsert API**

The above graph depicts the variation of latency over number of threads (the real time users). The Y axis represent latency and X axis represent number of threads. On analysing above data, it becomes evident that with increasing number of thread the total number of requests increases. The throughput in contrast to the increasing number of request indicates the good performance of fogflow for Upsert requests.

**Fogflow Subscription API**

The above graph depicts the variation of latency over number of threads (the real time users) in case of subscription. The Y axis represent latency and X axis represent number of threads. On analysing above data, it becomes evident that with increasing number of thread the total number of requests increases. The throughput in contrast to the increasing number of request indicates the average performance of fogflow for Subscription requests. Given the fact that fogflow subscriptions are interacting with fogflow component like fogflow discovery making it reliable but adding an extra tint of time in generating response.

Performance Comparison between Fogflow and Scorpio Broker
--------------------------------------------------------------

To compare response time of Fogflow upsert API with Scorpio Broker upsert API, we have created 36500 entities by using different no of thread 10, 50,100, 200, 400,500. The following graph repersent response time on Y-axis and no of thread on X-axis. 

**Fogflow Upsert API Vs. Scorpio Upsert API**

.. figure:: figures/1.2NewUpsert.png

.. figure:: figures/1.2upsertdata.png

**Comparison Result** : The above graph is combination of two graphs i.e. the blue marker represents Scorpio broker and orange marker represents Fogflow. With a detailed analysis of the response-time and number of thread graph, it is visible that Fogflow broker's Upsert API is a better performer than Scorpio broker's Upsert API. As shown in tabular data, it is evident that on icreasing the number of threads which utlimately increases number of requests are better handled in case of Fogflow. For example on executing 500 threads with a sum total of 7500 requests, Fogflow has an average throughput of 1156.69/s whereas Scorpio broker on same number of thread has an average throughput of 222.73/s.

*Hence Fogflow Upsert API is better performer than Scorpio Broker*

**Fogflow Subscription API Vs. Scorpio Subscription API**

.. figure:: figures/1.2Subscription.png

.. figure:: figures/1.2SubscriptionData.png

**Comparison Result** : The above graph is combination of two graphs i.e. the blue marker represents Scorpio broker and orange marker represents Fogflow. With a detailed analysis of the response-time and number of thread graph, it is visible that Fogflow broker's Subscription API is a better performer than Scorpio broker's Subscription API initially but gradually the graph of Fogflow rises which reflects the small increase in response time. Given the fact that fogflow subscriptions are interacting with fogflow component like fogflow discovery making it reliable but adding an extra tint of time in generating response gives it an edge over other brokers but trade's off quality over extra tint of response time. The reliability here is a key factor because of its distributed architecture that involves different edge nodes which is missing in scorpio btoker.

**To test how the performance can be scaled up with more FogFlow edge nodes**

.. figure:: figures/1.3upsert.png

.. figure:: figures/1.3upsertdata.png

**Comparison Result** : The above graph is combination of two graphs i.e. the blue marker represents Fogflow Edge node and orange marker represents Fogflow cloud node. With a detailed analysis of the response-time and number of thread graph, it is visible that Fogflow cloud node has a higher response time than Fogflow edge node because of the fact that the increased number of edge brokers speed up the process because they all are interally connected to discovery. The requests made to edge node are registered with discovery directly than having to follow up a longer path through cloud broker. Thus, the Upsert API has an increased throughput on same number of thread as for cloud. As shown in tabular data, the cloud achieves a throughput of 991.98/s for 100000 requests with 500 number of threads and edge node achieves 1743.3/s for same 100000 requests for 500 number of threads.

Throughput and latency to query entities
--------------------------------------------------

**To compare the performance of Query APIs with the other NGSI-LD brokers**

.. figure:: figures/2.1Id.png

.. figure:: figures/2.1IDData.png

.. figure:: figures/2.1SubID.png

.. figure:: figures/2.1SubBYIDData.png

**To test how the performance of Query APIs can be scaled up with more FogFlow edge nodes**

.. figure:: figures/2.3.png

.. figure:: figures/2.3QueryCloudEdge.png

Update Propagation from Context Producers to Context Consumer
------------------------------------------------------------------

to measure the delay of context update from the moment sent by a context producer to the time received by a subscriber

to measure how many updates can flow from the context producer to the subscriber per second

to compare the performance with the other NGSI-LD brokers

to test how the performance can be scaled up with more subscribers



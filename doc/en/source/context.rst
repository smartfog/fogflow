Distributed context management
===============================

The context management system is designed to provide a global view for all system components and running task instances 
to query, subscribe, and update context entities via the unified data model and the communication protocol *NGSI*. 
It plays a very important role to support the standard-based edge programming model in FogFlow. 
As compared to other existing brokers like MQTT-based Mosquitto or Apache Kafka, 
the distributed context management system in FogFlow has the following features: 

* separating context availability and context entity
* providing separated and standardized interfaces to manage both context data (via NGSI10) and context availability (via NGSI9). 
* supporting not only ID-based and topic-based query and subscription but also geoscope-based query and subscription

As illustrated by the following figure, in FogFlow a large number of distributed IoT Brokers work in parallel
under the coordination a global and centralized IoT Discovery. 

.. figure:: figures/distributed-brokers.png
   :scale: 100 %
   :alt: map to buried treasure


IoT Discovery
-------------
The centralized IoT Discovery provides a global view of context availability of context data and provides NGSI9 interfaces for registration, discovery, and subscription of context availability. 

IoT Broker
-------------
The IoT Broker in Fogflow is very light-weight, because it keeps only the lastest value of each context entity
and saves each entity data directly in the system memory. 
This brings high throughput and low latency for the data transfer from context produers to context consumers. 

Each IoT Broker manages a portion of the context data and registers data to the shared IoT Discovery.
However, all IoT Brokers can equally provide any requested context entity via NGSI10 
because they can find out which IoT Broker provides the entity through the shared IoT Discovery and then fetch the entity from that remote IoT Broker. 







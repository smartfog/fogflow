***************
Roadmap
***************

The FogFlow Generic Enabler is a distributed execution framework to support dynamic processing flows over cloud and edges.

Introduction
=================

This section elaborates on proposed new features or tasks which are expected to be added to the product in the foreseeable future. There should be no assumption of a commitment to deliver these features on specific dates or in the order given. The development team will be doing their best to follow the proposed dates and priorities, but please bear in mind that plans to work on a given feature or task may be revised. All information is provided as general guidelines only, and this section may be revised to provide newer information at any time.

Disclaimer:

- This section has been last updated in February 2020. Please take into account its content could be obsolete.
- Note we develop this software in Agile way, so development plan is continuously under review. Thus, this roadmap has to be understood as rough plan of features to be done along time which is fully valid only at the time of writing it. This roadmap has not be understood as a commitment on features and/or dates.
- Some of the roadmap items may be implemented by external community developers, out of the scope of GE owners. Thus, the moment in which these features will be finalized cannot be assured.
  
Short Term
---------------

The following list of features are planned to be addressed in the short term, and incorporated in a next release of the product:

1. FogFlow NGSIv2 support and integration with Wire cloud and Quantum Leap GE’s.
  - Since most of the other FIWARE GE now only support v2, it will be nice for FogFlow to support v2 as well. In this case, it will be easy to integrate FogFlow with IoTAgent for device integration, with Wire Cloud for visualization, and with Quantum Leap for saving historical data.
  - Functional Test cases automation for NGSIv2 support and integrate with FogFlow build.
2. Fogflow System Monitoring :- It will be quite useful to deploy Prometheus to monitor the resource and tasks in the entire system.
  - Show the monitoring information of each edge node
  - Be able to add or remove edge node dynamically
3. User Manual Updation
  - User manual updation for integration with other Fiware GE components.
  - FogFlow user manual support for new feature and bugs.
  
Medium Term
-------------------

The following list of features are planned to be addressed in the medium term, typically within the subsequent release(s) generated in the next 6 months after the next planned release.

1. Fogflow Persistent storage
  - Fogflow Persistent storage of all defined operators, service topologies, and fog functions.
2. NGSI-LD support
  - NGSI-LD support in Fogflow and integration with Scorpio Broker.
3. Firewall Support
  - Currently, the FogFlow edge node requires to have a public IP address to be accessible by the FogFlow cloud node. In the actual deployment environment, the FogFlow edge node is very often deployed behind the company firewall via NAT, we need to find a way to support this scenario. One way to address this is to find a proxy for such kind of edge nodes. For example, assign the FogFlow cloud broker to be the proxy for the FogFlow brokers at this type of edge nodes.
4. Edge AI
  - Support the edge node with Edge TPU.
  
Long term
-----------------

The following list of features are proposals regarding the longer-term evolution of the product even though the development of these features has not yet been scheduled for a release in the near future. Please feel free to contact us if you wish to get involved in the implementation or influence the roadmap:

1. Multi-tenancy support.
  - Support multiple users over the same cloud-edge infrastructure.
2. Digital twin support
  - make the current programming model to support the creation of digital twins and also the interaction between digital twins.
3. Semantics-based data integration
  - creating dynamic data processing pipelines to convert arbitrary raw data into standard-based entities.
4. semantics-based service composition
  - linking serverless functions based on their semantically-annotated inputs and outputs.

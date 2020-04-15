*****************************
System Overview
*****************************

FogFlow is a distributed execution framework to support dynamic processing flows over cloud and edges. It can dynamically and 
automatically composite multiple NGSI-based data processing tasks to form high level IoT services and then orchestrate and optimize 
the deployment of those services within a shared cloud-edge environment.

The shared environment of FogFlow cloud-edge can be created with one FogFlow cloud node and more than one FogFlow edge nodes as
illustrated in below figure. All integrated features that are running in FogFlow system, can be seen in this figure. 



.. figure:: figures/FogFlow_System_Design.png


In this page, a brief introduction is given about FogFlow integrations, more detail can be found in `Integration`_, section.

..`Integration`_: https://fogflow.readthedocs.io/en/latest/integration.html#integration

There are mainly two types of integration Northbound and Southbound, flow of data from a sensor device towards broker is known 
as Northbound Flow and when flow of data from broker towards actuator devices, then it is known as Southbound Flow.
more detail about Northbound and Southbound data flow can be checked via `Integrate an actuator device with Fogflow`_, page.


..`Integrate an actuator device with Fogflow`_: https://fogflow.readthedocs.io/en/latest/example5.html#integrate-an-actuator-device-with-fogflow


FogFlow Integration with Scorpio broker, Scorpio is an NGSI-LD compliant context broker, an NGSI-LD Adapter is built 
to enable FogFlow Ecosystem to talk with Scorpio context broker. The NGSI-LD Adapter converts NGSI data format to NGSI-LD and forward it
to Scorpio broker, more detail can be checked via `Integrate FogFlow with Scorpio Broker`_, page.


..`Integrate FogFlow with Scorpio Broker`_: https://fogflow.readthedocs.io/en/latest/scorpioIntegration.html#integrate-fogflow-with-scorpio-broker


Integration with Orion broker, FogFlow can be intergrated with Orion context broker using NGSI-V1 as well as NGSI-V2 APIs.
more detail can be checked via `Integrate FogFlow with FIWARE`_, page.

..`Integrate FogFlow with FIWARE`_: https://fogflow.readthedocs.io/en/latest/example4.html#integrate-fogflow-with-fiware

Similarly, FogFlow Integration with WireCloud is provided to visualize the data with the help of different widgets of wirecloud
and Integration with QuantumLeap is to store time series based historical data. more detail can be checked via  `Integrate FogFlow with WireCloud`_,
for wirecloud and `Integrate FogFlow with QuantumLeap`_, page for QuantumLeap.

..`Integrate FogFlow with WireCloud`_: https://fogflow.readthedocs.io/en/latest/wirecloudIntegration.html#integrate-fogflow-with-wirecloud
..`Integrate FogFlow with QuantumLeap`_: https://fogflow.readthedocs.io/en/latest/quantumleapIntegration.html#integrate-fogflow-with-quantumleap



FogFlow also provides a secure communication between the FogFlow cloud node and the FogFlow edge nodes, or between two edge nodes.
To acheive  HTTPs-based communication secure communication in FogFlow, it is necessary for FogFlow cloud node and the FogFlow edge
node to have their own domain names. Further the detail configuration and setup steps can be checked via `Security`_,.

..`Security`_: https://fogflow.readthedocs.io/en/latest/https.html#security



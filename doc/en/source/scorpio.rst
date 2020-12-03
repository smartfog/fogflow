Using NGSI-LD speciifcation implementation 
===============================================
Scorpio integration with FogFLow enable FogFlow task to communicate with scorpio Broker.
The figure below shows how data will transmit between socrpio broker, FOgFLow broker and FogFlow task.

.. figure:: figures/scorpioIntegration.png

Integration steps
===============================================
**Pre-Requisites:**

* Fogflow should be up and running with atleast one node.
* Scorpio Broker should be up and running.
**Implemantation Steps:**

* Create and trigger NGSI-LD task (`See Document`_).
.. _`See Document`: https://fogflow.readthedocs.io/en/latest/intent_based_program.html.
* **Send NGSI-LD subscription request to Scorpio Broker** to get notification form Scorpio Broker for every update on Scoprpio broker.


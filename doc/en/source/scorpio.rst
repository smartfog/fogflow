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
* Create and trigger NGSI-LD task (`See Document`_).
.. _`See Document`: https://fogflow.readthedocs.io/en/latest/intent_based_program.html.

**There are two type of Itegration**

* when a NGSI-LD device will send some update on scorpio broker then update should be notified to the FogFLow broker . After getting the notification from scorpio broker FogFlow will send this notification to the FogFLow task for furthur analysis. For This integration FogFLow broker will subscribe to the scorpio broker using following request.

.. code-block:: console

    curl -iX POST \
    'http://<Scorpio Broker>/ngsi-ld/v1/subscriptions/' \
      -H 'Content-Type: application/ld+json' \
      -H 'Accept: application/ld+json' \
      -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
      -d '
      {
         "type": "Subscription",
         "entities": [{
                "id" : "urn:ngsi-ld:Vehicle:A13",
                "type": "Vehicle"
           }],
          "watchedAttributes": ["*"],
          "notification": {
                 "attributes": ["*"],
                  "format": "keyValues",
                 "endpoint": {
                        "uri": "http://<FogFLow Broker>/ngsi-ld/v1/notifyContext/",
                        "accept": "application/json"
                }
         }
    }'



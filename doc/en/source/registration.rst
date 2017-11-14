.. _operator-registration:

How to register an operator docker image
=========================================

For each operator, once we create its docker image and push it to the FogFlow docker registry, 
we must register the operator in FogFlow. 
This can be done in one of the following two ways. 


.. note:: Please notice that each operator must have a unique name but the same operator can be associated with multiple docker images, 
			each of which is for one specific hardware or operating system but for implementing the same data processing logic. 
			During the runtime, FogFlow will select a proper docker image to run a scheduled task on an edge node, 
			based on the executuion environment of the edge node. 


Register it via FogFlow Task Designer
--------------------------------------------------------------

The following picture shows the list of all registered operator docker images and the key imformation of each image. 

.. figure:: figures/operator-registry-list.png
   :scale: 100 %
   :alt: map to buried treasure


After clicking the "register" button, you can see a form as below. 
Please fill out the required information and click the "register" button to finish the registration. 
The form is explained as the following. 

* Image: the name of your operator docker image
* Tag: the tag you used to publish your operator docker image; by default it is "latest"
* Hardware Type: the hardware type that your docker image supports, including X86 or ARM (e.g. Raspberry Pi)
* OS Type: the operating system type that your docker image supports; currently this is only limited to Linux
* Operator: the operator name, which must be unique and will be used when defining a service topology
* Prefetched: if this is checked, that means all edge nodes will start to fetch this docker image in advance; otherwise, the operator docker image is fetched on demand, only when edge nodes need to run a scheduled task associated with this operator. 

.. figure:: figures/operator-register.png
   :scale: 100 %
   :alt: map to buried treasure


Register it programmatically by sending a NGSI update 
-----------------------------------------------------------

You can also register an operator docker image by sending a constructed NGSI update message to the IoT Broker deployed in the cloud. 

Here is a Javascript-based code example to register an operator docker image. 
Within this code example, we use the Javascript-based library to interact with FogFlow IoT Broker. 
You can find out the library from the github code repository (designer/public/lib/ngsi). You must include ngsiclient.js into your web page. 

.. code-block:: javascript

    var image = {
        name: "counter",
        tag: "latest",
        hwType: "X86",
        osType: "Linux",
        operatorName: "counter",
        prefetched: false
    };

    //register a new docker image
    var newImageObject = {};

    newImageObject.entityId = {
        id : image.name + ':' + image.tag, 
        type: 'DockerImage',
        isPattern: false
    };

    newImageObject.attributes = {};   
    newImageObject.attributes.image = {type: 'string', value: image.name};        
    newImageObject.attributes.tag = {type: 'string', value: image.tag};    
    newImageObject.attributes.hwType = {type: 'string', value: image.hwType};      
    newImageObject.attributes.osType = {type: 'string', value: image.osType};          
    newImageObject.attributes.operator = {type: 'string', value: image.operatorName};      
    newImageObject.attributes.prefetched = {type: 'boolean', value: image.prefetched};                      
    
    newImageObject.metadata = {};    
    newImageObject.metadata.operator = {
        type: 'string',
        value: image.operatorName
    };               
	
	// assume the config.brokerURL is the IP of cloud IoT Broker
    var client = new NGSI10Client(config.brokerURL);	
    client.updateContext(newImageObject).then( function(data) {
        console.log(data);
    }).catch( function(error) {
        console.log('failed to register the new device object');
    });    	

.. _operator-implementation:

How to implement a dockerized operator
========================================

As explained in :ref:`flow-task`, FogFlow can automatically subscribe the required input streams
on behalf of data processing tasks. 
What operator developers should focus on is the internal data processing logic. 

Available templates
-----------------------

Currently, FogFlow provides the templates to implement operators for JavaScript (nodejs-based) and Python. 
You can reuse them when developing your own operators. 

	* `nodejs template <https://github.com/smartfog/fogflow/tree/master/application/operator/template/javascript/>`_.
	* `python template <https://github.com/smartfog/fogflow/tree/master/application/operator/template/python/>`_. 

In the following, we explain how to implement an operator based on the provided nodejs template. 

How to program an operator
-------------------------------

Once you download the nodejs template, you can follow the framework to implement your own application logic. 
Overall, you need to know the following important things: 

* what is the main logic part of the provided tempalte?

	The following code block shows the main logic part of our provided nodejs template.
	It uses two libraries: ngsiclient.js and ngsiagent.js. 
	ngsiclient.js is a nodejs-based based NGSI client to interact with FogFlow broker via NGSI10. 
	ngsiagent.js is a server to receive NGSI10 notify. 
	With the provided libraries, your application can perform NGSI10 update, subscribe, query. 
	
	In the main logic, three important functions are called. The responsibility of each function is explained as below. 
	
	* startApp: to do some initialization of your application; this function will be called right after ngsiagent starts listening to receive input data
	* handleAdmin: to handle the configuration commands issued by FogFlow
	* handleNotify: to handle all received notify messages; please note that the subscription of your input data is done by FogFlow on behalf of your application
	
.. code-block:: javascript

    const NGSIClient = require('./ngsi/ngsiclient.js');
    const NGSIAgent = require('./ngsi/ngsiagent.js');
    
    // get the listening port number from the environment variables given by the FogFlow edge worker
    var myport = process.env.myport;
    
    // set up the NGSI agent to listen on 
    NGSIAgent.setNotifyHandler(handleNotify);
    NGSIAgent.setAdminHandler(handleAdmin);
    NGSIAgent.start(myport, startApp);
    
    process.on('SIGINT', function() {
        NGSIAgent.stop();	
        stopApp();
        
        process.exit(0);
    });
    
    
    function startApp() 
    {
        console.log('your application starts listening on port ',  myport);
        
        // to add the initialization part of your application logic here    
    }
    
    function stopApp() 
    {
        console.log('clean up the app');
    }
	
	
	
* how to initialize something when your applicaiton starts?

.. code-block:: javascript

	function startApp() 
	{
	    console.log('your application starts listening on port ',  myport);
    
	    // to add the initialization part of your application logic here    
	}

* how is your application configured by FogFlow?

	For any generate results, your application can publish them by sending NGSI updates to the given IoT broker
	The given IoT broker and the entity type to be used for publishing your generate results 
	are configured by FogFlow. 
	
	Once FogFlow launches a task associated with your operator, FogFlow worker will
	send the configuration to the task via the listening port. 
	In the following code block, 
	we show how your operator should handle this provided configuration. 
	The following global variables are configured by FogFlow. 
	
	* brokerURL: URL of the IoT Broker assigned by FogFlow; 
	* ngsi10client: the NGSI10 client created based on the given brokerURL
	* myReferenceURL: the ip address and port number that your application is listening; Using this reference URL your application can subscribe to extra input data from the provided IoT Broker
	* outputs: the array of your output entities, 
	* isConfigured: to indicate whether the configuration has been done

.. code-block:: javascript

    // global variables to be configured 
    var ngsi10client;
    var brokerURL;
    var myReferenceURL;
    var outputs = [];
    var isConfigured = false;
    
    // handle the configuration commands issued by FogFlow
    function handleAdmin(req, commands, res) 
    {	
        handleCmds(commands);
        
        isConfigured = true;
        
        res.status(200).json({});
    }
    
    // handle all configuration commands
    function handleCmds(commands) 
    {
        for(var i = 0; i < commands.length; i++) {
            var cmd = commands[i];
            handleCmd(cmd);
        }	
    }
    
    // handle each configuration command accordingly
    function handleCmd(commandObj) 
    {	
        switch(commandObj.command) {
            case 'CONNECT_BROKER':
                connectBroker(commandObj);
                break;
            case 'SET_OUTPUTS':
                setOutputs(commandObj);
                break;
            case 'SET_YOUR_REFERENCE':
                setReferenceURL(commandObj);
                break;            
        }	
    }
    
    // connect to the IoT Broker
    function connectBroker(cmd) 
    {
        brokerURL = cmd.brokerURL;
        ngsi10client = new NGSIClient.NGSI10Client(brokerURL);
        console.log('connected to broker', cmd.brokerURL);
    }
    
    function setOutputs(cmd) 
    {
        var outputStream = {};
        outputStream.id = cmd.id;
        outputStream.type = cmd.type;
    
        outputs.push(outputStream);
    
        console.log('output has been set: ', cmd);
    }
    
    function setReferenceURL(cmd) 
    {
        myReferenceURL = cmd.referenceURL   
        console.log('your application can subscribe addtional inputs under the reference URL: ', myReferenceURL);
    }
            

* how to handle the received entity data?

.. code-block:: javascript

	// handle all received NGSI notify messages
	function handleNotify(req, ctxObjects, res) 
	{	
		console.log('handle notify');
		for(var i = 0; i < ctxObjects.length; i++) {
			console.log(ctxObjects[i]);
	        fogfunction.handler(ctxObjects[i], publish);
		}
	}

	// process the input data stream accordingly and generate output stream
	function processInputStreamData(data) 
	{
		var type = data.entityId.type;
		console.log('type ', type);
		
		// do the internal data processing
		if (type == 'PowerPanel'){
			// to handle this type of input
		} else if (type == 'Rule') {
			// to handle this type of input
		}	
	}

* how to send an update within your applicaiton?

	For any generate results, your application can publish them by sending NGSI updates to the given IoT broker
	The given IoT broker and the entity type to be used for publishing your generate results 
	are configured by FogFlow. 

.. code-block:: javascript

    // update context for streams
    function updateContext(anomaly) 
    {
        if (isConfigured == false) {
            console.log('the task is not configured yet!!!');
            return;
        }
            
        var ctxObj = {};
        
        ctxObj.entityId = {};
        
        var outputStream = outputs[0];
        
        ctxObj.entityId.id = outputStream.id;
        ctxObj.entityId.type = outputStream.type;
        ctxObj.entityId.isPattern = false;
        
        ctxObj.attributes = {};
        
        ctxObj.attributes.when = {		
            type: 'string',
            value: anomaly['when']
        };
        ctxObj.attributes.whichpanel = {
            type: 'string',
            value: anomaly['whichpanel']
        };  
            
        ctxObj.attributes.shop = {
            type: 'string',
            value: anomaly['whichshop']
        };  
        ctxObj.attributes.where = {
            type: 'object',
            value: anomaly['where']
        };  
        ctxObj.attributes.usage = {
            type: 'integer',
            value: anomaly['usage']
        };
            
        ctxObj.metadata = {};		
        ctxObj.metadata.shop = {
            type: 'string',
            value: anomaly['whichshop']
        };  
                
        ngsi10client.updateContext(ctxObj).then( function(data) {
            console.log('======send update======');
            console.log(ctxObj);
            console.log(data);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to update context');
        });    
    }
        

* what if your application requires some third-party libraries? 

	If your application needs some third-party libraries, please specify them and their version numbers in "package.json". 
	

* how to subscribe addtional inputs if your applicaiton needs?



How to dockerize your operator
-------------------------------

	.. code-block:: bash
		
		# please change package.json to fetch all third-party libraries you need
		npm install 
		
		# build your docker image
		./build 
		
		# push your docker image to the FogFlow docker registry
		docker push task1
		

How to debug and test your operator
---------------------------------------

You can also test your operator manually before using it in your service topology. 
Currently, you can define a subscription like the provided example in "subscription.json" 
and then issue a subscribe request to the known Cloud IoT broker. 

The step is the following: 

#. run your application indepent from FogFlow

	.. code-block:: bash
		
		# please change package.json to fetch all third-party libraries you need
		npm install 
	
		# you need to find out a port which is free and use this free port (e.g., )
		# number to set an environment variable
		export myport=100010
		
		# build your docker image
		node main.js 

#. issue a subscription in order to bring input data to your running application

	The bash script in "test.sh" shows how you can 

	.. code-block:: bash
		
		# please change package.json to fetch all third-party libraries you need
		./test.sh 



*****************************************
Test 
*****************************************

Please follow the following steps to deploy the entire FogFlow system on a single Linux machine before your test. 

`Set up all FogFlow component on a single machine`_

.. _`Set up all FogFlow component on a single machine`: https://fogflow.readthedocs.io/en/latest/setup.html

Once the FogFlow is up and running, an end-to-end function test can be done by performing test.sh in the /test folder. 

.. code-block:: console    
     
	#perform an end-to-end function test
	cd ./test
	./function_test.sh 

For a performance test, please install JMeter and load the provided JMeter test plan ./test/performance_test.jmx





*************************
Monitoring
*************************

Fogflow system health can be monitored by system monitoring tools Metricbeat, Elasticsearch and Grafana in short EMG. 
With these tools edges and Fogflow Docker service health can be monitored. 
Metricbeat deployed on Edge node. Elasticsearch and Grafana on Cloud node.

As illustrated by the following picture, in order to set up FogFlow System Monitoring tools to monitor system resource usage.



.. figure:: figures/Fogflow_System_Monitoring_Architecture.png


Set up Monitoring components on Cloud node
===========================================================


Fetch all required scripts
-------------------------------------------------------------

Download the docker-compose file and the configuration files as below.

.. code-block:: console    

	# the docker-compose file to start all Monitoring components on the cloud node
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/grafana/docker-compose.yml
	
	# the configuration file used by Grafana
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/grafana/grafana.yaml

	# the configuration file used by metric beat
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/grafana/metricbeat.docker.yml

        # JSON files to configure grafana dashboard 
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/grafana/dashboards




.. note:: It is supposed that FogFlow cloud components are in running state before setting up system monitoring.


Start all Monitoring components
----------------------------------

.. code-block:: console  
 
             docker-compose up -d

             #Check all the containers are Up and Running using "docker ps -a"
             docker ps -a



Configure Elasticsearch on Grafana Dashboard
===========================================================  


Grafana dashboard can be accessible on web browser to see the current system status via the URL: http://<output.elasticsearch.hosts>:3003/. 
The default username and password for Grafana login are admin and admin respectively.


- After successful login to grafana, click on "Add data source" on Home Dashboard to setup the source of data.
- Select Elasticsearch from Add Data Sourch page. Now the new page is Data Sources/Elasticsearch same as below figure.


.. figure:: figures/Elastic_config.png



1. Put a name for the Data Source.
2. In HTTP detail ,mention URL of your elasticsearch and Port. URL shall include HTTP. 
3. In Access select Server(default). URL needs to be accessible from the Grafana backend/server.
4. In Elasticsearch details, put @timestamp for Time field name. Here a default for the time field can be specified with the name of your Elasticsearch index. Use a time pattern for the index name or a wildcard.
5. Select Elasticsearch Version i.e. "7.0+".

Then click on "Save & Test" button.

On successful configuration the dashboard will return "Index OK. Time field name OK."


Grafana-based monitoring
===========================================================  
        
Grafana based system metrics can be seen from Home Dashboard page of Grafana. In the sidebar, take the cursor over Dashboards (squares) icon, and then click Manage. The dashboard appears in a Services folder.

Here are some basic Grafana visualization dashboard to monitor metrics of FogFlow cloud as well as edge nodes.

- **Below diagram illustrate steps to setup dashboard for containers list with maximum memory usage**.




.. figure:: figures/Container_max_memory_usage.png




- **Below diagram illustrate steps to setup dashboard to show system memory used in bytes**.




.. figure:: figures/System_Memory_Gauge.png



- **Below diagram illustrate steps to setup dashboard to show system metric data rate in packet per second**.



.. figure:: figures/System_Metric_filter.png



- **Below diagram illustrate steps to setup dashboard to show FogFlow Cloud and Edge nodes that are live**.


.. figure:: figures/Fogflow_Cloud_Edge_Nodes.png





Set up the Metricbeat on Edge node
-------------------------------------

Download the metric beat yml file for edge node.

.. code-block:: console  

            # the configuration file used by metric beat
            wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/grafana/metricbeat.docker.yml

**Optional** Edit "name" in metricbeat.docker.yml file to add particular name for better identification of edge node.

Copy below Docker run command, edit and replace <Cloud_Public_IP> with IP/URL of elasticsearch in output.elasticsearch.hosts=["<Cloud_Public_IP>:9200"]>. This command will deploy metric beat on edge node.

.. code-block:: console  

            docker run -d   --name=metricbeat   --user=root   --volume="$(pwd)/metricbeat.docker.yml:/usr/share/metricbeat/metricbeat.yml:ro"   --volume="/var/run/docker.sock:/var/run/docker.sock:ro"   --volume="/sys/fs/cgroup:/hostfs/sys/fs/cgroup:ro"   --volume="/proc:/hostfs/proc:ro"   --volume="/:/hostfs:ro"   docker.elastic.co/beats/metricbeat:7.6.0 metricbeat -e   -E output.elasticsearch.hosts=["<Cloud_Public_IP>:9200"]


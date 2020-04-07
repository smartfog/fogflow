Setup Grafana Dashboard for various system metrics
===========================================================  
        
To monitor metrics of FogFlow cloud as well as edge nodes in graphical format we need to setup dashboard.
Here are some basic Grafana visualization dashboard setting examples to monitor system resources.

- **Below diagram illustrate steps to setup dashboard for containers list with maximum memory usage**.


.. figure:: figures/Container_max_memory_usage.png


1. To create query for Elasticsearch select Query: Metrics: Average(docker.memory.usage.max), Group by: Terms(host.name), Terms(container.image.name), Date Histogram(@timestamp) from drop down list.
2. Click on Visualization select Graph from drop down , Draw Modes (Lines), Mode Options(Fill:1,Fill Gradient:0,Line Width:2), Stacking & Null value(Null value:connected)
   Axes- Left Y(Unit:bytes,Scale:linear), Right Y(Unit:short,Scale:linear), X-Axis(Mode:Time)
   Legend- Options(Show,As Table,To the right), Values(Max)
3. Click on General Title: Container memory usage max, write Description if there is any description.


- **Below diagram illustrate steps to setup dashboard to show system memory used in bytes**.


.. figure:: figures/System_Memory_Gauge.png


1. To create query for Elasticsearch select Query: memory, Metrics: Average(system.memory.actual.used.bytes), Group by: Terms(host.name), Date Histogram(@timestamp)from drop down list.
2. Click on Visualization select Gauge from drop down , Display (Show:Calculation, Calc:Last(not null), Labels, Markers), Field (Unit:bytes, Min:0, Max:100), Thresholds (50 (yellow), base (green)).
3. Click on General Title: System memory used in bytes, write Description if there is any description.

- **Below diagram illustrate steps to setup dashboard to show system metric data rate in packet per second**.

.. figure:: figures/System_Metric_filter.png

1. To create query for Elasticsearch select Query: Metrics: Average(system.memory.actual.used.bytes), Group by: Terms(agent.name), Date Histogram(@timestamp)from drop down list.
2. Click on Visualization select Graph from drop down , Draw Modes (Lines), Mode Options(Fill:1,Fill Gradient:0,Line Width:2), Hover tooltip(Mode: All series, Sort order:Increasing), Stacking & Null value(Null value:connected).
   Axes- Left Y(Unit:packets/sec, Scale:linear), Right Y(Unit:packets/sec, Scale:linear), X-Axis(Mode:Time)
   Legend- Options(Show,As Table,To the right), Values(Avg)
3. Click on General Title: System Metric filter, write Description if there is any description.


- **Below diagram illustrate steps to setup dashboard to show FogFlow Cloud and Edge nodes that are live**.


.. figure:: figures/Fogflow_Cloud_Edge_Nodes.png


1. To create query for Elasticsearch select Query: Metrics: Count(), Group by: Terms(agent.name), Date Histogram(@timestamp) from drop down list.
2. Click on Visualization select Graph from drop down , Draw Modes (Lines), Mode Options(Fill:1,Fill Gradient:0,Line Width:2).
   Axes- Left Y(Unit:bytes, Scale:linear), Right Y(Unit:short, Scale:linear), X-Axis(Mode:Time).
   Legend- Options(Show, As Table, To the right), Values(Avg, Max).

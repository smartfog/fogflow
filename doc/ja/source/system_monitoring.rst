*************************
モニタリング
*************************


Grafana ベースのモニタリング
===========================================================  
        
FogFlow クラウドのメトリックとエッジ ノードをグラフィック形式で監視するには、ダッシュボードを設定する必要があります。システム リソースを監視するための基本的な Grafana 視覚化ダッシュボード設定の例を次に示します。

- **次の図は、メモリ使用量が最大のコンテナー リストのダッシュボードを設定する手順を示しています。**


.. figure:: ../../en/source/figures/Container_max_memory_usage.png


1. Elasticsearch のクエリを作成するには、ドロップダウンリストから Query: Metrics: Average(docker.memory.usage.max), Group by: Terms(host.name), Terms(container.image.name), Date Histogram(@timestamp) を選択します。
2. ドロップダウンから Visualization select Graph をクリックします。Draw Modes (Lines), Mode Options(Fill:1,Fill Gradient:0,Line Width:2), Stacking & Null value(Null value:connected)
   Axes- Left Y(Unit:bytes,Scale:linear), Right Y(Unit:short,Scale:linear), X-Axis(Mode:Time)
   Legend- Options(Show,As Table,To the right), Values(Max)
3. General Title: Container memory usage max, をクリックし、説明がある場合は説明を入力します。


- **次の図は、使用されているシステム メモリをバイト単位で表示するようにダッシュボードを設定する手順を示しています。**


.. figure:: ../../en/source/figures/System_Memory_Gauge.png


1. Elasticsearch のクエリを作成するには、ドロップダウンリストから、Query: memory, Metrics: Average(system.memory.actual.used.bytes), Group by: Terms(host.name), Date Histogram(@timestamp) を選択します。
2. ドロップダウンから Visualization select Gauge  をクリックします。 Display (Show:Calculation, Calc:Last(not null), Labels, Markers), Field (Unit:bytes, Min:0, Max:100), Thresholds (50 (yellow), base (green)
3. General Title: System memory used in bytes をクリックし、説明がある場合は説明を入力します。

- **次の図は、ダッシュボードをセットアップして、システム メトリック データレートをパケット/秒で表示する手順を示しています。**


.. figure:: ../../en/source/figures/System_Metric_filter.png

1. Elasticsearch のクエリを作成するには、ドロップダウンリストから、Query: Metrics: Average(system.memory.actual.used.bytes), Group by: Terms(agent.name), Date Histogram(@timestamp) を選択します。
2. ドロップダウンから Visualization select Graph  をクリックします。Draw Modes (Lines), Mode Options(Fill:1,Fill Gradient:0,Line Width:2), Hover tooltip(Mode: All series, Sort order:Increasing), Stacking & Null value(Null value:connected).
   Axes- Left Y(Unit:packets/sec, Scale:linear), Right Y(Unit:packets/sec, Scale:linear), X-Axis(Mode:Time)
   Legend- Options(Show,As Table,To the right), Values(Avg)
3. General Title: System Metric filter をクリックし、説明がある場合は説明を入力します。


- **次の図は、ライブの FogFlow クラウド ノードとエッジ ノードを表示するダッシュボードを設定する手順を示しています**.


.. figure:: ../../en/source/figures/Fogflow_Cloud_Edge_Nodes.png


1. Elasticsearch のクエリを作成するには、ドロップダウンリストから、Query: Metrics: Count(), Group by: Terms(agent.name), Date Histogram(@timestamp) を選択します。
2. ドロップダウンから Visualization select Graph をクリックします。Draw Modes (Lines), Mode Options(Fill:1,Fill Gradient:0,Line Width:2).
   Axes- Left Y(Unit:bytes, Scale:linear), Right Y(Unit:short, Scale:linear), X-Axis(Mode:Time).
   Legend- Options(Show, As Table, To the right), Values(Avg, Max).

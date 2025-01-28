*************************
モニタリング
*************************


Fogflow システムの状態は、システム監視ツールの Metricbeat、Elasticsearch、Grafana つまり EMG によって監視できます。これらのツールを使用すると、エッジと Fogflow Docker サービスの状態を監視できます。エッジノードにデプロイされた Metricbeat。クラウド ノードにデプロイされた Elasticsearch と Grafana。

次の図に示すように、FogFlow システム監視ツールを設定して、システムリソースの使用状況を監視します。


.. figure:: figures/Fogflow_System_Monitoring_Architecture.png



Grafana ダッシュボードで Elasticsearch を構成
===========================================================  


Grafana ダッシュボード はWeb ブラウザーからアクセスでき、URL: http:///<output.elasticsearch.hosts>:3003/ を介して現在のシステム ステータスを確認できます 。Grafana ログインのデフォルトのユーザー名とパスワードは、それぞれ admin と admin です。

- grafana に正常にログインしたら、ホーム ダッシュボードの "Create your first data source" をクリックして、データソースを設定します。
- Add Data Source ページからElasticsearch を選択します。これで、下の図と同じページの Data Sources/Elasticsearch が表示されます。


.. figure:: figures/Elastic_config.png


1. データソースに名前を付けます。
2. HTTP の詳細で、elasticsearch とポートの URL に言及します。URL には HTTP を含める必要があります。
3. Access で Server(default) を選択します。URL は、Grafana バックエンド/サーバーからアクセスできる必要があります。
4. Elasticsearch の詳細で、Time フィールド名に @timestamp を入力します。ここで、時間フィールドのデフォルトを Elasticsearch インデックスの名前で指定できます。インデックス名またはワイルドカードには時間パターンを使用します。
5. Elasticsearch バージョンを選択します。

次に、"Save & Test" ボタンをクリックします。


Metricbeat を設定
===========================================================  


- 以下のように、metricbeat.docker.yml ファイルの Elasticsearch の詳細を変更します:


.. code-block:: json
        name: "<155.54.239.141_cloud>"
        metricbeat.modules:
        - module: docker
          #Docker module parameters that has to be monitored based on user requirement, example as below
          metricsets: ["cpu","memory","network"]
          hosts: ["unix:///var/run/docker.sock"]
          period: 10s
          enabled: true
        - module: system
          #System module parameters that has to be monitored based on user requirement, example as below
          metricsets: ["cpu","load","memory","network"]
          period: 10s
        output.elasticsearch:
          hosts: '155.54.239.141:9200'
	  
Fogflow システムの状態は、システム監視ツールの Metricbeat、Elasticsearch、Grafana、つまり EMG によって監視できます。これらのツールを使用すると、エッジとFogflow Docker サービスの状態を監視できます。エッジノードにデプロイされた Metricbeat。クラウド ノードにデプロイされた Elasticsearch と Grafana。

次の図に示すように、FogFlow システム監視ツールを設定して、システムリソースの使用状況を監視します。


.. figure:: figures/Fogflow_System_Monitoring_Architecture.png


1台のマシンですべての FogFlow コンポーネントをセットアップ
===========================================================


必要なすべてのスクリプトを取得
-------------------------------------------------------------

以下のように docker-compose ファイルと構成ファイルをダウンロードします。

.. code-block:: console    

	# the docker-compose file to start all FogFlow components on the cloud node
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/docker-compose.yml

	# the configuration file used by all FogFlow components
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/config.json

	# the configuration file used by the nginx proxy
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/nginx.conf

        # the configuration file used by metricbeat
        wget https://raw.githubusercontent.com/smartfog/fogflow/master/docker/core/http/metricbeat.docker.yml


環境に応じてelasticsearchとmetricbeatのIP構成を変更
---------------------------------------------------------------------------

ご使用の環境に応じて、docker-compose.yml の次の IP アドレスを変更する必要があります。

- **output.elasticsearch.hosts**: これは、metricbeat が csv 形式でデータを共有する elasticsearch のホストの場所です。

また、ご使用の環境に応じて、metricbeat.docker.yml の次の IP アドレスを変更する必要があります。

- **name**: これは、Grafana メトリック ダッシュボードのエッジ　ノードからのクラウド ノードの一意性に付けられた名前です。IP アドレスの代わりに任意の名前を指定できます。

- **hosts**: metricsearh データベースのホストの場所であり、metricbeat がメトリックデータを共有します。


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


Grafana ダッシュボードで Elasticsearch を構成
-------------------------------------------------------------

Grafana ダッシュボードは Web ブラウザーからアクセスでき、URL: http:///<output.elasticsearch.hosts>:3003/ を介して現在のシステム ステータスを確認できます 。Grafana ログインのデフォルトのユーザー名とパスワードは、それぞれ admin と admin です。

- Grafana に正常にログインしたら、ホーム ダッシュボードの "Create your first data source" をクリックして、データソースを設定します。
- Add Data Source ページから Elasticsearch を選択します。これで、下の図と同じページの Data Sources/Elasticsearch が表示されます。


.. figure:: figures/Elastic_config.png


1. データソースに名前を付けます。
2. HTTP の詳細で、elasticsearch とポートの URL に言及します。URL には HTTP を含める必要があります。
3. Access で Server(default) を選択します。URL は、Grafana バックエンド/サーバーからアクセスできる必要があります。
4. Elasticsearch の詳細で、Time フィールド名に @timestamp を入力します。ここで、Time フィールドのデフォルト をElasticsearch インデックスの名前で指定できます。インデックス名またはワイルドカードには時間パターンを使用します。
5. Elasticsearch バージョンを選択します。

次に、"Save & Test" ボタンをクリックします。


Metricbeatを設定
---------------------------------------------


- 以下のように、metricbeat.docker.yml ファイルの Elasticsearch の詳細を変更します:


.. code-block:: json
        name: "<155.54.239.141_cloud>"
        metricbeat.modules:
        - module: docker
          #Docker module parameters that has to be monitored based on user requirement, example as below
          metricsets: ["cpu","memory","network"]
          hosts: ["unix:///var/run/docker.sock"]
          period: 10s
          enabled: true
        - module: system
          #System module parameters that has to be monitored based on user requirement, example as below
          metricsets: ["cpu","load","memory","network"]
          period: 10s
        output.elasticsearch:
          hosts: '155.54.239.141:9200'


これは、FogFlow の1ページの入門チュートリアルです。FIWARE ベースのアーキテクチャでは、FogFlow を使用して、エッジ ノード (IoT ゲートウェイや Raspberry Pi など）で生データを変換および前処理する目的で、IoT デバイスとOrion Context Broker の間でデータ処理機能を動的にトリガーできます。

チュートリアルでは、温度センサー データのエッジで異常検出を行う簡単な例を使用して、一般的な FogFlow システムのセット アップを紹介します。FogFlow と Orion Context Broker を相互に統合して使用するユースケース実装の例について説明します。

ユースケースを実装するたびに、FogFlow は、オペレーター、Docker イメージ、フォグ ファンクション 、サービス トポロジーなどの内部 NGSI エンティティを作成します。したがって、これらのエンティティ データは FogFlow システムにとって非常に重要であり、どこかに保存する必要があります。メモリは揮発性であり、電源が失われるとコンテンツが失われるため、エンティティ データを FogFlow メモリに保存することはできません。この問題を解決するために、FogFlow は Dgraph  に永続ストレージを導入します。永続ストレージは、FogFlow エンティティ データをグラフの形式で保存します。


.. _`Dgraph`: https://dgraph.io/docs/get-started/


次の図に示すように、このユースケースでは、接続された温度センサーが更新メッセージを FogFlow システムに送信します。これにより、事前定義されたフォグ ファンクションの実行中のタスク インスタンスがトリガーされ、分析結果が生成されます。フォグ ファンクションは、FogFlow ダッシュボードで事前に指定されていますが、温度センサーがシステムに接続された場合にのみトリガーされます。実際の分散セットアップでは、実行中のタスク インスタンスは、温度センサーに近いエッジ ノードにデプロイされます。生成された分析結果が生成されると、FogFlow システムから Orion Context Broker に転送されます。これは、参照 URL として Orion Context Broker を使用したサブスクリプションが発行されたためです。


.. figure:: ../../en/source/figures/systemview.png


FogFlow システムの状態は、システム監視ツールの Metricbeat、Elasticsearch、Grafana、つまり EMG によって監視できます。これらのツールを使用すると、エッジとFogFlow Docker サービスの状態を監視できます。エッジ ノードにデプロイされた Metricbeat。クラウド ノードにデプロイされた Elasticsearch と Grafana。

次の図に示すように、FogFlow システム監視ツールをセットアップしてシステム リソースの使用状況を監視するため。


.. figure:: ../../en/source/figures/Fogflow_System_Monitoring_Architecture.png


FogFlow を実行するための前提条件のコマンドは次のとおりです:

1. docker

2. docker-compose


Ubuntu 16.04 の場合、docker-ce と docker-compose をインストールする必要があります。


Docker CE をインストールするには、`Install Docker CE`_ を参照してください。必要なバージョン > 18.03.1-ce です。


.. important:: 
	**また、ユーザーが sudo なしで Docker コマンドを実行できるようにしてください**


Docker Compose をインストールするには、`Install Docker Compose`_ を参照してください。必要なバージョン> 2.4.2 です。


.. _`Install Docker CE`: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04
.. _`Install Docker Compose`: https://www.digitalocean.com/community/tutorials/how-to-install-docker-compose-on-ubuntu-16-04



1台のマシンですべてのFogFlowコンポーネントをセットアップ
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
	
	
IP構成を変更
-------------------------------------------------------------

ご使用の環境に応じて、config.json で以下の IP アドレスを変更する必要があります。

- **coreservice_ip**: すべての FogFlow エッジ ノードが FogFlow クラウド ノードのコア サービス (ポート80 の nginx やポート 5672 の rabbitmq など) にアクセスするために使用します。通常、これは FogFlow クラウド ノードのパブリック IP になります。

- **external_hostip**: FogFlow クラウド ノードの構成の場合、これは、実行中の FogFlow コア サービスにアクセスするためにコンポーネント (Cloud Worker および Cloud Broker) によって使用される coreservice_ip と同じです。

- **internal_hostip**: これはデフォルトの Docker ブリッジの IP であり、Linux ホストの "docker0" ネットワーク インターフェースです。Windows または MacOS のDocker エンジンの場合、"docker0" ネットワークインターフェイスはありません。代わりに、特別なドメイン名 "host.docker.internal" を使用する必要があります。

- **site_id**: 各 FogFlow ノード (クラウド ノードまたはエッジ ノード) は、システム内で自身を識別するために一意の文字列ベースの ID を持っている必要があります。

- **physical_location**: FogFlow ノードの地理的位置。

- **worker.capacity**: FogFlow ノードが呼び出すことができる Docker コンテナーの最大数を意味します。


Elasticsearch と Metricbeat の IP 構成を変更
---------------------------------------------------------------------------

ご使用の環境に応じて、docker-compose.yml の次の IP アドレスを変更する必要があります。

- **output.elasticsearch.hosts**: metricbeat が csv 形式でデータを共有する elasticsearch のホストの場所です。

また、ご使用の環境に応じて、metricbeat.docker.yml の次の IP アドレスを変更する必要があります。

- **name**: Grafana metric ダッシュボードのエッジ ノードからのクラウド ノードの一意性に付けられた名前です。IP アドレスの代わりに任意の名前を指定できます。

- **hosts**: elasticsearh データベースのホストの場所であり、metricbeat がメトリックデータを共有します。


.. important:: 

        **coreservice_ip** および **external_hostip** の IP アドレスとして "127.0.0.1" を使用しないでください。これらは Docker コンテナー内で実行中のタスクによって使用されるためです。
	
	**Firewall rules:** external_ip を介して FogFlow Web ポータルにアクセスできるようにします。次のポートも開いている必要があります: TCP の場合は80 および 5672

	**Mac Users:** Macbook で FogFlow をテストする場合は、Docker デスクトップをインストールし、"host.docker.internal" を使用して構成ファイルの coreservice_ip、external_hostip、internal_hostip を置き換えてください。



すべての FogFlow コンポーネントを起動
-------------------------------------------------------------


すべての FogFlow コンポーネントの Docker イメージをプルし、FogFlow システムを起動します。


.. code-block:: console    

	#if you already download the docker images of FogFlow components, this command can fetch the updated images
	docker-compose pull  

	docker-compose up -d


セットアップを検証
-------------------------------------------------------------


FogFlow クラウド ノードが正しく開始されているかどうかを確認するには、次の2つの方法があります:


- "docker ps -a" を使用して、すべてのコンテナーが稼働していることを確認します。


.. code-block:: console    

	docker ps -a
	
	CONTAINER ID      IMAGE                       COMMAND                  CREATED             STATUS              PORTS                                                 NAMES
	90868b310608      nginx:latest            "nginx -g 'daemon of…"   5 seconds ago       Up 3 seconds        0.0.0.0:80->80/tcp                                       fogflow_nginx_1
	d4fd1aee2655      fogflow/worker          "/worker"                6 seconds ago       Up 2 seconds                                                                 fogflow_cloud_worker_1
	428e69bf5998      fogflow/master          "/master"                6 seconds ago       Up 4 seconds        0.0.0.0:1060->1060/tcp                               fogflow_master_1
	9da1124a43b4      fogflow/designer        "node main.js"           7 seconds ago       Up 5 seconds        0.0.0.0:1030->1030/tcp, 0.0.0.0:8080->8080/tcp       fogflow_designer_1
	bb8e25e5a75d      fogflow/broker          "/broker"                9 seconds ago       Up 7 seconds        0.0.0.0:8070->8070/tcp                               fogflow_cloud_broker_1
	7f3ce330c204      rabbitmq:3              "docker-entrypoint.s…"   10 seconds ago      Up 6 seconds        4369/tcp, 5671/tcp, 25672/tcp, 0.0.0.0:5672->5672/tcp     fogflow_rabbitmq_1
	9e95c55a1eb7      fogflow/discovery       "/discovery"             10 seconds ago      Up 8 seconds        0.0.0.0:8090->8090/tcp                               fogflow_discovery_1
        399958d8d88a      grafana/grafana:6.5.0   "/run.sh"                29 seconds ago      Up 27 seconds       0.0.0.0:3003->3000/tcp                               fogflow_grafana_1
        9f99315a1a1d      fogflow/elasticsearch:7.5.1 "/usr/local/bin/dock…" 32 seconds ago    Up 29 seconds       0.0.0.0:9200->9200/tcp, 0.0.0.0:9300->9300/tcp       fogflow_elasticsearch_1
        57eac616a67e      fogflow/metricbeat:7.6.0 "/usr/local/bin/dock…"   32 seconds ago     Up 29 seconds                                                                  fogflow_metricbeat_1


.. important:: 

        不足しているコンテナーがある場合は、"docker ps -a" を実行して、FogFlow コンポーネントが何らかの問題で終了していないかどうかを確認できます。ある場合は、"docker logs [container ID]" を実行して、出力ログをさらに確認できます。


- FogFlow DashBoard からシステム ステータスを確認します。

Web ブラウザで FogFlow ダッシュボードを開くと、次の URL を介して現在のシステム ステータスを確認できます: http://<coreservice_ip>/index.html


.. important:: 

i       FogFlow クラウド ノードがゲートウェイの背後にある場合は、ゲートウェイ IP から coreservice_ip へのマッピングを作成してから、ゲートウェイ IP を介して FogFlow ダッシュボードにアクセスする必要があります。
        FogFlow クラウド ノードが AzureCloud、Google Cloud、Amazon Cloud などのパブリッククラウド内の VM である場合は、VM のパブリック IP を介して FogFlow ダッシュボードにアクセスする必要があります。
	

FogFlow ダッシュボードにアクセスできるようになると、次の Web ページが表示されます:


.. figure:: ../../en/source/figures/dashboard.png



Grafana ダッシュボードで Elasticsearch を構成
-------------------------------------------------------------

Grafana ダッシュボードは Web ブラウザーからアクセスでき、URL: http://<output.elasticsearch.hosts>:3003/ を介して現在のシステム ステータスを確認できます。Grafana ログインのデフォルトのユーザー名とパスワードは、それぞれ admin と admin です。


- Grafana に正常にログインしたら、ホームダッシュボードの "Create your first data source" をクリックして、データソースを設定します。

- Add Data Sourch ページから Elasticsearch を選択します。これで、下の図と同じページの Data Sources/Elasticsearch が表示されます。


.. figure:: ../../en/source/figures/Elastic_config.png


1. データソースに名前を付けます。
2. HTTP の詳細で、elasticsearch とポートの URL に言及します。URL には HTTP を含める必要があります。
3. Access で Server(default) を選択します。URL は、Grafana バックエンド/サーバーからアクセスできる必要があります。
4. Elasticsearch の詳細で、Time フィールド名に @timestamp を入力します。ここで、時間フィールドのデフォルトを Elasticsearch インデックスの名前で指定できます。インデックス名またはワイルドカードには時間パターンを使用します。
5. Elasticsearch バージョンを選択します。

次に、"Save & Test" ボタンをクリックします。


Metricbeat を設定
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
	  

既存の IoT サービスを試す
===========================================================

FogFlow クラウド ノードがセットアップされると、FogFlow エッジ ノードを実行せずに既存の IoT サービスを試すことができます。たとえば、次のような簡単なフォグ ファンクションを試すことができます。


3回のクリックですべての定義済みサービスを初期化
-------------------------------------------------------------

- 上部のナビゲーター バーにある "Operator Registry" をクリックして、事前定義されたオペレーターの初期化をトリガーします。

最初に "Operator Registry" をクリックすると、事前定義されたオペレーターのリストが FogFlow システムに登録されます。2回クリックすると、次の図に示すように、更新されたリストが表示されます。


.. figure:: ../../en/source/figures/operator-list.png


- 上部のナビゲーター バーで "Service Topology" をクリックして、事前定義されたサービス トポロジーの初期化をトリガーします。

最初に "Service Topology" をクリックすると、事前定義されたトポロジーのリストが FogFlow システムに登録されます。2回クリックすると、次の図に示すように、更新されたリストが表示されます。


.. figure:: ../../en/source/figures/topology-list.png


- 上部のナビゲーターバーの "Fog Function" をクリックして、事前定義されたフォグ ファンクションの初期化をトリガーします。

最初に "Fog Function" をクリックすると、事前定義されたファンクションのリストが FogFlow システムに登録されます。2回クリックすると、次の図に示すように、更新されたリストが表示されます。


.. figure:: ../../en/source/figures/function-list.png


IoT デバイスをシミュレートしてフォグ ファンクションをトリガー
-------------------------------------------------------------

フォグファンクションをトリガーする方法は2つあります:

**1. FogFlow ダッシュボードを介して “Temperature” センサーエンティティを作成**

デバイス登録ページからデバイス エンティティを登録できます: "System Status" -> "Device" -> "Add"。次に、次の要素を入力して "Temperature" センサー エンティティを作成できます:

- **Device ID:** 一意のエンティティIDを指定します。
- **Device Type:** エンティティ タイプとして “Temperature” を使用します。
- **Location:** マップ上の場所 (location) を選択します。
 

.. figure:: ../../en/source/figures/device-registration.png

**2. NGSI エンティティの更新を送信して、“Temperature” センサーエンティティを作成**
 
エンティティの更新のために FogFlow Broker に curl リクエストを送信します:

.. code-block:: console    

	
	curl -iX POST \
		  'http://coreservice_ip/ngsi10/updateContext' \
		  -H 'Content-Type: application/json' \
		  -d '
		{
		    "contextElements": [
		        {
		            "entityId": {
		                "id": "Device.Temp001",
		                "type": "Temperature",
		                "isPattern": false
		                },
		            "attributes": [
		                    {
		                    "name": "temperature",
		                    "type": "float",
		                    "value": 73
		                    },
		                    {
		                    "name": "pressure",
		                    "type": "float",
		                    "value": 44
		                    }
		                ],
		            "domainMetadata": [
		                    {
		                    "name": "location",
		                    "type": "point",
		                    "value": {
		                    "latitude": -33.1,
		                    "longitude": -1.1
		                    }}
		                ]
		        }
		    ],
		    "updateAction": "UPDATE"
		}'


フォグ ファンクションがトリガーされているかどうかを確認
-------------------------------------------------------------

システム管理 (System Management) の "Task" の下にタスクが作成されているかどうかを確認します。

.. figure:: ../../en/source/figures/fog-function-task-running.png

システム管理の "Stream" の下にストリームが作成されているかどうかを確認します。

.. figure:: ../../en/source/figures/fog-function-streams.png



FogFlow を OrionBroker と統合
======================================


Orion を起動
-------------------------------------------------------------

Orion のドキュメントに従って、Orion Context Broker インスタンスをここからセットアップできます: `Installing Orion`.

.. _`Installing Orion`: https://fiware-orion.readthedocs.io/en/master/admin/install/index.html


以下のコマンドを使用して、Docker で Orion をセットアップすることもできます (docker はこの方法で必要です) 注意: Orion コンテナーは MongoDB データベースに依存しています。

**前提条件:** Docker がインストールされている必要があります。

最初に、以下のコマンドを使用して MongoDB コンテナーを起動します:

.. code-block:: console    

	sudo docker run --name mongodb -d mongo:3.4


そして、このコマンドで Orion を実行します。

.. code-block:: console    

	sudo docker run -d --name orion1 --link mongodb:mongodb -p 1026:1026 fiware/orion -dbhost mongodb


すべてが動作することを確認します。

.. code-block:: console    

	curl http://<Orion IP>:1026/version

注意: パブリック アクセスのためにファイアウォールのポート 1026 を許可します。



サブスクリプションを発行して、生成された結果を Orion Context Broker に転送
----------------------------------------------------------------------------------

次の curl リクエストを使用して、FogFlow Broker を FIWARE Orion にサブスクライブします:

.. code-block:: console    

	curl -iX POST \
	  'http://coreservice_ip/ngsi10/subscribeContext' \
	  -H 'Content-Type: application/json'  \
	  -H 'Destination: orion-broker'  \
	  -d '
	{
	  "entities": [
	    {
	      "id": ".*",
	      "type": "Result",
	      "isPattern": true
	    }
	  ],
	  "reference": "http://<Orion IP>:1026/v2/op/notify"
	}'


このサブスクリプション リクエストは制限や属性を使用せず、エンティティ タイプに基づく一般的なサブスクリプション リクエストであることに注意してください。


Orion Context Broker から結果をクエリ
-------------------------------------------------------------

ブラウザで次の URL にアクセスし、目的のコンテキスト エンティティを検索します:

.. code-block:: console    

	curl http://<Orion IP>:1026/v2/entities/

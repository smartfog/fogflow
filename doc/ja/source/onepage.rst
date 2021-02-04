これは、FogFlow の1ページの入門チュートリアルです。FIWARE ベースのアーキテクチャでは、FogFlow を使用して、エッジ ノード (IoT ゲートウェイや Raspberry Pi など）で生データを変換および前処理する目的で、IoT デバイスとOrion Context Broker の間でデータ処理機能を動的にトリガーできます。

チュートリアルでは、温度センサー データのエッジで異常検出を行う簡単な例を使用して、一般的な FogFlow システムのセット アップを紹介します。FogFlow と Orion Context Broker を相互に統合して使用するユースケース実装の例について説明します。

ユースケースを実装するたびに、FogFlow は、オペレーター、Docker イメージ、フォグ ファンクション 、サービス トポロジーなどの内部 NGSI エンティティを作成します。したがって、これらのエンティティ データは FogFlow システムにとって非常に重要であり、どこかに保存する必要があります。メモリは揮発性であり、電源が失われるとコンテンツが失われるため、エンティティ データを FogFlow メモリに保存することはできません。この問題を解決するために、FogFlow は Dgraph  に永続ストレージを導入します。永続ストレージは、FogFlow エンティティ データをグラフの形式で保存します。


.. _`Dgraph`: https://dgraph.io/docs/get-started/


次の図に示すように、このユースケースでは、接続された温度センサーが更新メッセージを FogFlow システムに送信します。これにより、事前定義されたフォグ ファンクションの実行中のタスク インスタンスがトリガーされ、分析結果が生成されます。フォグ ファンクションは、FogFlow ダッシュボードで事前に指定されていますが、温度センサーがシステムに接続された場合にのみトリガーされます。実際の分散セットアップでは、実行中のタスク インスタンスは、温度センサーに近いエッジ ノードにデプロイされます。生成された分析結果が生成されると、FogFlow システムから Orion Context Broker に転送されます。これは、参照 URL として Orion Context Broker を使用したサブスクリプションが発行されたためです。


.. figure:: ../../en/source/figures/systemview.png


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
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/release/3.2/cloud/docker-compose.yml

	# the configuration file used by all FogFlow components
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/release/3.2/cloud/config.json

	# the configuration file used by the nginx proxy
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/release/3.2/cloud/nginx.conf
	
	
IP構成を変更
-------------------------------------------------------------

ご使用の環境に応じて、config.json で以下の IP アドレスを変更する必要があります。

- **my_hostip **：これはホストマシンの IP であり、ホストマシンの Web ブラウザと Docker コンテナの両方からアクセスできる必要があります。これには "127.0.0.1" を使用しないでください。
- **site_id**: 各 FogFlow ノード (クラウド ノードまたはエッジ ノード) は、システム内で自身を識別するために一意の文字列ベースの ID を持っている必要があります。
- **physical_location**: FogFlow ノードの地理的位置。
- **worker.capacity**: FogFlow ノードが呼び出すことができる Docker コンテナーの最大数を意味します。


.. important:: 

       "127.0.0.1" を **my_hostip** の IP アドレスとして使用しないでください。これは、Docker コンテナ内で実行中のタスクにのみアクセスできるためです。
	
	**Firewall rules:** FogFlow Web ポータルにアクセスできるようにするには、TCP を介した次のポート 80 および 5672 が開いている必要があります。

	**Mac Users:** Macbook で FogFlow をテストする場合は、Docker デスクトップをインストールし、構成ファイルの my_hostip として "host.docker.internal" も使用してください。

        ポート番号を変更する必要がある場合は、変更がこれら3つの構成ファイルすべてで一貫していることを確認してください。


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
	795e6afe2857   nginx:latest            "/docker-entrypoint.…"   About a minute ago   Up About a minute   0.0.0.0:80->80/tcp                                                                               fogflow_nginx_1
	33aa34869968   fogflow/worker:3.2      "/worker"                About a minute ago   Up About a minute                                                                                                    fogflow_cloud_worker_1
	e4055b5cdfe5   fogflow/master:3.2      "/master"                About a minute ago   Up About a minute   0.0.0.0:1060->1060/tcp                                                                           fogflow_master_1
	cdf8d4068959   fogflow/designer:3.2    "node main.js"           About a minute ago   Up About a minute   0.0.0.0:1030->1030/tcp, 0.0.0.0:8080->8080/tcp                                                   fogflow_designer_1
	56daf7f078a1   fogflow/broker:3.2      "/broker"                About a minute ago   Up About a minute   0.0.0.0:8070->8070/tcp                                                                           fogflow_cloud_broker_1
	51901ce6ee5f   fogflow/discovery:3.2   "/discovery"             About a minute ago   Up About a minute   0.0.0.0:8090->8090/tcp                                                                           fogflow_discovery_1
	51eff4975621   dgraph/standalone       "/run.sh"                About a minute ago   Up About a minute   0.0.0.0:6080->6080/tcp, 0.0.0.0:8000->8000/tcp, 0.0.0.0:8082->8080/tcp, 0.0.0.0:9082->9080/tcp   fogflow_dgraph_1
	eb31cd255fde   rabbitmq:3              "docker-entrypoint.s…"   About a minute ago   Up About a minute   4369/tcp, 5671/tcp, 15691-15692/tcp, 25672/tcp, 0.0.0.0:5672->5672/tcp                           fogflow_rabbitmq_1


.. important:: 

        不足しているコンテナーがある場合は、"docker ps -a" を実行して、FogFlow コンポーネントが何らかの問題で終了していないかどうかを確認できます。ある場合は、"docker logs [container ID]" を実行して、出力ログをさらに確認できます。


- FogFlow DashBoard からシステム ステータスを確認します。

Web ブラウザで FogFlow ダッシュボードを開くと、次の URL を介して現在のシステム ステータスを確認できます: http://<coreservice_ip>/index.html


.. important:: 

i       FogFlow クラウド ノードがゲートウェイの背後にある場合は、ゲートウェイ IP から coreservice_ip へのマッピングを作成してから、ゲートウェイ IP を介して FogFlow ダッシュボードにアクセスする必要があります。
        FogFlow クラウド ノードが AzureCloud、Google Cloud、Amazon Cloud などのパブリッククラウド内の VM である場合は、VM のパブリック IP を介して FogFlow ダッシュボードにアクセスする必要があります。
	

FogFlow ダッシュボードにアクセスできるようになると、次の Web ページが表示されます:


.. figure:: ../../en/source/figures/dashboard.png



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
		  'http://my_hostip/ngsi10/updateContext' \
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

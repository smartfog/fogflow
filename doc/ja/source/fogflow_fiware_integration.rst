*****************************************
FogFlow を FIWARE と統合
*****************************************

FogFlow は、FIWARE エコシステムの公式 Generic Enabler (GE) として、Cloud-Edge Orchestrator として独自の地位を占めており、データの取り込み、変換、および高度な分析のために、クラウドとエッジ上で動的なデータ処理フローをシームレスに起動および管理します。

次の図に示すように、FogFlow は、他の FIWARE GEs と相互作用して、次の2つのレイヤーで FIWARE ベースの IoT プラットフォームを強化できます。
上位層では、FogFlow は次の2つの方法で標準化された NGSI インターフェースを介して Orion Context Broker と相互作用できます。

Orion Broker とFogFlow を統合
======================================


Orion を起動
-------------------------------------------------------------

Orion のドキュメントに従って、Orion Context Broker インスタンスをここからセットアップできます: `Installing Orion`.

.. _`Installing Orion`: https://fiware-orion.readthedocs.io/en/master/admin/install/index.html


以下のコマンドを使用して、Docker で Orion をセットアップすることもできます (Docker はこの方法で必要です)
注: Orion コンテナーは MongoDB データベースに依存しています。

**前提条件:** Docker がインストールされている必要があります

最初に、以下のコマンドを使用して MongoDB コンテナーを起動します:

.. code-block:: console    

	sudo docker run --name mongodb -d mongo:3.4


そして、このコマンドで Orion を実行します:

.. code-block:: console    

	sudo docker run -d --name orion1 --link mongodb:mongodb -p 1026:1026 fiware/orion -dbhost mongodb


すべてが動作することを確認します:

.. code-block:: console    

	curl http://<Orion IP>:1026/version


Note: パブリックアクセスのためにファイアウォールのポート 1026 を許可します。



サブスクリプションを発行して、生成された結果を Orion Context Broker に転送
----------------------------------------------------------------------------------

次の curl リクエストを使用して、Fogflow Broker を FIWARE Orion にサブスクライブします:

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

ブラウザで次の URL にアクセスし、目的のコンテキストエンティティを検索します:

.. code-block:: console    

	curl http://<Orion IP>:1026/v2/entities/


最初の方法は、Orion Context Broker を FogFlow IoT サービスによって生成されるコンテキスト情報の宛先と見なすことです。この場合、NGSI サブスクリプションは、外部アプリケーションまたは FogFlow ダッシュボードによって発行され、要求されたコンテキスト更新を指定された Orion Context Broker に転送するように FogFlow に要求する必要があります。

.. figure:: ../../en/source/figures/orion-integration.png


どのエンティティを Orion Context Broker に転送するかを FogFlow に指示する方法は2つあります。最初の方法は、FogFlow Broker への生のサブスクリプションを発行することです。2番目の方法は、これを行うための小さな JavaScript プログラムを作成することです。以下に例を示します。統合は Orion Context Broker の NGSIv2 インターフェースを使用していることに注意してください。

.. important::

	* **fogflowBroker**: FogFlow Broker の IP アドレス。構成ファイルの "webportal_ip" または "coreservice_ip" の場合があります。これは、FogFlow システムにアクセスする場所までです。
	* **orionBroker**: 実行中の Orion インスタンスのアクセス可能な IP アドレス。


.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://fogflowBroker:8080/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -H 'Destination: orion-broker' \			
              -d '{"entities": [{"type": "PowerPanel", "isPattern": true}],
					"reference": "http://orionBroker:1026/v2/op/notify"} '           


   .. code-tab:: javascript

	    // please refer to the JavaScript library, located at  https://github.com/smartfog/fogflow/tree/master/designer/public/lib/ngsi
	
	    //  entityType: the type of context entities to be pushed to Orion Context Broker
	    //  orionBroker: the URL of your running Orion Context Broker
	    function subscribeFogFlow(entityType, orionBroker)
	    {
	        var subscribeCtxReq = {};    
	        subscribeCtxReq.entities = [{type: entityType, isPattern: true}];
	        subscribeCtxReq.reference =  'http://' + orionBroker + '/v2/op/notify';
	        
	        client.subscribeContext4Orion(subscribeCtxReq).then( function(subscriptionId) {
	            console.log(subscriptionId);   
	            ngsiproxy.reportSubID(subscriptionId);		
	        }).catch(function(error) {
	            console.log('failed to subscribe context');
	        });	
	    }
	    
	    
	    // client to interact with IoT Broker
	    var client = new NGSI10Client(config.brokerURL);
	    
	    subscribeFogFlow('PowerPanel', 'cpaasio-fogflow.inf.um.es:1026');
	

要求されたデータがOrion Context Broker にプッシュされているかどうかを確認するには、次の NGSIv2 クエリを送信して確認します。

.. code-block:: console 

    curl http://orionBroker:1026/v2/entities?type=PowerPanel -s -S -H 'Accept: application/json'
    

データソースとしての Orion Context Broker
---------------------------------------

2番目の方法は、Orion Context Broker を追加情報を提供するデータソースと見なすことです。この場合、簡単なフォグ ファンクションを実装して、必要な情報を FogFlow システムにフェッチできます。いずれの場合も、既存の Orion ベースの FIWARE システムに変更を加える必要はありません。したがって、このタイプの統合は、ほとんど労力をかけずに高速に実行できます。

下位層では、MQTT、COAP、OneM2M、OPC-UA、LoRaWAN などの NGSI 以外でサポートされているデバイスと統合するために、FogFlow は既存の IoT Agents のモジュールを再利用し、フォグ ファンクション プログラミング モデルに基づいて FogFlow アダプターに変換できます。これらのアダプターを使用すると、FogFlow は、デバイス統合に必要なアダプターをエッジで直接起動できます。このようにして、FogFlow はさまざまな IoT デバイスと通信できます。

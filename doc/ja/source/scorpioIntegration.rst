*****************************************
FogFlow を NGSI-LD Broker と統合
*****************************************

チュートリアルでは、FogFlow を高度なデータ分析フレームワークとして利用して、Scorpio、Orion-LD、Stellio などの NGSI-LD Broker
でキャプチャされた生データに加えてオンデマンド データ分析を可能にする方法を紹介します。次の図は、これを行う方法の簡単な例を
詳細に示しています。主に、8つのステップを持つ3つの側面が含まれています。

* NGSI-LD Broker から FogFlow システムに生データをフェッチする方法 (**ステップ1-3**)
* FogFlow のサーバーレス機能を使用してカスタマイズされたデータ分析を行う方法 (**ステップ4**)
* 生成された分析結果を NGSI-LD Broker にプッシュして、さらに共有する方法 (**ステップ5-8**)

.. figure:: ../../en/source/figures/fogflow-ngsild-broker.png


詳細な手順を検討する前に、以下の情報に従って FogFlow システムとNGSI-LD Broker をセットアップしてください。

まず、`FogFlow on a Single Machine`_ を参照 して、単一のホストマシン上に FogFlow システムをセットアップしてください。

NGSI-LD Broker に関しては、Scorpio、Orion-LD、Stellio のさまざまな選択肢があります。ここでは、詳細な手順を示す具体的な例として
Orion-LD を取り上げます。同じホストマシンに Orion-LD ブローカーをセットアップするには、次の手順を参照してください。
他のブローカー (Scorpio、Stellio など)との統合では、リクエストのポート番号と構成ファイルを少し変更する必要があります。

.. code-block:: console

	# fetch the docker-compose file 
	wget https://raw.githubusercontent.com/smartfog/fogflow/development/test/orion-ld/docker-compose.yml

	# start the orion-ld broker
	docker-compose pull
	docker-Compose up -d 

次の手順を開始する前に、Orion-LD Broker と FogFlow システムが正しく実行されているかどうかを確認してください。

# Orion-LD Broker が実行されているかどうかを確認します

.. code-block:: console

	curl localhost:1026/ngsi-ld/ex/v1/version

# FogFlow システムが正しく実行されているかどうかを確認します

	ブラウザから FogFlow ダッシュボードを開きます

Orion-LD から FogFlow にデータをフェッチする方法
================================================================

ステップ 1: Orion-LD へのサブスクリプションを発行します

.. code-block:: console

	curl -iX POST \
		  'http://localhost:1026/ngsi-ld/v1/subscriptions' \
		  -H 'Content-Type: application/json' \
		  -H 'Accept: application/ld+json' \
		  -H 'Link: <https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
		  -d ' {
             	"type": "Subscription",
             	"entities": [{
                   "type": "Vehicle"
             	}],
             	"notification": {
                   "format": "keyValues",
                   "endpoint": {
                       "uri": "http://192.168.0.59:8070/ngsi-ld/v1/notifyContext/",
                       "accept": "application/ld+json"
             	    }
            	}
 			}'


ステップ 2: エンティティの更新を Orion-LD に送信します


.. code-block:: console    


	curl --location --request POST 'http://localhost:1026/ngsi-ld/v1/entityOperations/upsert?options=update' \
	--header 'Content-Type: application/json' \
	--header 'Link: <https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	--data-raw '[
	{
	   "id": "urn:ngsi-ld:Vehicle:A106",
	   "type": "Vehicle",
	   "brandName": {
	                  "type": "Property",
	                  "value": "Mercedes"
	    },
	    "isParked": {
	                  "type": "Relationship",
	                  "object": "urn:ngsi-ld:OffStreetParking:Downtown1",
	                  "providedBy": {
	                                  "type": "Relationship",
	                                  "object": "urn:ngsi-ld:Person:Bob"
	                   }
	     },
	     "speed": {
	                "type": "Property",
	                "value": 120
	      },
	     "location": {
	                    "type": "GeoProperty",
	                    "value": {
	                              "type": "Point",
	                              "coordinates": [-8.5, 41.2]
	                    }
	     }
	}
	]'



ステップ 3: FogFlow がサブスクライブされたエンティティを受信するかどうかを確認します


FogFlow thinBroker から "Vehicle" エンティティをクエリするための curl コマンドを準備してください。


.. code-block:: console    

	curl -iX GET \
		  'http://localhost:8070/ngsi-ld/v1/entities?type=Vehicle' \
		  -H 'Content-Type: application/ld+json' \
		  -H 'Accept: application/ld+json' \
		  -H 'Link: <https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' 



データ分析機能をプログラムして適用する方法
================================================================

ステップ 4: fogfunction1 を適用して、カスタマイズされたデータ分析を実行します


"/application/operator/alert" のコードを変更して、簡単な分析を行ってください。たとえば、
車両の速度がしきい値を超えたときにアラート メッセージを生成します。


生成された結果を NGSI-LD Broker にプッシュ バックする方法
=========================================================================

オペレーター、Docker イメージ、Fog ファンクションを含む fogfunction2 をデフォルトで登録してください。
それらをデザイナの初期化リストに入れます。


ステップ 5: fogfunction2 をトリガーする更新メッセージを送信します


.. code-block:: console    

	#please write the curl message to trigger fogfunction2


ステップ 6: fogfunction2 が作成されているかどうかを確認します


fogfunction2 がトリガーされているかどうかをユーザーが確認できる場所を説明します。


ステップ 7: Orion-LD が転送された結果を受信したかどうかを確認します


.. code-block:: console    

	curl -iX GET \
		  'http://localhost:8070/ngsi-ld/v1/entities?type=Alert' \
		  -H 'Content-Type: application/ld+json' \
		  -H 'Accept: application/ld+json' \
		  -H 'Link: <https://fiware.github.io/data-models/context.jsonld>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' 

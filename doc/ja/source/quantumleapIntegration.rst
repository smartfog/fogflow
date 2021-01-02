*****************************************
FogFlow を QuantumLeap と統合
*****************************************

`QuantumLeap`_ は、NGSIv2 の時空間データを保存、クエリ、取得するための REST サービスです。QuantumLeap は、NGSI の半構造化データを表形式に変換し、時系列データベースに保存します。

.. _`QuantumLeap`: https://quantumleap.readthedocs.io/en/latest/

次の図は、fogflow と QuantumLeap の統合を示しています。

.. figure:: ../../en/source/figures/quantum-leap-fogflow-integration.png

1. ユーザーは、NGSIv2 の FogFlow Broker にサブスクリプション要求を送信します。
2. ユーザーは、NGSIv1 の FogFlow Broker にアップデート要求を送信します。
3. FogFlow Broker は、NGSIv2 の QuantumLeap にノーティファイします。

統合手順
===============================================

**前提条件:**

* FogFlow は、少なくとも1つのノードで稼働している必要があります。
* QuantumLeap は、少なくとも1つのノードで稼働している必要があります。

FogFlow Broker に QuantumLeap の **サブスクリプション リクエストを送信します**:

.. code-block:: console

	curl -iX POST \
	'http://<FogFlow Broker>:8070/v2/subscriptions' \
	 -H 'Content-Type: application/json' \
	 -d '
 	{
	"description": "A subscription to get info about Room1",
	"subject": {
		"entities": [{
			"id": "Room4",
			"type": "Room",
			"isPattern": false
		}],
		"condition": {
			"attrs": [
				"temperature"
			]
		}
	},
	"notification": {
		"http": {
			"url": "http://<Quantum-leap-Host-IP>:8668/v2/notify"
		},
		"attrs": [
			"temperature"
		]
	},
	"expires": "2040-01-01T14:00:00.00Z",
	"throttling": 5
    }'

上記のサブスクリプションで定義されたタイプと属性のエンティティを使用して、FogFlow Broker に **更新リクエストを送信します**。
リクエストの例を以下に示します:

.. code-block:: console

	curl -iX POST \
  	'http://<Fogflow broker>:8070/ngsi10/updateContext' \
 	 -H 'Content-Type: application/json' \
  	-d '
      {
	"contextElements": [{
		"entityId": {
			"id": "Room4",
			"type": "Room",
			"isPattern": false
		},
		"attributes": [{
			"name": "temperature",
			"type": "Integer",
			"value": 155
		}],
		"domainMetadata": [{
			"name": "location",
			"type": "point",
			"value": {
				"latitude": 49.406393,
				"longitude": 8.684208
			}
		}]
	}],
	"updateAction": "UPDATE"
     }'

**FogFlow Broker** は Quantumleap にノーティフィケーションを送信します。以下のコマンドで結果を確認します:

.. code-block:: console

	http://<QuantuLeap-Host-Ip>:8668/v2/entities/Room4/attrs/temperature

**結果:**

.. figure:: ../../en/source/figures/quantum-leap-result.png

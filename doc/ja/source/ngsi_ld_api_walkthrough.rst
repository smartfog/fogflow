*****************************************
NGSI-LD API ウォークスルー
*****************************************

このチュートリアルは、主に FogFlow でサポートされている NGSI-LD APIs に焦点を当てています。これには、エンティティ、コンテキスト レジストレーション、およびサブスクリプションの API が含まれます。これらについては、次のセクションで詳しく説明します。FogFlow で NGSI-LD APIs を使用するには、Docker Hub から最新の Docker イメージ "fogflow/broker:3.1" をチェック アウトしてください。

FogFlow は、NGSI-LD データモデルに従い、継続的に改善されています。NGSI-LD データモデルの理解を深めるには、こちら (`this`_) を参照してください。

.. _`this`: https://fiware-datamodels.readthedocs.io/en/latest/ngsi-ld_howto/index.html


エンティティ (Entities)
=========================

エンティティは、環境内のオブジェクトを表すための単位であり、それぞれに独自のプロパティがあり、他のオブジェクトとのリレーションシップもあります。これがリンクトデータの形成方法です。


エンティティを作成 (Create entities)
------------------------------------------

FogFlow Broker で NGSI-LD エンティティを作成する方法はいくつかあります:

* Link ヘッダーでコンテキストが提供されている場合: ペイロードを解決するためのコンテキストは、Link ヘッダーを介して提供されます。
* ペイロードでコンテキストが提供されている場合: コンテキストはペイロード自体にあるため、リクエストに Link ヘッダーを添付する必要はありません。
* リクエストペイロードがすでに拡張されている場合: 一部のペイロードは、何らかのコンテキストを使用してすでに拡張されています。

さまざまな方法で FogFlow Broker でエンティティを作成するための curl リクエストを以下に示します。これらはすべて、FogFlow Broker への POST リクエストです。ブローカーは、新しいエンティティの作成が成功した場合は "201 Created"、既存のエンティティの作成が成功した場合は "409 Conflict" の応答を返します。

**Link ヘッダーでコンテキストが提供されている場合:**

.. code-block:: console

        curl -iX POST \
        'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/' \
        -H 'Content-Type: application/ld+json' \
        -H 'Accept: application/ld+json' \
        -H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
        -d '
        {
                "id": "urn:ngsi-ld:Vehicle:A100",
                "type": "Vehicle",
                "brandName": {
                        "type": "Property",
                        "value": "BMW",
                        "observedAt": "2017-07-29T12:00:04"
                },
                "isParked": {
                        "type": "Relationship",
                        "object": "urn:ngsi-ld:OffStreetParking:Downtown",
                        "observedAt": "2017-07-29T12:00:04",
                        "providedBy": {
                                "type": "Relationship",
                                "object": "urn:ngsi-ld:Person:Bob"
                        }
                },
                "speed": {
                        "type": "Property",
                        "value": 81,
                        "observedAt": "2017-07-29T12:00:04"
                },
                "location": {
                        "type": "GeoProperty",
                        "value": {
                                "type": "Point",
                                "coordinates": [-8.5, 41.2]
                        }
                }
        }'

**ペイロードでコンテキストが提供される場合:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-d '
	{
		"@context": [{
			"Vehicle": "https://uri.etsi.org/ngsi-ld/default-context/Vehicle",
			"brandName": "https://uri.etsi.org/ngsi-ld/default-context/brandName",
			"speed": "https://uri.etsi.org/ngsi-ld/default-context/speed",
			"isParked": {
				"@type": "@id",
				"@id": "https://uri.etsi.org/ngsi-ld/default-context/isParked"
			}
		}],
		"id": "urn:ngsi-ld:Vehicle:A200",
		"type": "Vehicle1",
		"brandName": {
			"type": "Property",
			"value": "Mercedes"
		},
		"isParked": {
			"type": "Relationship",
			"object": "urn:ngsi-ld:OffStreetParking:Downtown1",
			"observedAt": "2017-07-29T12:00:04",
			"providedBy": {
				"type": "Relationship",
				"object": "urn:ngsi-ld:Person:Bob"
			}
		},
		"speed": {
			"type": "Property",
			"value": 80
		},
		"createdAt": "2017-07-29T12:00:04",
		"location": {
			"type": "GeoProperty",
			"value": {
				"type": "Point",
				"coordinates": [-8.5, 41.2]
			}
		}
	}'

**リクエスト ペイロードがすでに拡張されている場合:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-d '
	{
		"https://uri.etsi.org/ngsi-ld/default-context/brandName": [
			{
				"@type": [
					"https://uri.etsi.org/ngsi-ld/Property"
				],
				"https://uri.etsi.org/ngsi-ld/hasValue": [
					{
						"@value": "Mercedes"
					}
				]
			}
		],
		"https://uri.etsi.org/ngsi-ld/createdAt": [
			{
				"@type": "https://uri.etsi.org/ngsi-ld/DateTime",
				"@value": "2017-07-29T12:00:04"
			}
		],
		"@id": "urn:ngsi-ld:Vehicle:A300",
		"https://uri.etsi.org/ngsi-ld/default-context/isParked": [
			{
				"https://uri.etsi.org/ngsi-ld/hasObject": [
					{
						"@id": "urn:ngsi-ld:OffStreetParking:Downtown1"
					}
				],
				"https://uri.etsi.org/ngsi-ld/observedAt": [
					{
						"@type": "https://uri.etsi.org/ngsi-ld/DateTime",
						"@value": "2017-07-29T12:00:04"
					}
				],
				"https://uri.etsi.org/ngsi-ld/default-context/providedBy": [
					{
						"https://uri.etsi.org/ngsi-ld/hasObject": [
							{
								"@id": "urn:ngsi-ld:Person:Bob"
							}
						],
						"@type": [
							"https://uri.etsi.org/ngsi-ld/Relationship"
						]
					}
				],
				"@type": [
					"https://uri.etsi.org/ngsi-ld/Relationship"
				]
			}
		],
		"https://uri.etsi.org/ngsi-ld/location": [
			{
				"@type": [
					"https://uri.etsi.org/ngsi-ld/GeoProperty"
				],
				"https://uri.etsi.org/ngsi-ld/hasValue": [
					{
						"@value": "{ \"type\":\"Point\", \"coordinates\":[ -8.5, 41.2 ] }"
					}
				]
			}
		],
		"https://uri.etsi.org/ngsi-ld/default-context/speed": [
			{
				"@type": [
					"https://uri.etsi.org/ngsi-ld/Property"
				],
				"https://uri.etsi.org/ngsi-ld/hasValue": [
					{
						"@value": 80
					}
				]
			}
		],
		"@type": [
			"https://uri.etsi.org/ngsi-ld/default-context/Vehicle"
		]
	}'


エンティティを更新 (Update entities)
-----------------------------------------------

エンティティは、属性 (プロパティとリレーションシップ) をアップデートすることで更新でき、属性は次の方法で更新できます。

* エンティティに属性を追加: 既存のエンティティに、プロパティまたはリレーションシップ、あるいはその両方を追加できます。これは、エンティティに属性を追加するためのブローカーへの POST http リクエストです。
* エンティティの既存の属性をアップデート: エンティティの既存のプロパティまたはリレーションシップ、あるいはその両方を更新できます。これは、FogFlow Broker への PATCH http リクエストです。
* エンティティの特定の属性をアップデート: エンティティの既存の属性のフィールドを更新できます。この更新は、部分アップデート (partial update) とも呼ばれます。これは、FogFlow Broker への PATCH リクエストでもあります。

FogFlow Broker は、属性の更新が成功すると "204 NoContent" を返し、存在しないエンティティの場合は "404 NotFound" を返します。既存のエンティティの属性をアップデートしているときに、リクエスト ペイロードで提供された属性の一部が存在しない可能性があります。このような場合、FogFlow Broker は "207 MultiStatus" エラーを返します。

これらのアップデートの curl リクエストは次のとおりです。

**エンティティに属性を追加:**

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs' \
	-H 'Content-Type: application/ld+json' \
	-d '
	{
		"@context": {
			"brandName1": "https://uri.etsi.org/ngsi-ld/default-context/brandName1",
			"isParked1": "https://uri.etsi.org/ngsi-ld/default-context/isParked1"
		},
		"brandName1": {
			"type": "Property",
			"value": "Audi"
		},
		
		"isParked1": {
			"type": "Relationship",
			"object": "Audi"
		}
	}'

**エンティティの既存の属性をアップデート:**

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs' \
	-H 'Content-Type: application/ld+json' \
	-d '
	{
		"@context": {
			"isParked": "https://uri.etsi.org/ngsi-ld/default-context/isParked"
		},
		"brandName": {
			"type": "Property",
			"object": "Audi"
		}
	}'

**エンティティの特定の属性をアップデート:**

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs/<Attribute-Name>' \
	-H 'Content-Type: application/ld+json' \
	-d '
		{
		"@context": {
			"brandName": "https://uri.etsi.org/ngsi-ld/default-context/brandName"
		},
		"value": "Suzuki"
	}'


エンティティを取得 (Get entities)
-----------------------------------------------

このセクションでは、FogFlow Broker から作成済みのエンティティを取得する方法について説明します。エンティティは、以下にリストされているさまざまなフィルターに基づいて FogFlow から取得できます。

* エンティティ Id に基づく: リクエスト URL で id が渡されたエンティティを返します。
* 属性名に基づく: リクエスト URL のクエリ パラメーターで渡される属性名を含むすべてのエンティティを返します。
* エンティティ Id とエンティティ タイプに基づく: クエリ パラメーターで指定されたものと同じエンティティ Id を持つエンティティとタイプの一致を返します。
* エンティティ タイプに基づく: 要求されたタイプのすべてのエンティティを返します。
* Link ヘッダー付きのエンティティ タイプに基づく: 要求されたタイプのすべてのエンティティを返しますが、ここでは、要求 URL のクエリ パラメーターでタイプを別の方法で指定できます。次のセクションでこの要求を参照してください。
* エンティティ IdPattern とエンティティ タイプに基づく: IdPattern 範囲内にあるすべてのエンティティと、クエリ パラメーターに記載されている一致するタイプを返します。

上記のリクエストで少なくとも1つのエンティティが正常に取得されると、FogFlow Brokerは "200 OK" 応答を返します。存在しないエンティティの場合、"404 NotFound" エラーが返されます。


**エンティティ Id に基づく:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**属性名に基づく:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?attrs=<Expanded-Attribute-Name>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**エンティティ Id とエンティティ タイプに基づく:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?id=<Entity-Id>&type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**エンティティ タイプに基づく:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'

**Link ヘッダー付きのエンティティ タイプに基づく:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?type=<Unexpanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"'

**エンティティ IdPattern およびエンティティ タイプに基づく:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities?idPattern=<Entity-IdPattern>&type=<Expanded-Entity-Type>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json'


エンティティを削除 (Delete entities)
-----------------------------------------------

エンティティを削除するか、そのエンティティの特定の属性を削除することができます。削除に成功すると、"204 NoContent" 応答が返されますが、存在しない属性またはエンティティの場合は、"404 NotFound" エラーが返されます。

**エンティティの特定の属性の削除:**

.. code-block:: console

	curl -iX DELETE \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>/attrs/<Attribute-Name>'

**エンティティの削除:**

.. code-block:: console

	curl -iX DELETE \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/entities/<Entity-Id>'


サブスクリプション (Subscriptions)
==================================

サブスクライバーは、FogFlow Broker へのサブスクリプション要求を使用してエンティティをサブスクライブできます。


サブスクリプションを作成 (Create subscriptions)
------------------------------------------------

サブスクリプションは、エンティティ Id またはエンティティ IdPattern のいずれかに対して作成できます。そのサブスクリプションに対してエンティティのアップデートがある場合は常に、FogFlow Broker はアップデートされたエンティティをサブスクライバーに自動的に通知します。"201 Created" 応答は、Broker でサブスクリプションが成功すると、サブスクリプション Id とともに返されます。これは、後でサブスクリプションを取得、アップデート、または削除するために使用できます。

次の curl リクエストを参照してください。ただし、サブスクリプションを実行する前に、ノーティフィケーションの内容を簡単に表示できるノーティファイ レシーバーが実行されていることを確認してください。すでにサブスクライブしているエンティティの場合、エンティティの作成またはアップデートが行われると、サブスクライバーがノーティフィケーションを受信します。既存のエンティティへのサブスクリプションの場合、サブスクライバーもノーティフィケーションを受け取ります。

**エンティティ Id のサブスクライブ**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	-d '
	{
		"type": "Subscription",
		"entities": [{
			"id" : "urn:ngsi-ld:Vehicle:A100",
			"type": "Vehicle"
		}],
		"watchedAttributes": ["brandName"],
		"notification": {
			"attributes": ["brandName"],
			"format": "keyValues",
			"endpoint": {
				"uri": "http://my.endpoint.org/notify",
				"accept": "application/json"
			}
		}
	}'

**IdPattern のサブスクライブ:**

.. code-block:: console

	curl -iX POST \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	-d '
	{
		"type": "Subscription",
		"entities": [{
			"idPattern" : ".*",
			"type": "Vehicle"
		}],
		"watchedAttributes": ["brandName"],
		"notification": {
			"attributes": ["brandName"],
			"format": "keyValues",
			"endpoint": {
				"uri": "http://my.endpoint.org/notify",
				"accept": "application/json"
			}
		}
	}'


サブスクリプションを更新 (Update subscriptions)
------------------------------------------------

FogFlow Broker の既存のサブスクリプションは、以下の curl リクエストを使用してidで更新できます。

.. code-block:: console

	curl -iX PATCH \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/<Subscription-Id>' \
	-H 'Content-Type: application/ld+json' \
	-H 'Accept: application/ld+json' \
	-H 'Link: <{{link}}>; rel="https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld"; type="application/ld+json"' \
	-d '
	{
		"type": "Subscription",
		"entities": [{
			"type": "Vehicle1"
		}],
		"watchedAttributes": ["https://uri.etsi.org/ngsi-ld/default-context/brandName11"],
		"notification": {
			"attributes": ["https://uri.etsi.org/ngsi-ld/default-context/brandName223"],
			"format": "keyValues",
			"endpoint": {
				"uri": "http://my.endpoint.org/notify",		
				"accept": "application/json"
			}
		}
	}'
	

サブスクリプションを取得 (Get subscriptions)
---------------------------------------------

すべてのサブスクリプションまたは特定の Id を持つサブスクリプションは、どちらも "200 OK" の応答で FogFlow Broker から取得できます。curl のリクエストは以下のとおりです。

**すべてのサブスクリプション:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/' \
	-H 'Accept: application/ld+json'

**特定のサブスクリプション:**

.. code-block:: console

	curl -iX GET \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/<Subscription-Id>' \
	-H 'Accept: application/ld+json'


サブスクリプションを削除 (Delete subscriptions)
-----------------------------------------------

サブスクリプションは、"204 NoContent "の応答で FogFlow Broker に次のリクエストを送信することで削除できます。

.. code-block:: console

	curl -iX DELETE \
	'http://<Thin_Broker_IP>:8070/ngsi-ld/v1/subscriptions/<Subscription-Id>'



**FogFlow での NGSI-LD のサポートにも、いくつかの制限があります。改善は続けられています。**

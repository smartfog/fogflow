*****************************************
APIs とその使用例
*****************************************

FogFlow Discovery の APIs
===================================

近くのブローカーを探す
-----------------------------------------------

外部アプリケーションまたは IoT デバイスの場合、FogFlow Discovery から必要な唯一のインターフェースは、自身の場所に基づいて近くのブローカーを見つけることです。その後、彼らは割り当てられた近くのブローカーと対話する必要があるだけです。


**POST /ngsi9/discoverContextAvailability**

==============   ===============
パラメーター     説明
==============   ===============
latitude         場所 (location) の緯度
longitude        場所 (location)の緯度
limit            予想されるブローカーの数
==============   ===============


以下の例を確認してください。

.. note:: JavaScript のコード例では、ライブラリ ngsiclient.js が必要です。application/device/powerpanel のコード リポジトリを参照してください。

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8071/ngsi9/discoverContextAvailability' \
              -H 'Content-Type: application/json' \
              -d '
                {
                   "entities":[
                      {
                         "type":"IoTBroker",
                         "isPattern":true
                      }
                   ],
                   "restriction":{
                      "scopes":[
                         {
                            "scopeType":"nearby",
                            "scopeValue":{
                               "latitude":35.692221,
                               "longitude":139.709059,
                               "limit":1
                            }
                         }
                      ]
                   }
                }             


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        
        var discoveryURL = "http://localhost:8071/ngsi9";
        var myLocation = {
                "latitude": 35.692221,
                "longitude": 139.709059
            };
        
        // find out the nearby IoT Broker according to my location
        var discovery = new NGSI.NGSI9Client(discoveryURL)
        discovery.findNearbyIoTBroker(myLocation, 1).then( function(brokers) {
            console.log('-------nearbybroker----------');    
            console.log(brokers);    
            console.log('------------end-----------');    
        }).catch(function(error) {
            console.log(error);
        });

  
       

FogFlow Broker の APIs
===============================

コンテキストの作成/アップデート
-----------------------------------------------

.. note:: コンテキスト エンティティを作成またはアップデートするのと同じ API です。コンテキスト アップデートの場合、既存のエンティティがない場合は、新しいエンティティが作成されます。


**POST /ngsi10/updateContext**

==============   ===============
パラメーター     説明
==============   ===============
latitude         場所 (location) の緯度
longitude        場所 (location)の緯度
limit            予想されるブローカーの数
==============   ===============

例: 

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -iX POST \
              'http://localhost:8070/ngsi10/updateContext' \
              -H 'Content-Type: application/json' \
              -d '
                {
                    "contextElements": [
                        {
                            "entityId": {
                                "id": "Device.temp001",
                                "type": "Temperature",
                                "isPattern": false
                            },
                            "attributes": [
                            {
                              "name": "temp",
                              "type": "integer",
                              "value": 10
                            }
                            ],
                            "domainMetadata": [
                            {
                                "name": "location",
                                "type": "point",
                                "value": {
                                    "latitude": 49.406393,
                                    "longitude": 8.684208
                                }
                            },{
                                "name": "city",
                                "type": "string",
                                "value": "Heidelberg"                             
                            }
                            ]
                        }
                    ],
                    "updateAction": "UPDATE"
                }'          


   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"
    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
    
        var profile = {
                "type": "PowerPanel",
                "id": "01"};
        
        var ctxObj = {};
        ctxObj.entityId = {
            id: 'Device.' + profile.type + '.' + profile.id,
            type: profile.type,
            isPattern: false
        };
        
        ctxObj.attributes = {};
        
        var degree = Math.floor((Math.random() * 100) + 1);        
        ctxObj.attributes.usage = {
            type: 'integer',
            value: degree
        };   
        ctxObj.attributes.shop = {
            type: 'string',
            value: profile.id
        };       
        ctxObj.attributes.iconURL = {
            type: 'string',
            value: profile.iconURL
        };                   
        
        ctxObj.metadata = {};
        
        ctxObj.metadata.location = {
            type: 'point',
            value: profile.location
        };    
       
        ngsi10client.updateContext(ctxObj).then( function(data) {
            console.log(data);
        }).catch(function(error) {
            console.log('failed to update context');
        }); 



GET を介してコンテキストをクエリ
-----------------------------------------------


ID でコンテキスト エンティティを取得
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi10/entity/#eid**

==============   ===============
パラメーター     説明
==============   ===============
eid              エンティティ ID
==============   ===============

例: 

.. code-block:: console 

   curl http://localhost:8070/ngsi10/entity/Device.temp001

特定のコンテキスト エンティティの特定の属性を取得
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi10/entity/#eid/#attr**

==============   ===============
パラメーター     説明
==============   ===============
eid              エンティティ ID
attr             フェッチする属性名を指定
==============   ===============

例: 

.. code-block:: console 

   curl http://localhost:8070/ngsi10/entity/Device.temp001/temp


単一のブローカー上のすべてのコンテキスト エンティティを確認
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**GET /ngsi10/entity**

例: 

.. code-block:: console 

    curl http://localhost:8070/ngsi10/entity



POST を介してコンテキストをクエリ
-----------------------------------------------

**POST /ngsi10/queryContext**

==============   ===============
パラメーター     説明
==============   ===============
entityId         特定のエンティティ ID、ID pattern、またはタイプを定義できるエンティティ フィルターを指定
restriction      スコープのリストと各スコープは、ドメイン メタデータに基づいてフィルターを定義
==============   ===============

エンティティ ID pattern によるコンテキストのクエリ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{"entities":[{"id":"Device.*","isPattern":true}]}'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{id:'Device.*', isPattern: true}];           
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });          


エンティティ タイプによるコンテキストのクエリ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{"entities":[{"type":"Temperature","isPattern":true}]}'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'Temperature', isPattern: true}];           
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });          


ジオスコープ (geo-scope) によるコンテキストのクエリ (サークル)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]};
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


ジオスコープ (geo-scope) によるコンテキストのクエリ (ポリゴン)
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
               "entities":[
                  {
                     "id":".*",
                     "isPattern":true
                  }
               ],
               "restriction":{
                  "scopes":[
                     {
                        "scopeType":"polygon",
                        "scopeValue":{
                           "vertices":[
                              {
                                 "latitude":34.4069096565206,
                                 "longitude":135.84594726562503
                              },
                              {
                                 "latitude":37.18657859524883,
                                 "longitude":135.84594726562503
                              },
                              {
                                 "latitude":37.18657859524883,
                                 "longitude":141.51489257812503
                              },
                              {
                                 "latitude":34.4069096565206,
                                 "longitude":141.51489257812503
                              },
                              {
                                 "latitude":34.4069096565206,
                                 "longitude":135.84594726562503
                              }
                           ]
                        }
                    }]
                }
            }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {
               "scopes":[
                  {
                     "scopeType":"polygon",
                     "scopeValue":{
                        "vertices":[
                           {
                              "latitude":34.4069096565206,
                              "longitude":135.84594726562503
                           },
                           {
                              "latitude":37.18657859524883,
                              "longitude":135.84594726562503
                           },
                           {
                              "latitude":37.18657859524883,
                              "longitude":141.51489257812503
                           },
                           {
                              "latitude":34.4069096565206,
                              "longitude":141.51489257812503
                           },
                           {
                              "latitude":34.4069096565206,
                              "longitude":135.84594726562503
                           }
                        ]
                     }
                  }
               ]
            }
                    
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


ドメイン メタデータ値のフィルターを使用したコンテキストのクエリ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. note:: 条件ステートメント (conditional statement) は、コンテキスト エンティティのドメイン メタデータでのみ定義できます。当面は、特定の属性値に基づいてエンティティを除外することはサポートされていません。

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "restriction": {
                        "scopes": [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {scopes: [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};        
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


複数のフィルターを使用したコンテキストのクエリ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/queryContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var queryReq = {}
        queryReq.entities = [{type:'.*', isPattern: true}];  
        queryReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};          
        
        ngsi10client.queryContext(queryReq).then( function(deviceList) {
            console.log(deviceList);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to query context');
        });    


コンテキストを削除
-----------------------------------------------

ID で特定のコンテキスト エンティティを削除
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /ngsi10/entity/#eid**

==============   ===============
パラメーター     説明
==============   ===============
eid              エンティティ ID
==============   ===============

例: 

.. code-block:: console 

    curl -iX DELETE http://localhost:8070/ngsi10/entity/Device.temp001





コンテキストをサブスクライブ
-----------------------------------------------

**POST /ngsi10/subscribeContext**

==============   ===============
パラメーター     説明
==============   ===============
entityId         特定のエンティティ ID、ID pattern、またはタイプを定義できるエンティティ フィルターを指定
restriction      スコープのリストと各スコープは、ドメイン メタデータに基づいてフィルターを定義
reference        ノーティフィケーションを受信する宛先
==============   ===============

エンティティ ID のパターンによるコンテキストのサブスクライブ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities":[{"id":"Device.*","isPattern":true}],
                    "reference": "http://localhost:8066"
                }'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        var mySubscriptionId;
        
        var subscribeReq = {}
        subscribeReq.entities = [{id:'Device.*', isPattern: true}];           
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });

エンティティ タイプごとのコンテキストのサブスクライブ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities":[{"type":"Temperature","isPattern":true}]
                    "reference": "http://localhost:8066"                    
                  }'          

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'Temperature', isPattern: true}];           
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });       


ジオスコープ (geo-scope) によるコンテキストのサブスクライブ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "reference": "http://localhost:8066",                    
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'.*', isPattern: true}];  
        subscribeReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            }
                        }]};
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });   

ドメイン メタデータ値のフィルターを使用してコンテキストをサブスクライブ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. note:: 条件ステートメント (conditional statement) は、コンテキスト エンティティのドメイン メタデータでのみ定義できます。当面は、特定の属性値に基づいてエンティティを除外することはサポートされていません。

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "reference": "http://localhost:8066",                    
                    "restriction": {
                        "scopes": [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'.*', isPattern: true}];  
        subscribeReq.restriction = {scopes: [{
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};        
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });      


複数のフィルターでコンテキストをサブスクライブ
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. tabs::

   .. group-tab:: curl

        .. code-block:: console 

            curl -X POST 'http://localhost:8070/ngsi10/subscribeContext' \
              -H 'Content-Type: application/json' \
              -d '{
                    "entities": [{
                        "id": ".*",
                        "isPattern": true
                    }],
                    "reference": "http://localhost:8066", 
                    "restriction": {
                        "scopes": [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]
                    }
                  }'
                  

   .. code-tab:: javascript

        const NGSI = require('./ngsi/ngsiclient.js');
        var brokerURL = "http://localhost:8070/ngsi10"    
        var ngsi10client = new NGSI.NGSI10Client(brokerURL);
        
        var subscribeReq = {}
        subscribeReq.entities = [{type:'.*', isPattern: true}];  
        subscribeReq.restriction = {scopes: [{
                            "scopeType": "circle",
                            "scopeValue": {
                               "centerLatitude": 49.406393,
                               "centerLongitude": 8.684208,
                               "radius": 10.0
                            } 
                        }, {
                            "scopeType": "stringQuery",
                            "scopeValue":"city=Heidelberg" 
                        }]};          
        
        // use the IP and Port number your receiver is listening
        subscribeReq.reference =  'http://' + agentIP + ':' + agentPort;  
        
        
        ngsi10client.subscribeContext(subscribeReq).then( function(subscriptionId) {		
            console.log("subscription id = " + subscriptionId);   
    		mySubscriptionId = subscriptionId;
        }).catch(function(error) {
            console.log('failed to subscribe context');
        });   

サブスクリプション ID でサブスクリプションをキャンセル
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

**DELETE /ngsi10/subscription/#sid**


==============   ===============
パラメーター     説明
==============   ===============
sid              サブスクリプションの発行時に作成されるサブスクリプション ID
==============   ===============


curl -iX DELETE http://localhost:8070/ngsi10/subscription/#sid


FogFlow Orchestrator の APIs
=======================================

オペレーターを登録
-----------------------------------------------


サービス トポロジーを登録
-----------------------------------------------


サービス トポロジーをトリガー
-----------------------------------------------


サービス トポロジーの実行中のタスクをクエリ
----------------------------------------------------


フォグ ファンクションを登録
-----------------------------------------------


サービス トポロジーをトリガー
-----------------------------------------------


サービス トポロジーの実行中のタスクをクエリ
----------------------------------------------------


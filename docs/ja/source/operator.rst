.. _operator-implementation:

Docker 化されたオペレーターを実装する方法
=========================================

:ref:`flow-task`, で説明されているように、FogFlow は、データ処理タスクに代わって必要な入力ストリームを自動的にサブスクライブできます。オペレーター開発者が注目すべきは、内部データ処理ロジックです。

利用可能なテンプレート
-----------------------

現在、FogFlow は、JavaScript (Node.js ベース) および Python のオペレーターを実装するためのテンプレートを提供しています。独自のオペレーターを開発するときにそれらを再利用できます。

	* `nodejs template <https://github.com/smartfog/fogflow/tree/master/application/operator/template/javascript/>`_.
	* `python template <https://github.com/smartfog/fogflow/tree/master/application/operator/template/python/>`_. 

以下では、提供されている Node.js テンプレートに基づいてオペレーターを実装する方法について説明します。

オペレーターをプログラムする方法
-------------------------------

Node.js テンプレートをダウンロードすると、フレームワークに従って独自のアプリケーションロジックを実装できます。全体として、次の重要なことを知っておく必要があります:

* 提供されたテンプレートの主要なロジック部分は何ですか？

        次のコード ブロックは、提供されている Node.js テンプレートの主要なロジック部分を示しています。ngsiclient.js と ngsiagent.js の2つのライブラリを使用します。ngsiclient.js は、Node.js ベースの NGSI クライアントであり、NGSI10 を介して FogFlow Broker と対話します。ngsiagent.js は、NGSI10 ノーティファイを受信するサーバーです。提供されているライブラリを使用すると、アプリケーションで NGSI10 のアップデート、サブスクライブ、クエリを実行できます。
	
        メインロジックでは、3つの重要なファンクションが呼び出されます。各機能の役割は以下のとおりです。
	
	* startApp: アプリケーションの初期化を行います。このファンクションは、ngsiagent が入力データの受信をリッスンし始めた直後に呼び出されます
	* handleAdmin: FogFlow によって発行された構成コマンド (configuration commands) を処理します。
	* handleNotify: 受信したすべてのノーティファイ メッセージを処理します。入力データのサブスクリプションは、アプリケーションに代わって FogFlow によって行われることに注意してください。
	
.. code-block:: javascript

    const NGSIClient = require('./ngsi/ngsiclient.js');
    const NGSIAgent = require('./ngsi/ngsiagent.js');
    
    // get the listening port number from the environment variables given by the FogFlow edge worker
    var myport = process.env.myport;
    
    // set up the NGSI agent to listen on 
    NGSIAgent.setNotifyHandler(handleNotify);
    NGSIAgent.setAdminHandler(handleAdmin);
    NGSIAgent.start(myport, startApp);
    
    process.on('SIGINT', function() {
        NGSIAgent.stop();	
        stopApp();
        
        process.exit(0);
    });
    
    
    function startApp() 
    {
        console.log('your application starts listening on port ',  myport);
        
        // to add the initialization part of your application logic here    
    }
    
    function stopApp() 
    {
        console.log('clean up the app');
    }
	
	
	
* アプリケーションの起動時に何かを初期化する方法は？

.. code-block:: javascript

	function startApp() 
	{
	    console.log('your application starts listening on port ',  myport);
    
	    // to add the initialization part of your application logic here    
	}

* アプリケーションは FogFlow によってどのように構成されていますか？

        生成結果 (any generate results) の場合、アプリケーションは、特定の IoT Broker に NGSI アップデートを送信することで、それらを公開できます。特定の IoT Broker と、生成結果の公開に使用されるエンティティ タイプは、FogFlow によって構成されます。

        FogFlow がオペレーターに関連付けられたタスクを起動すると、FogFlow woker はリスニング ポートを介して構成をタスクに送信します。次のコード ブロックでは、オペレーターがこの提供された構成をどのように処理する必要があるかを示します。次のグローバル変数は FogFlow によって構成されます。
	
	* brokerURL: FogFlow によって割り当てられた IoT Broker の URL。
	* ngsi10client: 指定された brokerURL に基づいて作成された NGSI10 クライアント
	* myReferenceURL: アプリケーションがリッスンしている IP アドレスとポート番号。この参照 URL を使用して、アプリケーションは提供された IoT Broker からの追加の入力データをサブスクライブできます。
	* outputs: 出力エンティティの配列。
	* isConfigured: 構成が完了したかどうかを示します。

.. code-block:: javascript

    // global variables to be configured 
    var ngsi10client;
    var brokerURL;
    var myReferenceURL;
    var outputs = [];
    var isConfigured = false;
    
    // handle the configuration commands issued by FogFlow
    function handleAdmin(req, commands, res) 
    {	
        handleCmds(commands);
        
        isConfigured = true;
        
        res.status(200).json({});
    }
    
    // handle all configuration commands
    function handleCmds(commands) 
    {
        for(var i = 0; i < commands.length; i++) {
            var cmd = commands[i];
            handleCmd(cmd);
        }	
    }
    
    // handle each configuration command accordingly
    function handleCmd(commandObj) 
    {	
        switch(commandObj.command) {
            case 'CONNECT_BROKER':
                connectBroker(commandObj);
                break;
            case 'SET_OUTPUTS':
                setOutputs(commandObj);
                break;
            case 'SET_YOUR_REFERENCE':
                setReferenceURL(commandObj);
                break;            
        }	
    }
    
    // connect to the IoT Broker
    function connectBroker(cmd) 
    {
        brokerURL = cmd.brokerURL;
        ngsi10client = new NGSIClient.NGSI10Client(brokerURL);
        console.log('connected to broker', cmd.brokerURL);
    }
    
    function setOutputs(cmd) 
    {
        var outputStream = {};
        outputStream.id = cmd.id;
        outputStream.type = cmd.type;
    
        outputs.push(outputStream);
    
        console.log('output has been set: ', cmd);
    }
    
    function setReferenceURL(cmd) 
    {
        myReferenceURL = cmd.referenceURL   
        console.log('your application can subscribe additional inputs under the reference URL: ', myReferenceURL);
    }
            

* 受信したエンティティ データを処理する方法は？

.. code-block:: javascript

	// handle all received NGSI notify messages
	function handleNotify(req, ctxObjects, res) 
	{	
		console.log('handle notify');
		for(var i = 0; i < ctxObjects.length; i++) {
			console.log(ctxObjects[i]);
	        fogfunction.handler(ctxObjects[i], publish);
		}
	}

	// process the input data stream accordingly and generate output stream
	function processInputStreamData(data) 
	{
		var type = data.entityId.type;
		console.log('type ', type);
		
		// do the internal data processing
		if (type == 'PowerPanel'){
			// to handle this type of input
		} else if (type == 'Rule') {
			// to handle this type of input
		}	
	}

* アプリケーション内でアップデートを送信する方法は？

        生成結果( any generate results) の場合、アプリケーションは、特定の IoT Broker に NGSI アップデートを送信することで、それらを公開できます。特定の IoT Broker と、生成結果の公開に使用されるエンティティ タイプは、FogFlow によって構成されます。

.. code-block:: javascript

    // update context for streams
    function updateContext(anomaly) 
    {
        if (isConfigured == false) {
            console.log('the task is not configured yet!!!');
            return;
        }
            
        var ctxObj = {};
        
        ctxObj.entityId = {};
        
        var outputStream = outputs[0];
        
        ctxObj.entityId.id = outputStream.id;
        ctxObj.entityId.type = outputStream.type;
        ctxObj.entityId.isPattern = false;
        
        ctxObj.attributes = {};
        
        ctxObj.attributes.when = {		
            type: 'string',
            value: anomaly['when']
        };
        ctxObj.attributes.whichpanel = {
            type: 'string',
            value: anomaly['whichpanel']
        };  
            
        ctxObj.attributes.shop = {
            type: 'string',
            value: anomaly['whichshop']
        };  
        ctxObj.attributes.where = {
            type: 'object',
            value: anomaly['where']
        };  
        ctxObj.attributes.usage = {
            type: 'integer',
            value: anomaly['usage']
        };
            
        ctxObj.metadata = {};		
        ctxObj.metadata.shop = {
            type: 'string',
            value: anomaly['whichshop']
        };  
                
        ngsi10client.updateContext(ctxObj).then( function(data) {
            console.log('======send update======');
            console.log(ctxObj);
            console.log(data);
        }).catch(function(error) {
            console.log(error);
            console.log('failed to update context');
        });    
    }
        

* アプリケーションにサードパーティのライブラリが必要な場合はどうなりますか？

        アプリケーションでサードパーティのライブラリが必要な場合は、"package.json" でそれらとそのバージョン番号を指定してください。
	

* アプリケーションが必要な場合、追加の入力をサブスクライブする方法は？



オペレーターを Docker 化する方法
--------------------------------

	.. code-block:: bash
		
		# package.json を変更して、必要なすべてのサードパーティ ライブラリをフェッチしてください。
		npm install 
		
		# Docker イメージをビルドします。
		./build 
		
		# FogFlow Docker レジストリーに Docker イメージをプッシュします。
		docker push task1
		

オペレーターをデバッグおよびテストする方法
------------------------------------------

サービス トポロジーで使用する前に、オペレータを手動でテストすることもできます。現在、"subscription.json" で提供されている例のようにサブスクリプションを定義してから、既知の Cloud IoT Broker にサブスクライブ リクエストを発行できます。

手順は次のとおりです:

#. FogFlow からアプリケーションを実行します。

	.. code-block:: bash
		
		# package.json を変更して、必要なすべてのサードパーティ ライブラリをフェッチしてください。
		npm install 
	
		# 空きポートを見つけて、この空きポートを使用する必要があります (例)
		# 環境変数を設定する数値
		export myport=100010
		
		# Docker イメージをビルドします。
		node main.js 

#. 実行中のアプリケーションに入力データを提供するためにサブスクリプションを発行します。

	The bash script in "test.sh" shows how you can 

	.. code-block:: bash
		
		# package.json を変更して、必要なすべてのサードパーティ ライブラリをフェッチしてください。
		./test.sh 

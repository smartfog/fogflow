*****************************
システム概要
*****************************

FogFlow は、クラウドとエッジ上の動的処理フローをサポートする分散実行フレームワークです。複数の NGSI ベースのデータ処理タスクを動的かつ自動的に合成して高レベルの IoT サービスを形成し、共有クラウド エッジ環境内でそれらのサービスの展開を調整および最適化できます。

FogFlow クラウド エッジの共有環境は、次の図に示すように、1つの FogFlow クラウド ノードと複数の FogFlow エッジ ノードで作成できます。この図では、FogFlow システムで実行されているすべての統合機能を確認できます。

.. figure:: ../../en/source/figures/FogFlow_System_Design.png

FogFlow は現在、後で使用するために内部エンティティ データを保存するためのグラフデータベースをサポートしています。FogFlow システムでは、グラフデータベースの相互作用ポイントはDesigner ーです。したがって、FogFlow Web ブラウザーまたは curl や python クライアントなどの他のクライアントを介して作成する FogFlow エンティティは、エンティティ作成要求が直接 Designer に送信されます。次に、 Designer はこのエンティティ リクエストを Cloud Broker に送信してエンティティを登録し、同時に Designer はこのデータを Dgraph に送信してデータベースに保存します。

FogFlow システムが再起動するたびに、Designer は Dgraphへ のクエリ リクエストをトリガーし、保存されているすべてのエンティティを Dgraph から取得してから、これらのエンティティを Cloud Broker に送信して登録します。

.. note:: FogFlow Designer は、Discovery と CloudBroker に依存しています。

FogFlow Designer と Dgraph の統合は、gRPC を介して行われます。ここで、Dgraph クライアントは、gRPC を使用して以下のコードのように Designer サーバーで実装されます。gRPC のデフォルトポートは9080で、FogFlow では9082ポートが使用されます。

.. code-block:: console

   /*
   creating grpc client for making connection with dgraph
   */
   function newClientStub() {
       return new dgraph.DgraphClientStub(config.HostIp+":"+config.grpcPort, grpc.credentials.createInsecure());
   }

   // Create a client.
   function newClient(clientStub) {
      return new dgraph.DgraphClient(clientStub);
   }
   
新しい FogFlow エンティティが Web ブラウザーによって作成されるたびに、要求は Designer に送信されるか、ユーザーは任意のクライアントを使用して Designer 上にエンティティを直接作成します。次に、Designer は2つのタスクを実行します:   

1. エンティティを登録するために Cloud Broker にリクエストを送信します。

2. Dgraph クライアントを呼び出してエンティティ データを保存します。Dgraph クライアントは、スキーマを作成した後、Dgraph サーバーとの接続を作成し、データを Dgraph に送信します。これとは別に、FogFlow システムが再起動すると、Designer から別のフローがトリガーされます。このフローでは、Designer は Dgraph からすべての保存されたエンティティ データをクエリし、これらのエンティティを登録するために Cloud Broker に転送します。

スキーマを作成し、Dgraph にデータを挿入するためのコードを垣間見ることができます。

.. code-block:: console

   /*
   create schema for node
   */
   async function setSchema(dgraphClient) {
       const schema = `
            attributes: [uid] .
            domainMetadata: [uid] .
            entityId: uid .
            updateAction: string .
            id: string .
            isPattern: bool .
            latitude: float .
            longitude: float .
            name: string .
            type: string .
	          value: string . 
       `;
       const op = new dgraph.Operation();
       op.setSchema(schema);
       await dgraphClient.alter(op);
   }
   
   /*
   insert data into database
   */
   async function createData(dgraphClient,ctx) {
       const txn = dgraphClient.newTxn();
       try {
           const mu = new dgraph.Mutation();
           mu.setSetJson(ctx);
           const response = await txn.mutate(mu);
           await txn.commit();
       }
	    finally {
          await txn.discard();
       }
   }
   
詳細なコードについては、 https://github.com/smartfog/fogflow/blob/development/designer/dgraph.js を参照してください。

このページでは、FogFlow 統合について簡単に紹介します。詳細については、リンクを参照してください。

統合には主にノースバウンドとサウスバウンドの2種類があり、センサー デバイスからブローカーへのデータ フローはノースバウンド フローと呼ばれ、ブローカーからアクチュエータ デバイスへのデータ フローはサウスバウンド フローと呼ばれます。ノースバウンドおよびサウスバウンドのデータ フローの詳細は、この (`this`_) ページで確認できます。

.. _`this`: https://fogflow.readthedocs.io/en/latest/integration.html

FogFlow は Scorpio Broker と統合できます。Scorpio は、NGSI-LD 準拠のコンテキスト ブローカーです。そのため、NGSI-LD アダプターは、FogFlow エコシステムが Scorpio Context Broker と対話できるように構築されています。NGSI-LD アダプターは、NGSI データ形式を NGSI-LD に変換し、それを Scorpio Broker に転送します。詳細は、 `Integrate FogFlow with Scorpio Broker`_ のページで確認できます。

.. _`Integrate FogFlow with Scorpio Broker`: https://fogflow.readthedocs.io/en/latest/scorpioIntegration.html


**FogFlow NGSI-LD サポート:** FogFlow は、NGSI9 および NGSI10 とともに NGSI-LD APIs サポートを提供しています。NGSI-LD 形式は、FogFlow または他の GEs のコンポーネント間のコンテキスト共有通信を明確かつよりよく理解する目的で  **FIWARE が使用するリンクトデータモデル** を利用することを目的としています。より効率的な方法で情報を推測するためにデータ間の関係を確立することにより、NGSIv1 モデルおよび NGSIv2 モデルの間でデータを維持する複雑さを軽減します。

- このモデルを組み込む理由は、エッジ コンピューティングのバックボーンを形成しているリンクトデータの関連付けが直接必要なためです。
- これにより、FogFlow と Scorpio Broker 間の相互作用のように、相互作用が可能になったため、FogFlow と他の GEs の間のギャップが埋められます。

NGSI-LD APIs の詳細は、 `API Walkthrough`_  ページで確認できます。

.. _`API Walkthrough`: https://fogflow.readthedocs.io/en/latest/api.html#ngsi-ld-supported-apis


FogFlow は、NGSI APIs を使用して Orion Context Broker と統合することもできます。詳細は、FogFlow とFIWARE の統合ページ (`Integrate FogFlow with FIWARE`_) で確認できます。

.. _`Integrate FogFlow with FIWARE`: https://fogflow.readthedocs.io/en/latest/fogflow_fiware_integration.html


同様に、WireCloud との FogFlow 統合は、WireCloud のさまざまなウィジェットを使用してデータを視覚化するために提供されています。QuantumLeap との FogFlow 統合は、時系列ベースの履歴データを保存することです。詳細については 、WireCloud の場合は FogFlow を WireCloud と統合 (`Integrate FogFlow with WireCloud`_) のページ、QuantumLeap の場合は FogFlow を QuantumLeap と統合 (`Integrate FogFlow with QuantumLeap`_) のページで確認できます。

.. _`Integrate FogFlow with WireCloud`: https://fogflow.readthedocs.io/en/latest/wirecloudIntegration.html
.. _`Integrate FogFlow with QuantumLeap`: https://fogflow.readthedocs.io/en/latest/quantumleapIntegration.html


FogFlow は、FogFlow クラウド ノードと FogFlow エッジ ノード間、および2つのエッジ ノード間の安全な通信も提供します。FogFlow で HTTP ベースの安全な通信を実現するには、FogFlow クラウド ノードと FogFlow エッジ ノードが独自のドメイン名を持っている必要があります。さらに、詳細な構成とセットアップの手順は、セキュリティ (`Security`_) のページで確認できます。

.. _`Security`: https://fogflow.readthedocs.io/en/latest/https.html

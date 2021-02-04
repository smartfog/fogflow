*******************************
Kubernetes のインテグレーション
*******************************

FogFlow のコンポーネントは、ソースコードを介して構築することも、docker-compose ツールを使用して docker 環境で構築することもできます。Docker 環境では、FogFlow の各コンポーネントが単一インスタンスとして実行されています。単一のコンポーネント コンテナがダウンした場合、FogFlow システム全体を再起動する必要があり、単一のサービスが過負荷になると、負荷を処理するように拡張できなくなります。これらの問題を克服するために、FogFlow は Kubernetes に移行しました。FogFlow コンポーネントは、エンドユーザーの要件に基づいて Kubernetes クラスター環境にデプロイされます。さまざまなクラスター構成を展開できます:

1.      同じノード上のマスターとワーカー
2.      シングル マスターとシングル ワーカー
3.      シングル マスター ノードとマルチ ワーカー ノード
4.      複数のマスター ノードとマルチ ワーカー ノード


クラスターに加えて、K8s の次の機能が FogFlow に実装されています:

1. **高可用性と負荷分散**:
高可用性とは、単一障害点が発生しないように、Kubernetes とそのサポート コンポーネントをセットアップすることです。環境セットアップで複数のアプリケーションが単一のコンテナーで実行されている場合、そのコンテナーは簡単に失敗する可能性があります。Kubernetes の高可用性のための仮想マシンと同じように、コンテナーの複数のレプリカを実行できます。負荷分散は、バックエンド サーバーのグループ全体に着信ネットワーク トラフィックを分散するのに効率的です。ロードバランサーは、サーバーのクラスター全体にネットワークまたはアプリケーション トラフィックを分散するデバイスです。ロードバランサーは、クラスターの高可用性とパフォーマンスの向上を実現するために大きな役割を果たします。
 
2. **自己回復 (Self-healing)**: ポッドのいずれかが手動で削除された場合、またはポッドが誤って削除されたか、再起動された場合。Kubernetes にはポッドを自動修復する機能があるため、デプロイによってポッドが確実に戻されます。

3. **自動ロールアウトとロールバック**: ローリング アップデートによって実現できます。ローリング アップデートは、実行中のアプリのバージョンをアップデートするためのデフォルトの戦略です。以前のポッドのサイクルを更新し、新しいポッドを段階的に取り込みます。導入された変更によって本番環境が中断された場合は、その変更をロールバックする計画が必要です。Kubernetes と kubectl は、デプロイなどのリソースへの変更をロールバックするためのシンプルなメカニズムを提供します。

4. **Helmサポートでデプロイを容易にする**: Helm は、Kubernetes アプリケーションのインストールと管理を合理化するツールです。Kubernetes アプリケーションの管理に役立ちます。Helm Chart は、最も複雑な Kubernetes アプリケーションでさえも定義、インストール、アップグレードするのに役立ちます。FogFlow ドキュメントは、Kubernetes 環境を十分に理解してアクセスできるように、上記の機能の機能の詳細で更新されます。


**FogFlow K8s インテグレーションの制限**

以下は、FogFlow Kubernetes ンテグレーションのいくつかの制限です。これらの制限は、将来 FogFlow で実装される予定です。

1. FogFlow ワーカーが起動するタスク インスタンスはポッドに実装されていません。K8s ポッドを介した起動タスク インスタンスの移行は、FogFlow OSS の将来の範囲に含まれます。

2. FogFlow Edge ノード K8s のサポート

3. K8s 環境のセキュリティとネットワークポリシー

4. Taints と Trait

5. パフォーマンス評価

6. その他の機能


Kubernetes の FogFlow クラウド アーキテクチャ図
----------------------------------------------




.. figure:: figures/k8s-architecture.png




Dgraph、Discovery、Broker、Designer、Master、Worker、Rabbitmq などの FogFlow クラウド ノード コンポーネントは、クラスター ノードに分散されています。FogFlow コンポーネント間の通信とその動作は以前と同じであり、ワーカー ノードは Docker コンテナーでタスクインスタンスを起動します。


こちら `here`_ のリンクをたどって、Kubernetesコンポーネントがどのように機能するかを確認してください。

.. _`here`: https://kubernetes.io/docs/concepts/overview/components/



K8s で FogFlow を実行するための前提条件コマンドは次のとおりです。

1. docker
2. Kubernetes
3. Helm

.. important:: 
	**ユーザーが sudo なしで Docker コマンドを実行できるようにしてください。**
	
Kubernetes をインストールするには、`Kubernetes Official Site`_ を参照する か、代替の `Install Kubernetes`_ を確認してください。

Helm をインストールするには、Install Helm`_ を参照してください。

.. _`Kubernetes Official Site`: https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

.. _`Install Kubernetes`: https://medium.com/@vishal.sharma./installing-configuring-kubernetes-cluster-on-ubuntu-18-04-lts-hosts-f37b959c8410

.. _`Install Helm`: https://helm.sh/docs/intro/install/


K8s 環境に FogFlow クラウド コンポーネントをデプロイ
--------------------------------------------------

Dgraph、Discovery、Broker、Designer、Master、Worker、Rabbitmq などの FogFlow クラウド ノード コンポーネントは、クラスター ノードに分散されています。FogFlow コンポーネント間の通信とその動作は通常どおりであり、ワーカー ノードは Docker コンテナーでタスク インスタンスを起動します。


**必要なすべてのスクリプトを取得します**

以下のように、Kubernetes ファイルと構成ファイルをダウンロードします。

.. code-block:: console    

	# the Kubernetes yaml file to start all FogFlow components on the cloud node
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/helm/fogflow-chart
	
	# the configuration file used by all FogFlow components
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/yaml/config.json

	# the configuration file used by the nginx proxy
	wget https://raw.githubusercontent.com/smartfog/fogflow/master/yaml/nginx.conf

	
   
環境に応じて IP 構成を変更します
-------------------------------------------------------------

ご使用の環境に応じて、config.json で以下の IP アドレスを変更する必要があります。

- **coreservice_ip**: すべての FogFlow エッジ ノードが FogFlow クラウド ノードのコア サービス (ポート80 の nginx やポート5672のrabbitmq など) にアクセスするために使用します。通常、これは FogFlow クラウドノードのパブリック IP になります。
- **external_hostip**: FogFlow クラウド ノードの構成の場合、これは、実行中の FogFlow コア サービスにアクセスするためにコンポーネント (CloudWorker および CloudBroker) によって使用される coreservice_ip と同じです。
- **internal_hostip**: これはデフォルトの K8s ネットワーク インターフェースの IP であり、Linux ホストの "cni0" ネットワークインターフェースです。
- **site_id**: 各 FogFlow ノード (クラウド ノードまたはエッジ ノード) は、システム内でそれ自体を識別するために、一意の文字列ベースの ID を持っている必要があります。
- **physical_location**: FogFlow ノードの地理的位置 (geo-location)
- **worker.capacity**: FogFlow ノードが呼び出すことができる Docker コンテナーの最大数を意味します。


values.yaml ファイルを変更
---------------------------

- 要件に従って名前空間 (namespace) を編集します。必要なreplicaCount 数を追加します。

- 環境 hostPath に従って、values.yaml ファイルの dgraph、configJson、および nginxConf パスを変更します。

- 環境に応じて externalIP を変更します。

.. code-block:: console

      #Kubernetes namespace of FogFlow components
      namespace: default

      #replicas will make sure that no. of replicaCount mention in values.yaml
      #are running all the time for the deployment
      replicaCount: 1

      serviceAccount:
      #Specifies whether a service account should be created
        create: true
      #Annotations to add to the service account
        annotations: {}
      #The name of the service account to use.
      #If not set and create is true, a name is generated using the fullname template
        name: ""

      #hostPath for dgraph volume mount
      dgraph:
        hostPath:
          path: /mnt/dgraph

      #hostPath for config.json
      configJson:
        hostPath:
          path: /home/necuser/fogflow/helm/files/fogflow-chart/config.json

      #hostPath for nginx.conf
      nginxConf:
        hostPath:
          path: /home/necuser/fogflow/fogflow/yaml/nginx.conf

      #External IP to expose cluster
      Service:
       spec:
        externalIPs:
        - XXX.XX.48.24

	  
Helm Chart を使用してすべての Fogflow コンポーネントを開始
-------------------------------------------------------------

Helm-Chart フォルダーの外部から Helm コマンドを実行して、FogFlow コンポーネントを起動します。ここでは helm-chart 名は "fogflow-chart" です。

コマンドラインから設定を渡すには、helm install コマンドで "--set" フラグを追加します。

.. code-block:: console
 
          helm install ./fogflow-chart --set externalIPs={XXX.XX.48.24} --generate-name


詳細については、Helmの公式 `link_` を参照してください。

.. _`link`: https://helm.sh/docs/helm/

セットアップを検証
-------------------------------------------------------------

FogFlow クラウド ノードが正しく開始されているかどうかを確認するには、次の2つの方法があります:

- "kubectl get pods --namespace = <namespace_name>" を使用して、すべてのポッドが稼働していることを確認します

.. code-block:: console  

         kubectl get pods --namespace=fogflow
		 
		 
        NAME                           READY   STATUS              RESTARTS   AGE
        cloud-broker-c78679dd8-gx5ds   1/1     Running             0          8s
        cloud-worker-db94ff4f7-hwx72   1/1     Running             0          8s
        designer-bf959f7b7-csjn5       1/1     Running             0          8s
        dgraph-869f65597c-jrlqm        1/1     Running             0          8s
        discovery-7566b87d8d-hhknd     1/1     Running             0          8s
        master-86976888d5-drfz2        1/1     Running             0          8s
        nginx-69ff8d45f-xmhmt          1/1     Running             0          8s
        rabbitmq-85bf5f7d77-c74cd      1/1     Running             0          8s

		
- FogFlow DashBoard からシステムステータスを確認します

システムステータスは、Web ブラウザの FogFlow ダッシュボードから確認して、次の URL で現在のシステム ステータスを確認することもできます: http://<coreservice_ip>/index.html


ソースコードからすべてをビルド
=========================================

FogFlow は、ARM プロセッサと x86 プロセッサ (32ビットと64ビット) の両方の Linux でビルドおよびインストールできます。

依存関係をインストール
----------------------

#. FogFlowをビルドするには、最初に次の依存関係をインストールします。

	- git クライアントのインストール: https://www.digitalocean.com/community/tutorials/how-to-install-git-on-ubuntu-16-04 の指示に従ってください。
	
	- Docker CE のインストール: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-16-04 の指示に従ってください。
	
		.. note:: すべてのスクリプトは、Docker を sudo なしで実行できることを前提に作成されています。
	

	- 最新バージョンの golang(>v1.9) をインストールします: https://golang.org/doc/install の情報に従ってダウンロードしてインストールしてください。

	- Node.js (>6.11) および npm (>3.10) をインストールします: https://nodejs.org/en/download/ の情報に従ってダウンロードしてインストールしてください。


#. インストールされているバージョンを確認します。


	.. code-block:: bash

		go version   #output  go version go1.9 linux/amd64 
  		nodejs -v    #output 	v6.10.2
  		npm -v       #output  3.10.10


#. 環境変数 GOPATH を設定します。


	.. note:: GOPATH は、go ベースのプロジェクトのワークスペースを定義します。go ワークスペース フォルダーには "src" フォルダーが必要であり、FogFlow コード リポジトリはこの "src" フォルダーに複製される必要があることに注意してください。たとえば、ホームフォルダが "/home/smartfog" であると仮定して、ワークスペースとして新しいフォルダ "go" を作成します。この場合、最初に "/home/smartfog/go" の下に "src" (正確にこの名前である必要があります) を作成してから、"/home/smartfog/go/src フォルダー内の FogFlow コード リポジトリをチェック アウトする必要があります。

	.. code-block:: bash	

		export GOPATH="/home/smartfog/go"


#. コード リポジトリをチェック アウトします。

	.. code-block:: bash	
		
		cd /home/smartfog/go/src/	
		git clone https://github.com/smartfog/fogflow.git
		
		
#. 以下のようにソースコードからすべてのコンポーネントをビルドします。


IoT Discovery をビルド
------------------------

	- ネイティブ実行可能プログラムをビルドします。
	
		.. code-block:: bash	
			
			# go the discovery folder
			cd /home/smartfog/go/src/fogflow/discovery
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- Docker イメージを作成します。

		.. code-block:: bash			
		
			# Simply ./build  can be run to perform the following commands
		
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/discovery" .										
		
			
IoT Broker をビルド
--------------------------

	- ネイティブ実行可能プログラムをビルドします。
	
		.. code-block:: bash	
			
			# go the broker folder
			cd /home/smartfog/go/src/fogflow/broker
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- Docker イメージを作成します。
		
		.. code-block:: bash			
		
			# simply ./build can be run to perform the following commands		
				
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/broker" .			



Topology Master をビルド
--------------------------

	- ネイティブ実行可能プログラムをビルドします。
	
		.. code-block:: bash	
			
			# go the master folder
			cd /home/smartfog/go/src/fogflow/master
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- Docker イメージを作成します。
		
		.. code-block:: bash							
		
			# simply ./build can be run to perform the following commands		
					
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/master" .			



Worker をビルド
--------------------------

	- ネイティブ実行可能プログラムをビルドします。
	
		.. code-block:: bash	
			
			# go the worker folder
			cd /home/smartfog/go/src/fogflow/worker
			# download its third-party library dependencies
			go get
			# build the source code
			go build
	
	- Docker イメージを作成します。
		
		.. code-block:: bash	
					
			# simply ./build  can be run to perform the following commands									
			
			# download its third-party library dependencies
			go get
			# build the source code and link all libraries statically
			CGO_ENABLED=0 go build -a
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/worker" .			


Task Designer をビルド
--------------------------

	- サードパーティのライブラリの依存関係をインストールします。
	
		.. code-block:: bash	
			
			# go the designer folder
			cd /home/smartfog/go/src/fogflow/designer
			
			# install all required libraries
			npm install
	
	- Docker イメージを作成します。
		
		.. code-block:: bash	
		
			# simply ./build can be run to perform the following commands					

			# install all required libraries
			npm install
			
			# create the docker image; sudo might have to be used to run this command 
			# if the docker user is not in the sudo group
			docker build -t "fogflow/designer"  .

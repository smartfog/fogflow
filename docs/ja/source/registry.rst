イメージ管理用の Docker レジストリーを設定
==============================================


自己署名証明書 (Self Signed Certificate)を作成
----------------------------------------------
プライベート Docker レジストリーに使用するには、サーバー上に自己署名証明書を作成する必要があります。

.. code-block:: console    

	mkdir registry_certs
	openssl req -newkey rsa:4096 -nodes -sha256 \
                -keyout registry_certs/domain.key -x509 -days 356 \
                -out registry_certs/domain.cert
	ls registry_certs/

最後に、2つのファイルがあります:

- domain.cert – このファイルは、プライベート レジストリーを使用してクライアントに処理できます。
- domain.key – これは TLS でプライベート レジストリーを実行するために必要な秘密キーです。


TLS を使用してプライベート Docker レジストリーを実行
--------------------------------------------------
これで、ローカル ドメイン証明書とキーファイルを使用してレジストリーを開始できます:

.. code-block:: console    

	docker run -d -p 5000:5000 \
		 	-v $(pwd)/registry_certs:/certs \
 		 	-e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.cert \
 		 	-e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
 			--restart=always --name registry registry:2

ここでは、フォルダー /registry_certs をボリュームとして Docker レジストリー コンテナーにマップします。証明書とキーファイルを指す環境変数を使用します。

これで、ローカル イメージを新しいレジストリーにプッシュできます:

.. code-block:: console    

	docker push localhost:5000/proxy:1.0.0


フォグ ノード (Fog node) からリモート レジストリーにアクセス
---------------------------------------------------------

これで、プライベート レジストリーが TLS サポートで開始されるため、ドメイン証明書を持つ任意のクライアントからレジストリーにアクセスできます。
そのため、証明書ファイル "domain.cert" はクライアントのファイル内に配置する必要があります。

.. code-block:: console    

	/etc/docker/certs.d/<registry_address>/ca.cert

ここで、<registry_address> はサーバーのホスト名です。証明書が更新されたら、ローカルの Docker デーモンを再起動する必要があります:


.. code-block:: console    

	mkdir -p /etc/docker/certs.d/dock01:5000 
	cp domain.cert /etc/docker/certs.d/dock01:5000/ca.crt
	service docker restart
	
	
これで、画像を新しいプライベート レジストリーにプッシュできます:

.. code-block:: bash

	docker tag imixs/proxy dock01:5000/proxy:dock01
	docker push dock01:5000/proxy:dock01
	

Docker レジストリー フロントエンドを起動します
------------------------------------------------

プロジェクト konradkleine/docker-registry-frontend は、Web ブラウザーを介したレジストリーへのアクセスを簡素化するために使用できるクールな Web フロントエンドを提供します。
docker-registry-frontend は、Docker コンテナーとして開始できます。レジストリーがで実行されていると仮定します。

https://yourserver.com:5000

次の docker run コマンドを使用して、フロントエンド コンテナーを起動します:

.. code-block:: console    

	docker run \
 		-d \
 		-e ENV_DOCKER_REGISTRY_HOST=yourserver.com \
 		-e ENV_DOCKER_REGISTRY_PORT=5000 \
 		-e ENV_DOCKER_REGISTRY_USE_SSL=1 \
 		-p 0.0.0.0:80:80 \
 		konradkleine/docker-registry-frontend:v2
		
これで、Web ブラウザの URL を介してレジストリーにアクセスできます:

http://localhost:80/

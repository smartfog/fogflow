# FogFlow

![CI/CD Status](https://github.com/smartfog/fogflow/workflows/CI/CD%20Status/badge.svg?branch=development)
[![FIWARE Security](https://nexus.lab.fiware.org/static/badges/chapters/processing.svg)](https://www.fiware.org/developers/catalogue/)
[![License: BSD-4-Clause](https://img.shields.io/badge/license-BSD%204%20Clause-blue.svg)](https://spdx.org/licenses/BSD-4-Clause.html)
[![Docker Status](https://img.shields.io/docker/pulls/fogflow/discovery.svg)](https://hub.docker.com/r/fogflow)
[![Support badge](https://img.shields.io/badge/tag-fiware--fogflow-orange.svg?logo=stackoverflow)](https://stackoverflow.com/search?q=%5Bfiware%5D%20fogflow)
<br>
[![Documentation badge](https://img.shields.io/readthedocs/fogflow.svg)](http://fogflow.readthedocs.org/en/latest/)
![Status](https://nexus.lab.fiware.org/repository/raw/public/static/badges/statuses/fogflow.svg)
[![Build Status](https://travis-ci.org/smartfog/fogflow.svg?branch=master)](https://travis-ci.org/smartfog/fogflow)
[![Swagger Validator](https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/OAI/OpenAPI-Specification/master/examples/v2.0/json/petstore-expanded.json.svg)](https://app.swaggerhub.com/apis/fogflow/broker/1.0.0)

FogFlow は、コンテキストによって駆動されるクラウドおよびエッジ上で動的データ処理
フローを自動的に調整する IoT エッジコンピューティング フレームワークです。
すべてのレイヤーから利用可能なシステムリソースのシステム コンテキスト、利用可能な
すべてのデータ エンティティの登録済みメタデータのデータ コンテキスト、および
ユーザーによって定義された予想される QoS の使用コンテキストを含みます。

このプロジェクトは [FIWARE](https://www.fiware.org/) の一部です。詳細については、
FIWARE Catalogue エントリの
[Processing](https://github.com/Fiware/catalogue/tree/master/processing)
を確認してください 。

| :books: [ドキュメント](https://fogflow.readthedocs.io/ja/latest/) | :mortar_board: [Academy](https://fiware-academy.readthedocs.io/en/latest/processing/fogflow) |:whale: [Docker Hub](https://hub.docker.com/r/fogflow) | :dart: [ロードマップ](https://github.com/smartfog/fogflow/blob/master/doc/roadmap.ja.md) |
| --- | --- | --- | --- |

## コンテンツ

-   [バックグラウンド](#background)
-   [インストール](#installation)
-   [チュートリアル](https://fogflow.readthedocs.io)
-   [API](#api)
-   [テスト](#testing)
-   [品質保証](#quality-assurance)
-   [ロードマップ](./doc/roadmap.ja.md)
-   [詳細情報](#more-information)
-   [ライセンス](#license)

<a name="background"/>

## バックグラウンド

FogFlow は、サービス プロバイダーがクラウドとエッジを介して IoT サービスを簡単に
プログラムおよび管理するための標準ベースのデータ処理フレームワークです。
以下は、FogFlow のモチベーション、機能、および利点です。

-   _なぜ FogFlow が必要なのですか？_

    -   クラウドのみのソリューションのコストは高すぎて、1,000を超える地理分散
        デバイスで大規模な IoT システムを実行できません。
    -   多くの IoT サービスでは、エンドツーエンドのレイテンシが10ミリ秒未満などの
        高速な応答時間が必要です。
    -   サービス プロバイダーは、クラウド エッジ環境で IoT サービスを迅速に設計
        および展開するために、非常に複雑でコストに直面しています。
    -   ビジネスの需要は時間とともに急速に変化しており、サービスプロバイダーは、
        共有クラウドエッジインフラストラクチャ上で新しいサービスを高速で試して
        リリースする必要があります。
    -   地理的に分散した ICT インフラストラクチャ上で IoT サービスを迅速に設計
        および展開するためのプログラミングモデルの欠如
    -   さまざまなアプリケーション間でデータと派生した結果を共有および再利用する
        ための相互運用性とオープン性の欠如

-   _FogFlow は何を提供しますか？_

    -   効率的なプログラミングモデル: サービスのプログラミングは、レゴ ブロックを
        構築するようなものです。
    -   動的サービスオーケストレーション: 必要な場合にのみ必要なデータ処理を
        開始します。
    -   最適化されたタスクの展開: プロデューサーとコンシューマーのローカリティに
        基づいてクラウドとエッジ間でタスクを割り当てます。
    -   スケーラブルなコンテキスト管理: プロデューサーとコンシューマーの間で
        柔軟な情報交換 (トピックベースとスコープベースの両方) を可能にします。

-   _利用者は FogFlow からどのように利益を得ることができますか？_

    -   地理的に分散した共有 ICT インフラストラクチャを介して新しいサービスを
        実現およびリリースする際の市場投入までの時間の短縮。
    -   さまざまなサービスを運用する際の運用コストと管理の複雑さの軽減。
    -   低遅延と高速応答時間を必要とするサービスを提供可能。

<a name="installation"/>

## インストール

FogFlow のインストール手順は、
[ インストールガイド](https://fogflow.readthedocs.io/en/latest/setup.html)
に記載されています。

<a name="api"/>

## API

API とその使用例は、
[こちら](https://fogflow.readthedocs.io/en/latest/api.html)
にあります。

<a name="testing"/>

## テスト

基本的なエンド ツー エンドのテストを実行するには、
[こちら](https://fogflow.readthedocs.io/en/latest/test.html)
の詳細な手順に従ってください。

<a name="quality-assurance"/>

## 品質保証

このプロジェクトは [FIWARE](https://fiware.org/) の一部であり、次のように評価されています:

-   **テストされたバージョン:**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Version&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.version&colorB=blue)
-   **ドキュメンテーション:**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Completeness&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.docCompleteness&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Usability&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.docSoundness&colorB=blue)
-   **即応性 (Responsiveness):**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Time%20to%20Respond&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.timeToCharge&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Time%20to%20Fix&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.timeToFix&colorB=blue)
-   **FIWARE テスト (FIWARE Testing):**
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Tests%20Passed&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.failureRate&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Scalability&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.scalability&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Performance&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.performance&colorB=blue)
    ![ ](https://img.shields.io/badge/dynamic/json.svg?label=Stability&url=https://fiware.github.io/catalogue/json/fogflow.json&query=$.stability&colorB=blue)

<a name="more-information"/>

## 詳細情報

-   [チュートリアル](http://fogflow.readthedocs.io/en/latest/index.html)
-   [IoT-J ペーパー](http://ieeexplore.ieee.org/document/8022859/)

<a name="license"/>

## ライセンス

FogFlow は、
[BSD-4-Clause](https://spdx.org/licenses/BSD-4-Clause.html)
の下でライセンスされてい ます。

© 2017-2020 NEC Laboratories Europe GmbH

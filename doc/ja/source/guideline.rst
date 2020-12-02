************************************************
コントリビューションとコーディングのガイドライン
************************************************

FogFlow コントリビューション ガイド
==================================

このドキュメントでは、FogFlow にコントリビューションするためのガイドラインについて説明します。コードにコントリビューションすることを計画している場合は、このドキュメントを読んで、その内容に精通する必要があります。

一般的な原則
---------------------

* FogFlow は Go プログラミング言語を使用します (ただし、テストツールや他のスクリプトなどの他のツールは python、java、bash で記述できます)。
* 効率的なコード (つまり、パフォーマンスを向上させるコード) は、非効率的なコードより優先されます。複雑なコードよりも単純なコード (つまり、よりクリーンで短いコード) が推奨されます。単純さのわずかなペナルティで効率の大幅な節約が可能です。効率のわずかなペナルティを伴う単純さの大幅な節約も可能です。
* FogFlow に提供されるコードは、コードで作業するすべての開発者に共通のプログラミングスタイルを設定するために、コード スタイル ガイドライン (`code style guidelines`_) に従う必要があります。

.. _`code style guidelines`: https://github.com/smartfog/fogflow/blob/fogflow_document_reconstruct/doc/en/source/guideline.rst#code-style-guidelines

コントリビューション ワークフロー 自体 (プルリクエストなど) は、別のドキュメント FIWARE 開発ガイドライン (`FIWARE Development Guidelines`_) で説明されていることに注意してください。

.. _`FIWARE Development Guidelines`: https://forge.fiware.org/plugins/mediawiki/wiki/fiware/index.php/Developer_Guidelines



ブランチ 管理ガイドライン
-------------------------------

.. figure:: ../../en/source/figures/gitGuideline.jpg

コミュニティには、寿命が無限の2つの主要なブランチがあります:

1. **Master branch**: これは非常に安定したブランチであり、常に本番環境に対応しており、本番環境でのソースコードの最新リリースバージョンが含まれています。
2. **Development branch**: Master ブランチ (Master branch) から派生した Development ブランチ (Development branch) は、次のリリースで計画されているさまざまな機能を統合するためのブランチとして機能します。このブランチは、Master ブランチほど安定している場合とそうでない場合があります。これは、開発者がコラボレーションして機能ブランチ (Feature Branch) をマージする場所です。すべての変更は、何らかの方法でマスターにマージして戻し、リリース番号でタグ付けする必要があります。


これらの2つの主要なブランチとは別に、ワークフローには他のブランチがあります:

- **Feature Branch**: 機能開発、つまり拡張またはドキュメント化のために Development ブランチからフォークしたブランチです。機能の開発または拡張の実装後に、Development ブランチにマージされます。

- **Bug Branch**: Development ブランチから分岐したブランチです。バグ修正後、Development ブランチにマージされます。

- **Hotfix branch**: ホットフィックス ブランチは、Master ブランチから作成されます。現在のプロダクト リリースであり、深刻なバグが原因で問題が発生していますが、開発の変更はまだ不安定です。その後、ホットフィックス ブランチから分岐して、問題の修正を開始する場合があります。重大なバグのみの場合で、これは最もまれな機会です。

**注意**: ホットフィックス ブランチを作成およびマージする権限を持っているのは NEC Laboratories Europe (NLE)、および、NEC Technologies India (NECTI) のメンバーのみです。

.. list-table::  **ブランチの命名規則** 
   :widths: 20 40 40
   :header-rows: 1

   * - ブランチ
     - ブランチ命名ガイドライン
     - 備考
     
   * - Feature branches
     - *development* から分岐する必要があります。*development* にマージして戻す必要があります。ブランチの命名規則: *feature-feature_id*
     - *feature_id* は、**https://github.com/smartfog/fogflow/issues** の Github Issue ID です。

   * - Bug Branches
     - *development* から分岐する必要があります。*development* にマージして戻す必要があります。ブランチの命名規則: *bug-bug_id*
     - *bug_id* は、**https://github.com/ScorpioBroker/ScorpioBroker/issues** の Github Issue ID です。

   * - Hotfix Branches
     - *master branch* から分岐する必要があります。*master branch* にマージして戻す必要があります。ブランチの命名規則: *hotfix-bug number*
     - *Bug number* は、**https://github.com/ScorpioBroker/ScorpioBroker/issues** の Github Issue ID です。

ブランチへのアクセス許可
*******************************

- **Master** - Master ブランチでマージしてプルリクエストを受け入れることができるのは、NLE メンバーと NECTI の特権メンバーのみであるという非常に厳しい傾向があります。マスターへのプルリクエストは、NECTI または NLE メンバーのみが作成できます。

- **Development** - コミュニティ メンバーは誰でもプルリクエストをDevelopment ブランチに上げることができますが、NLE または NECTI メンバーが確認する必要があります。Development ブランチのコミットは、travis.yml で定義されているすべてのテストケースが正常に実行された場合にのみMaster ブランチに移動されます。


コード スタイルのガイドライン
-----------------------------

現在、すべての FogFlow の既存のコードベースがこれらのルールに準拠しているわけではないことに注意してください。ガイドラインが確立される前に書かれたコードのいくつかの部分があります。ただし、すべての新しいコードの貢献はこれらのルールに従わなければならず、最終的には、古いコードはガイドラインに準拠するように変更されます。

**‘従わなければならない (MUST follow)’ ルール**

**M1 (ヘッダー ファイルを含む):**

*ルール*: すべてのヘッダーまたはソースファイルには、必要なすべてのヘッダー ファイルが含まれている必要があり、他のヘッダー ファイルは含まれていません。他のヘッダー ファイルのインクルードに依存してはなりません。また、すべてのヘッダー ファイルとソースファイルには、それ自体を必要としないヘッダー ファイルを含めてはなりません (MUST NOT)。

*根拠*: 各ファイルは、他のファイルに含まれるものと含まれないものに依存してはなりません。また、ヘッダー ファイルに必要以上のファイルが含まれている場合、その 'クライアント' には、それらの 'extra' (追加) ファイルも含める以外に選択肢はありません。これは競合につながることがあるため、回避する必要があります。さらに、コンパイル時間が長くなります。

*確認方法*: 手動

**M2 (著作権ヘッダー)**

*ルール*: すべてのファイルは、ソースコードであるかどうかに関係なく、著作権ヘッダーを持っている必要があります。

Golang ファイルの場合:

.. code-block:: console  
        /*
        *
        * Copyright 20xx The FogFlow Authors.
        *
        * This file is part of FogFlow.
        *
        * Licensed under the Apache License, Version 2.0 (the "License");
        * you may not use this file except in compliance with the License.
        *
        * You may obtain a copy of the License at
        *
        *      http://www.apache.org/licenses/LICENSE-2.0
        *
        * Unless required by applicable law or agreed to in writing, software,
        * distributed under the License is distributed on an "AS IS" BASIS,
        * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
        * See the License for the specific language governing permissions and
        * limitations under the License.
        */

Python、bash スクリプトなどの場合:

.. code-block:: console

        # Copyright 20XX FogFlow Authors.

        # This file is part of FogFlow.
        #
        # Licensed under the Apache License, Version 2.0 (the "License");
        # you may not use this file except in compliance with the License.
        # You may obtain a copy of the License at
        #
        #    http://www.apache.org/licenses/LICENSE-2.0
        #
        # Unless required by applicable law or agreed to in writing, software
        # distributed under the License is distributed on an "AS IS" BASIS,
        # WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
        # implied.
        # See the License for the specific language governing permissions and
        # limitations under the License.
        # For Python, bash script  etc.:

        # Author: <the author>

*根拠*: すべてのファイルに同種の著作権ヘッダーを設定します。

*確認方法*: 手動

**M3 (関数ヘッダー)**

*ルール*: すべての関数にはヘッダーが必要です (MUST)。ヘッダーには、関数の機能の簡単な説明、パラメーターの説明リスト、および戻り値が含まれている必要があります (SHOULD) 。

例:

.. code-block:: console  

        /* ****************************************************************************
        *
        * parseUrl - parse a URL and return its pieces
        *
        *  [ Short description if necessary ]
        *
        * PARAMETERS
        *   - url         The URL to be examined
        *   - host        To output the HOST of the URL
        *   - port        To output the PORT of the URL
        *   - path        To output the PATH of the URL
        *   - protocol    To output the PROTOCOL of the URL
        *
        * RETURN VALUE
        *   parseUrl returns TRUE on successful operation, FALSE otherwise
        *
        * NOTE
        *   About the components in a URL: according to
        *   https://tools.ietf.org/html/rfc3986#section-3,
        *   the scheme component is mandatory, i.e. the 'http://' or 'https://' must
        *   be present, otherwise the URL is invalid.
        */


*根拠*: このように準備すると、コードが読みやすくなります。

*確認方法*: 手動

**M4 (インデント)**

*ルール*: スペースのみを使用し（つまり、タブを使用せず）、一度に2つのスペース (TWO spaces) をインデントします。

*根拠*: 2つの空白で十分です。行が長くなりすぎない。

*確認方法*: 手動

**M5 (変数宣言):**

*ルール*: 宣言された各変数は、別々の行に配置する必要があります。

.. code-block:: console

        var  i  int;
        var  j  int;


次の使用は避けなければなりません (MUST):

.. code-block:: console  

        var  i, j, k int;


*根拠*: 読みやすくなります。

*確認方法*: 手動

**M6 (命名規則):**

*ルール*: 次の命名規則が適用されます。

* 名前は文字で始まる必要があり、任意の数の追加の文字と数字を含めることができます。
* 関数名を数字で始めることはできません。
* 関数名にスペースを含めることはできません。
* 名前が大文字で始まる関数が他のパッケージにエクスポートされる場合。関数名が小文字で始まる場合、他のパッケージにはエクスポートされませんが、同じパッケージ内でこの関数を呼び出すことができます。
* 関数名が複数の単語で構成されている場合は、キャメル ケースを使用してそのような名前を表します (例：empName、empAddress など）。
* 関数名では大文字と小文字が区別されます (car、Car、および CAR は3つの異なる変数です)。

*根拠*: このルールにより、理解が容易になります。

*確認方法*: 手動

**M7 (インデントやその他のフォーマットのためにコミットする前に gofmt を使用してください):**

*ルール*: gofmt -r '(a) -> a' -w FileName

* gofmt を適用する前のコード

.. code-block:: console  

        package main
                  import "fmt"
        // this is demo to format code
                   // with gofmt command
        var a int=2;
                   var b int=5;
                            var c string= `hello world`;
                   func print(){
                                fmt.Println("Value for a,b and c is : ");
                                  fmt.Println(a);
                                       fmt.Println((b));
                                             fmt.Println(c);
                               }

* ルール適用後のコード

.. code-block:: console

        package main
 
        import "fmt"
 
        // this is demo to format code
        // with gofmt command
        var a int = 2
        var b int = 5
        var c string = `hello world`
 
        func print() {
                fmt.Println("Value for a,b and c is : ")
                fmt.Println(a)
                fmt.Println((b))
                fmt.Println(c)
                     }


パッケージのフォーマットには gofmt /path/to/package を使用することに注意してください。

*根拠*: これにより、コードが再フォーマットされ、ファイルが更新されます。

*確認方法*: 手動

**M8 (コマンドと演算子の分離):**

*ルール*: 演算子 (+, *, =, == etc) の後には1つのスペースが続きます。カンマの後には1つのスペースが続きます。

.. code-block:: console

        FogFunction(va`r1, var2, var3) {
	        if (var1 == var2) {
  		         var2 = var3;
	         }
        }


ルール未適用

.. code-block:: console

        FogFunction(var1,var2,var3) {
	        if (var1==var2) {
  		        var1=var3;
	         }
        }


*根拠*: 目に優しい。

*確認方法*: 手動

**‘従わなければならない (MUST follow)’ ルール**

**S1 (エラー管理):**

*ルール*: 2番目の引数で返されたエラーは管理する必要があります。

* 悪い実装

.. code-block:: console

        FogContextElement, _ := preprocess(UpdateContextElement)

* 良い実装

.. code-block:: console

        preprocessed, err := preprocess(bytes)
        if err != nil {
          return Message{}, err
         }


**S2 (メッセージの印刷エラー):**

*ルール*: Golang の標準に従って、エラー文字列を大文字にしたり、句読点で終わらせたりしないでください。

* 悪い実装

.. code-block:: console

        if len(in) == 0 {
         return "", fmt.Errorf("Input is empty")
         }


* 良い実装

.. code-block:: console

        if len(in) == 0 {
	        return nil, errors.New("input is empty")
         }

**S3 (ネストを避ける):**

*ルール*: コードの記述中にネストを避ける。

* 悪い実装

.. code-block:: console

        func FogLine(msg *Message, in string, ch chan string) {
            if !startWith(in, stringComment) {
                token, value := parseLine(in)
                  if token != "" {
                      f, contains := factory[string(token)]
                        if !contains {
                            ch <- "ok"
                        } else {
                           data := f(token, value)
                           enrichMessage(msg, data)
                           ch <- "ok"
                        }
                        } else {
                            ch <- "ok"
                            return
                         }
                         } else {
                            ch <- "ok"
                            return
                   }
          }

* 良い実装

.. code-block:: console

        func FogLine(in []byte, ch chan interface{}) {
            // Filter empty lines and comment lines
            if len(in) == 0 || startWith(in, bytesComment) {
               ch <- nil
               return
            }
 
            token, value := parseLine(in)
            if token == nil {
                  ch <- nil
                  log.Warnf("Token name is empty on line %v", string(in))
                  return
            }

            sToken := string(token)
            if f, contains := factory[sToken]; contains {
                  ch <- f(sToken, value)
                  return
            }

            log.Warnf("Token %v is not managed by the parser", string(in))
            ch <- nil
        }


**S4 (前提条件)**

*ルール*: 処理を開始する前に、関数がパラメーターを評価し、必要に応じてエラーを返すことを強くお勧めします。

* 悪い実装

.. code-block:: console

        a, err := f1()
        if err == nil {
        b, err := f2()
            if err == nil {
                return b, nil
            } else {
                return nil, err
          }
        } else {
            return nil, err
         }

* 良い実装

.. code-block:: console

        a, err := f1()
        if err != nil {
            return nil, err
        }
        b, err := f2()
        if err != nil {
            return nil, err
        }
        return b, nil


**S5 (If 条件)**

*ルール*: Golang には、if 条件でいくつかの改善されたバージョンがあります。


* Golang での悪い実装

.. code-block:: console

        f, contains := array[index]
        if contains {
            // Do something
        }


* 良い実装

.. code-block:: console

        if f, contains := array[index]; contains {
            // Do something
        }

**S5 (Switch)**

*ルール*: スイッチ条件では常にデフォルトを使用します。


* 悪い実装

.. code-block:: console

        switch simpleToken.token {
        case tokenTitle:
            msg.Title = value
        case tokenAdep:
            msg.Adep = value
        case tokenAltnz:
            msg.Alternate = value 
         // Other cases
        }


* 良い実装 

.. code-block:: console

        switch simpleToken.token {
        case tokenTitle:
            msg.Title = value
        case tokenAdep:
            msg.Adep = value
        case tokenAltnz:
            msg.Alternate = value
        // Other cases    
        default:
            log.Errorf("unexpected token type %v", simpleToken.token)
            return Message{}, fmt.Errorf("unexpected token type %v", simpleToken.token)
        }

**S5 (コンスタント管理)**

*ルール*: コンスタント値は ADEXP と ICAO メッセージによって管理されるべきです。

* 悪い実装

.. code-block:: console

        const (
            AdexpType = 0 // TODO constant
            IcaoType  = 1
        )

* 良い実装 

.. code-block:: console

        const (
            AdexpType = iota
            IcaoType 
        )

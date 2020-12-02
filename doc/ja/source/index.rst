FogFlow
=======

.. figure:: https://nexus.lab.fiware.org/repository/raw/public/badges/chapters/processing.svg
  :target: https://www.fiware.org/developers/catalogue/ 
.. figure:: https://nexus.lab.fiware.org/repository/raw/public/badges/stackoverflow/fiware.svg
  :target: https://stackoverflow.com/questions/tagged/fiware/


FogFlowは、次のようなさまざまなコンテキスト (**context**) に基づいて、クラウドとエッジ上の動的データ処理フローを自動的に調整する IoT edge  computing framework です:

- *system context*: すべてのレイヤーから利用可能なシステム リソース;
- *data context*: 利用可能なすべてのデータエンティティの登録済みメタデータ; 
- *usage context*: QoS、遅延、および帯域幅コストの観点からユーザーが定義した予想される使用意図 (usage intention)
    
高度なインテント ベースのプログラミング モデル (intent-based programming) とコンテキスト駆動型のサービス オーケストレーションのおかげで、FogFlow は、**最小限の開発労力とほぼゼロの運用オーバーヘッドで最適化された QoS** を提供できます。現在、FogFlow は、小売、スマートシティ、スマートインダストリーの分野でさまざまなビジネス ユースケースに適用されています。

.. toctree::
    :maxdepth: 1
    :caption: イントロダクション
    :numbered:

    introduction.rst

.. toctree::
    :maxdepth: 1
    :caption: ビギナーガイド
    :numbered:
   
    onepage.rst
  
  
.. toctree::
    :maxdepth: 1
    :caption: 開発者ガイド
    :numbered:

    core_concept.rst
    intent_based_program.rst
    intent_model.rst
    guideline.rst


.. toctree::
    :maxdepth: 1
    :caption: オペレーターガイド
    :numbered:
    
    system_overview.rst
    setup.rst
    integration.rst
    fogflow_fiware_integration.rst
    scorpioIntegration.rst
    quantumleapIntegration.rst          
    wirecloudIntegration.rst
    system_monitoring.rst
    https.rst
   
.. toctree::
    :maxdepth: 1
    :caption: アドバンスユーザーガイド
    :numbered:    

    system_design.rst
    programming.rst	
    context.rst
    api.rst
    build.rst
    test.rst
    roadmap.rst

.. toctree::
    :maxdepth: 1
    :caption: その他
    
    publication.rst
    troubleshooting.rst
    contact.rst   




   
  

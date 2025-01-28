.. _trigger-topology:

サービス トポロジーをトリガーする方法
========================================

サービス トポロジーが FogFlow に送信されると、NGSI アプリケーションを使用して、FogFlow で "requirement (要件) と呼ばれる特定の NGSI エンティティのアップデートを送信することにより、このサービス トポロジーの完全または部分的なデータ処理をトリガーできます。

これは、NGSI10 アップデートを介してカスタマイズされた Requirement オブジェクトを送信するための JavaScript ベースのコード例です。この Requirement エンティティは、そのエンティティ ID によって識別されます。後で、同じ要件に変更を加えたい場合は、エンティティ ID に基づいてこの要件エンティティに NGSI10 アップデートを送信できます。

.. code-block:: javascript

    var rid = 'Requirement.' + uuid();    
   
    var requirementCtxObj = {};    
    requirementCtxObj.entityId = {
        id : rid, 
        type: 'Requirement',
        isPattern: false
    };
    
    var restriction = { scopes:[{scopeType: geoscope.type, scopeValue: geoscope.value}]};
                
    requirementCtxObj.attributes = {};   
    requirementCtxObj.attributes.output = {type: 'string', value: 'Stat'};
    requirementCtxObj.attributes.scheduler = {type: 'string', value: 'closest_first'};    
    requirementCtxObj.attributes.restriction = {type: 'object', value: restriction};    
                        
    requirementCtxObj.metadata = {};               
    requirementCtxObj.metadata.topology = {type: 'string', value: curTopology.entityId.id};
    
    console.log(requirementCtxObj);
            
	// assume the config.brokerURL is the IP of cloud IoT Broker
    var client = new NGSI10Client(config.brokerURL);				
    client.updateContext(requirementCtxObj).then( function(data) {
        console.log(data);
    }).catch( function(error) {
        console.log('failed to send a requirement');
    });    

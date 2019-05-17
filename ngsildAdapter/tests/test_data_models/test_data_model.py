
ngsi_data= \
 {
       'originator':u'',
           'subscriptionId': '195bc4c6-882e-40ce-a98f-e9b72f87bdfd',
           'contextResponses':
              [
                    {
                      'contextElement': {'attributes':
                            [
                                   {
                                          'contextValue': 'ford5',
                                           'type': 'string',
                                           'name': 'brand40'
                                   },
                                   {
                                          'contextValue': 'ford6',
                                          'type': 'string',
                                          'name': 'brand50'
                                   }
                        ],
                        'entityId':
                             {
                                   'type': 'Car',
                                    'id': 'Car31',
                                    'isPattern':True
                                 },
                      'domainMetadata':
                       [
                             {
                                  'type': 'point',
                                  'name': 'location',
                                  'value':
                                     {
                                               'latitude': 49.406393,
                                                'longitude': 8.684208
                                        }
                                }
                          ]
                  },
                  'statusCode':
                   {
                           'code': 200,
                           'reasonPhrase': 'OK'
                   }
         }
   ]
}
convert_data_output= \
    {
   '@context':
        [
                   'https://forge.etsi.org/gitlab/NGSI-LD/NGSI-LD/raw/master/coreContext/ngsi-ld-core-context.jsonld',
                        {
                                     'Car':     'http://example.org/Car',
                                     'brand40': 'http://example.org/brand40',
                                     'brand50': 'http://example.org/brand50'
                            }
                ],
                    'brand50':
                        {
                            'type': 'Property',
                            'value': 'ford6'
                        },
                     'brand40':
                         {
                             'type': 'Property',
                             'value': 'ford5'
                         },
                'type': 'Vehicle',
                'id': 'urn:ngsi-ld:Car31'
 }


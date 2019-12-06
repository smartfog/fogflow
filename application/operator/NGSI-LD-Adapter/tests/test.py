import sys, os
sys.path.append('opt/ngsildAdapter/module')

from common_utilities.rest_client import Rest_client 
from common_utilities import rest_client
from data_model.ld_generate  import ngsi_data_creation
from data_model.orian_ld_genrate import orian_convert_data
import json
import unittest
import requests
import mock 
from mock import patch
from consts import constant
from common_utilities.config import config_data
from data_model.metadata_jsonToString import metadata_converter
import data

#Input for ld_generate with DomainMetadata

class TestStringMethods(unittest.TestCase):

    def test_orian_converter(self):
        obj=orian_convert_data(data.orian_notify_data)
        returndata=obj.get_data()
        self.assertEqual(returndata,data.orian_notify_output_data)

    # Test ld_generate for point data

    def test_get_ngsi_ldPoint(self):
        obj=ngsi_data_creation(data.ngsi_data)
        result_data=obj.get_ngsi_ld()
        self.assertEqual(result_data,data.convert_data_output)
   
    # Test ld_generate for polygon data

    def testPolygonData(self):
        testObj=ngsi_data_creation(data.polygonDataInput)
        resultData=testObj.get_ngsi_ld()
        self.assertEqual(resultData,data.polygonDataOutput)

    # Test ld_generate for circle data

    def testCircleData(self):
        testObj=ngsi_data_creation(data.circleDataInput)
        resultData=testObj.get_ngsi_ld()   
        self.assertEqual(resultData,data.circleDataOutput)

   # Test ld_generate without vertices

    def testpolygonDataWithoutvertices(self):
        testObj=ngsi_data_creation(data.polygonTestDataWV_input)
        resultData=testObj.get_ngsi_ld()
        self.assertEqual(resultData,data.polygonTestDataWV_output)

    #Test ld_generate for structured attributes value

    def teststructuredvalue(self):
        testobj=ngsi_data_creation(data.TestDataForObject_input) 
        resultData=testobj.get_ngsi_ld()
        self.assertEqual(resultData,data.TestDataForObject_output)

    # Test case for entity id

    def test_get_entityId(self):
        obj=ngsi_data_creation(data.ngsi_data)
        entity_id=obj.get_entityId()
        self.assertEqual(entity_id,data.id_value)

    #mocking append request for NGB

    @patch('common_utilities.rest_client.requests.post')
    def test_mock_post(self,mock_get):
        mock_get.return_value.status_code = 201
        configobj=config_data()
        entity_url=configobj.get_entity_url()
        url1 =constant.http+entity_url+constant.entity_uri
        payload=data.convert_data_output
        payload=json.dumps(payload)
        obj=Rest_client(url1,payload)
        response=obj.post_request()
        self.assertEqual(response.status_code, 201)

    # mocking update request for NGB

    @patch('common_utilities.rest_client.requests.patch')
    def test_mock_patch(self,mock_get):
        mock_get.return_value.status_code = 204
        obj=ngsi_data_creation(data.ngsi_data)
        entity_id=obj.get_entityId()
        configobj=config_data()
        entity_url=configobj.get_entity_url()
        url=constant.http+entity_url+constant.entity_uri+entity_id+'/attrs'
        payload=data.patch_data_output
        payload=json.dumps(payload)
        obj=Rest_client(url,payload)
        response=obj.patch_request()
        self.assertEqual(response.status_code, 204)

    #testing the data without domainMetaData

    def test_data_without_domain_metadata(self):
        obj=ngsi_data_creation(data.wdmdata)
        result_data=obj.get_ngsi_ld()
        self.assertEqual(result_data,data.Eodata)
    
    #Testing the data with point domainMetaData

    def test_for_pointMetadata(self):
        test_obj=metadata_converter("",data.point_metaData)
        ngb_pointMetaData=test_obj.point_type()
        self.assertEqual(ngb_pointMetaData,data.ngb_pointMetaDataOutput)
    
    def test_get_coordinate(self):
        test_obj=metadata_converter("",data.point_metaData) 
        coordinate=test_obj.get_coordinate(data.get_coordinate_input) 
        self.assertEqual(coordinate,data.get_coordinate_output)

    def test_get_coordinate_withoutType(self):
        test_obj=metadata_converter("",data.get_coordinate_input_forWithoutType)
        coordinate=test_obj.get_coordinate(data.get_coordinate_input_forWithoutType)
        self.assertEqual(coordinate,data.get_coordinate_output_forWithoutType)

    def test_get_coordinate_withoutValue(self):
        test_obj=metadata_converter("",data.get_coordinate_input_forWithoutValue)
        coordinate=test_obj.get_coordinate(data.get_coordinate_input_forWithoutValue)
        self.assertEqual(coordinate,data.get_coordinate_output_forWithoutValue)

    def test_get_coordinate_withoutTypeValue(self):
        test_obj=metadata_converter("",data.get_coordinate_input_forWithoutTypeValue)
        coordinate=test_obj.get_coordinate(data.get_coordinate_input_forWithoutTypeValue)
        self.assertEqual(coordinate,data.get_coordinate_output_forWithoutTypeValue)    

    def test_domainMetaDataWithpolygon(self):
        test_obj=metadata_converter("",data.domainMetaData_with_polygon_input)
        ngb_polygonMetaData=test_obj.polygon_type()
        self.assertEqual(ngb_polygonMetaData,data.domainMetaData_with_withpolygon_output)

    def test_domainMetaDataWithpolygon_withoutVertices(self):
        test_obj=metadata_converter("",data.domainMetaData_with_polygon_withoutVertices_input)
        ngb_polygonMetaData=test_obj.polygon_type()
        self.assertEqual(ngb_polygonMetaData,data.domainMetaData_with_polygon_withoutVertices_output) 

    # test circle data

    def testCircleMetaData(self):
        test_obj=metadata_converter("",data.circle_metaData_input)
	ngb_circleMetaData=test_obj.circle_type() 
        self.assertEqual(ngb_circleMetaData,data.circle_metaData_output) 


if __name__ == '__main__':
    unittest.main()



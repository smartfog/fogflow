import sys,os
sys.path.append('/opt/ngsildAdapter/module')
from common_utilities.LogerHandler import Handler 
import json
from consts import constant

class metadata_converter:

    def __init__(self,context,metadata):
        self.meta_data=metadata
        self.context=context
        logger_obj=Handler()
        self.logger=logger_obj.get_logger()

# validating the MetaData

    def get_coordinate(self,meta_data):
        self.logger.info("start checking metadata")
        if self.meta_data.has_key('value')==True and self.meta_data.has_key('type')==True:
	    self.logger.debug("Reading value of metaData")
            coordinate=self.meta_data['value']
            return coordinate
        else:
           self.logger.debug("metadata is not valid")
           return constant.type_is_not_valid

# converting point metadata to ngb point metadata

    def point_type(self):
        self.logger.info("start converting point metadata to NGB string")
        self.logger.debug("Metadata type is point")
        coordinate=self.get_coordinate(self.meta_data)
        if coordinate==constant.type_is_not_valid:
            return constant.type_is_not_valid
        coordinate_data=[]
        coordinate_data.append(coordinate['latitude'])
        coordinate_data.append(coordinate['longitude'])
        d={}
        meta_data_type=self.meta_data['type']
        d["type"]="point"
        d["coordinates"]=coordinate_data
        self.logger.debug("Converting python metadata to json metadata inside file point_type")
        d2=json.dumps(d) 
        location={} 
        location['type']="GeoProperty" 
        location['value']=d2 
        #self.context['location']=location 
        self.logger.info("Creation of entity has been done")
        return location

# for polygon type

    def polygon_type(self):
         coordinate=self.get_coordinate(self.meta_data)
         # If the metaData is not valid
         if coordinate==constant.type_is_not_valid: # and coordinate.has_key('vertices')==false:
             return constant.type_is_not_valid

         # If the update domain metaData is circle
         if coordinate.has_key('vertices')==False and coordinate.has_key('centerLatitude')==True and coordinate.has_key('centerLongitude')==True and coordinate.has_key('radius')==True:
	    location =self.circle_type()
            return location
         coordinate_data=coordinate['vertices']
         coordinate_points=[]
         
         if coordinate_data is not None:
             attribute_length=len(coordinate_data) 
         else:
             attribute_length=-1 # for null vertices
         for vertex in range(attribute_length):
             point=[]
             context_attribute=coordinate_data[vertex]
             point.append(context_attribute['latitude'])
             point.append(context_attribute['longitude'])
             coordinate_points.append(point)
         d={}
         meta_data_type=self.meta_data['type']
         d["type"]="polygon"
         d["coordinates"]=coordinate_points
         self.logger.debug("Converting python metadata to json metadata inside file polygon_type")
         d2=json.dumps(d)
         location={}
         location['type']="GeoProperty"
         location['value']=d2
         self.logger.info("Creation of entity has been done")
         return location

# for circle type
    def circle_type(self):
        coordinate=self.get_coordinate(self.meta_data)
        if coordinate==constant.type_is_not_valid:
            return constant.type_is_not_valid
        coordinate_data=[]
        coordinate_data.append(coordinate['centerLatitude'])
        coordinate_data.append(coordinate['centerLongitude'])
        coordinate_data.append(coordinate['radius'])
        d={}
        meta_data_type=self.meta_data['type']
        d["type"]="polygon"
        d["coordinates"]=coordinate_data
        self.logger.debug("Converting python object to json object inside file circle_type")
        d2=json.dumps(d)
        location={}
        location['type']="GeoProperty"
        location['value']=d2
        self.logger.info("Creation of entity has been done")
        return location

    def get_converted_metadata(self):
        self.logger.info("get_converted_metadata is started")
        if self.meta_data.has_key('type')==True:

            # condition for point
            if self.meta_data['type']=="point":
                location=self.point_type()

           # condition for polygon
            if self.meta_data['type']=='polygon':
                location=self.polygon_type()

          # condition for circle
            if self.meta_data['type']=="circle":
                location=self.circle_type()

            self.logger.info("get_converted_metadata function has been end")
            if location==constant.type_is_not_valid:
                return self.context
            else:
                self.context['location']=location
                return self.context 
            
        else:
            return self.context

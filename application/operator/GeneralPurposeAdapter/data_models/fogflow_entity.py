import colorlog

from data_models.entity_data_model import Entity

logger = colorlog.getLogger('FogFlow Entity Creator')


class FogFlowEntity(Entity):
    """
    Create a FogFlow entity from a FIWARE entity
    """
    # Some entities don't have a location attribute, but instead have a latitude & longitude attributes,
    # that attributes are reported separately. For this case we have a auxiliary dictionary to store the last
    # reported latitude or longitude, when it upload both attributes we add it as domainMetadata entity field,
    # upload it and delete from the dict.
    # TODO Should have a timeout to delete devices with no latitude/longitude to avoid using too much memory space?
    device_location = {}

    def __init__(self, _id: str, _type: str, ispattern: str, attributes: list):
        super().__init__()
        self.entityId = {"id": _id, "type": _type, "isPattern": ispattern is True}
        self.attributes = []
        self.domainMetadata = []

        for attr in attributes:
            # if the entity have LOCATION/CITY attribute
            if attr.get('name') == 'location' or attr.get('name') == 'city':
                self.domainMetadata.append(
                    {"name": attr.get('name'), "type": attr.get('type'), "value": attr.get('value')})
            else:
                if FogFlowEntity.device_location.get(_id) is None:  # if is a new device register it here
                    FogFlowEntity.device_location[_id] = {}
                    logger.info("mapped new entity: {}".format(FogFlowEntity.device_location))

                # in case of entity dont have a location attribute but have longitude or latitude
                if attr.get('name') == 'longitude':
                    if attr.get('value') is not " ":
                        FogFlowEntity.device_location[_id].update({'longitude': float(attr.get('value'))})

                elif attr.get('name') == 'latitude':
                    if attr.get('value') is not " ":
                        FogFlowEntity.device_location[_id].update({'latitude': float(attr.get('value'))})

                # in other case is a normal attribute.
                else:
                    self.attributes.append(
                        {
                            "name": attr.get('name'),
                            "type": attr.get('type'),
                            "contextValue": attr.get('value')
                        }
                    )

        # If the device is stored in the device_location variable, have both longitude and latitude then add it to the
        # domainMetadata field.
        device_with_location = FogFlowEntity.device_location.get(_id)
        if device_with_location:
            if device_with_location.get('latitude') and device_with_location.get('longitude'):
                self.domainMetadata.append(
                    {
                        "name": "location",
                        "type": "point",
                        "value":
                            {
                                "latitude": device_with_location['latitude'],
                                "longitude": device_with_location['longitude']
                            }
                    }
                )
                logger.info("metadata updated, deleting entity from map {}".format(_id))
                del FogFlowEntity.device_location[_id]

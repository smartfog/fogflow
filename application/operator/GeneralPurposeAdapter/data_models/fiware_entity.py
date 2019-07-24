from data_models.entity_data_model import Entity


class FiwareEntity(Entity):
    """
    Create a FIWARE entity from a FogFlow entity
    """

    def __init__(self, entity_id: dict, attributes: list, domain_metadata=None):
        super().__init__()
        self.id = entity_id['id']
        self.type = entity_id['type']
        self.isPattern = entity_id['isPattern']
        self.attributes = []
        for attr in attributes:
            self.attributes.append(
                {
                    "name": attr.get('name'),
                    "type": attr.get('type'),
                    "value": attr.get('contextValue')
                }
            )
        if domain_metadata:
            for domain_attr in domain_metadata:
                self.attributes.append(
                    {"name": domain_attr.get('name'), "type": domain_attr.get('type'), "value": domain_attr.get('value')})


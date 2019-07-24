import secrets

import colorlog

logger = colorlog.getLogger('DeviceRegistration')


class DeviceRegistry:

    def __init__(self, device_id, iotbroker_uri, fiware_service, fiware_service_path):
        registration_id = secrets.token_hex(12)
        self.registry = {device_id: {"registrationId": registration_id, "iotbroker_uri": iotbroker_uri,
                                     "fiware-service": fiware_service,
                                     "fiware_service_path": fiware_service_path}}

    def __str__(self):
        return self.registry.__str__()

    """def __getitem__(self, item):
        return self.registry
    """

    def __get__(self, instance, owner):
        return self.registry


class DeviceRegistration:
    __devices = []

    def __init__(self):
        pass

    def add_device(self, device_id: str, iotbroker_uri: str, fiware_service: str, fiware_service_path: str):
        """for dev in DeviceRegistration.__devices:
            if device_id in dev.registry:"""
        if self.is_registered(device_id):
            return self.__update_device(device_id, iotbroker_uri, fiware_service, fiware_service_path)
        return self.__insert_device(device_id, iotbroker_uri, fiware_service, fiware_service_path)

    @staticmethod
    def __update_device(device_id: str, iotbroker_uri: str, fiware_service: str, fiware_service_path: str):
        for device in DeviceRegistration.__devices:
            if device_id in device.registry:
                DeviceRegistration.__devices.remove(device)
                DeviceRegistration.__insert_device(device_id, iotbroker_uri, fiware_service, fiware_service_path)
                return True
        return False

    @staticmethod
    def __insert_device(device_id: str, iotbroker_uri, fiware_service: str, fiware_service_path: str):
        new_element = DeviceRegistry(device_id, iotbroker_uri, fiware_service, fiware_service_path)
        DeviceRegistration.__devices.append(new_element)
        return True

    @staticmethod
    def devices_count():
        return len(DeviceRegistration.__devices)

    @staticmethod
    def empty():
        DeviceRegistration.__devices.clear()

    @staticmethod
    def get_list():
        return DeviceRegistration.__devices.copy()

    @staticmethod
    def is_registered(device_id):
        for dev in DeviceRegistration.__devices:
            if device_id in dev.registry:
                return True
        return False

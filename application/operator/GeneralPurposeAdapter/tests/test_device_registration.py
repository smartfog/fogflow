import unittest

import colorlog

from OLD.device_registration import DeviceRegistration

logger = colorlog.getLogger("test device registration")

brokerURL = "192.168.1.48:8080"


class TestDeviceRegistration(unittest.TestCase):
    device_registration = DeviceRegistration()

    def setUp(self):
        TestDeviceRegistration.device_registration.empty()

    def test_insert(self):
        self.assertEqual(
            TestDeviceRegistration.device_registration.add_device("Device.temp001", brokerURL, "devices",
                                                                  "/devices"), True)

    def test_insert2(self):
        TestDeviceRegistration.device_registration.add_device("Device.temp001", brokerURL, "devices",
                                                              "/devices")
        TestDeviceRegistration.device_registration.add_device("Device.temp002", brokerURL, "devices",
                                                              "/devices")
        self.assertEqual(TestDeviceRegistration.device_registration.devices_count(),
                         2, "devices_count should be 2")

    def test_update(self):
        TestDeviceRegistration.device_registration.add_device("Device.temp001", brokerURL, "devices",
                                                              "/devices")
        TestDeviceRegistration.device_registration.add_device("Device.temp001", brokerURL, "devices",
                                                              "/devices")
        self.assertEqual(TestDeviceRegistration.device_registration.devices_count(),
                         1, "devices_count should be 1")

    def test_update_2(self):
        TestDeviceRegistration.device_registration.add_device("Device.temp001", brokerURL, "devices",
                                                              "/devices")
        TestDeviceRegistration.device_registration.add_device("Device.temp001", brokerURL, "devicessss",
                                                              "/devices")
        self.assertEqual(TestDeviceRegistration.device_registration.devices_count(),
                         1, "devices_count should be 1")

        dic = {"Device.temp001": {"iotbroker_uri": brokerURL, "fiware-service": "devicessss",
                                  "fiware_service_path": "/devices"}}
        dic2 = TestDeviceRegistration.device_registration.get_list()[0].registry
        del dic2['Device.temp001']['registrationId']
        self.assertDictEqual(dic, dic2)


if __name__ == '__main__':
    unittest.main()

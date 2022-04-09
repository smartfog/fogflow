import socket
from common_utilities.config import config_data
import sys
sys.path.append('opt/ngsildAdapter/module')


class check_ip_port:
    def __init__(self):
        configobj = config_data()
        scarpio_broker_uri = configobj.get_entity_url()
        scarpio_broker_uri = scarpio_broker_uri.split(':')
        self.scarpio_broker_ip = scarpio_broker_uri[0]
        self.scarpio_broker_port = int(scarpio_broker_uri[1])
        fog_uri = configobj.get_fogflow_subscription_endpoint()
        fog_uri = fog_uri.split(':')
        self.fog_ip = fog_uri[0]
        self.fog_port = int(fog_uri[1])
        self.s = socket.socket()

    # check fogFlow_end_point

    def check_fog_end_point(self):
        try:
            self.s.connect((self.fog_ip, self.fog_port))
            return True
        except socket.error, e:
            print("Connection to %s on port %s failed: %s" %
                  (self.fog_ip, self.fog_port, e))
            return False

    # check_scarpio_broker_end_point

    def check_scarpio_broker_endpoint(self):
        try:
            self.s.connect((self.scarpio_broker_ip, self.scarpio_broker_port))
            return True
        except socket.error, e:
            print("Connection to %s on port %s failed: %s" %
                  (self.scarpio_broker_ip, self.scarpio_broker_port, e))
            return False

   # checking function

    def check_all_configration(self):
        fog_status = self.check_fog_end_point()
        scarpio_status = self.check_scarpio_broker_endpoint()
        if fog_status == False:
            return False
        elif scarpio_status == False:
            return False
        else:
            return True

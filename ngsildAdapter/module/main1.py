import ConfigParser
config = ConfigParser.ConfigParser()
config.readfp(open(r'config/config.ini'))
path1 = config.get('My Section', 'ngsi-ld-broker')
print(type(path1))

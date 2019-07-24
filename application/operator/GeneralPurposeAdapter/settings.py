import json
import os
current_path = os.path.dirname(os.path.abspath(__file__))
with open(current_path+'/config/configuration.json') as json_data_file:
    cfg = json.load(json_data_file)
# print(cfg)

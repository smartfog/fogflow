import logging
import datetime
import sys,os
sys.path.append('/opt/ngsildAdapter/module')
from consts import constant

class Handler:
    def __init__(self):
        pass

# creating logger

    def get_logger(self):
        log_path=self.get_file_name()
        logger_format='[ %(asctime)s ]  %(levelname)s   %(filename)s        %(message)s'
        logging.basicConfig(filename=log_path,
                    format=logger_format,
                    filemode='a+')
        logger=logging.getLogger()
        logger.setLevel(constant.logging_level)
        return logger

# creating file according to the date

    def get_file_name(self):
        file_name = datetime.date.today()
        file_location=constant.log_path+'log_file_'+str(file_name)
        return file_location
logger_obj=Handler()
logger=logger_obj.get_logger()

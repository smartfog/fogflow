import logging

import colorlog
from colorlog import ColoredFormatter

handler = colorlog.StreamHandler()

formatter = ColoredFormatter(
   "%(log_color)s%(levelname)-7s %(asctime)s.%(msecs)03d%(black)s - [%(name)s] → %(message)s",
   datefmt="%d %b %Y-%H:%M:%S",
   reset=True,
   log_colors={
       'DEBUG': 'cyan',
       'INFO': 'green',
       'WARNING': 'yellow',
       'ERROR': 'red',
       'CRITICAL': 'red,bg_white',
   },
   secondary_log_colors={
       'message': {
           'DEBUG': 'cyan',
           'INFO': 'green',
           'WARNING': 'yellow',
           'ERROR': 'red',
           'CRITICAL': 'red,bg_white'
       }
   },
   style='%'
)
# logging.FileHandler('debug.log')

handler.setFormatter(formatter)
logging.basicConfig(
   level=logging.INFO,
   handlers=[handler]
)

logger = colorlog.getLogger('Main logger')

import nxt
import sys
import time

#import nxt.locator
#from nxt.motor import *

#
#b = nxt.locator.find_one_brick()
b = nxt.find_one_brick()
mx = nxt.Motor(b, nxt.PORT_A)
mx.run(100)
time.sleep(2)
mx.brake()

print("finished")


import nxt
import sys
import time

b = nxt.find_one_brick()
mx = nxt.Motor(b, nxt.PORT_A)
mx.run(int(sys.argv[1]))
time.sleep(int(sys.argv[2]))
mx.brake()

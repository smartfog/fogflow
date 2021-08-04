import matplotlib.pyplot as plt
import numpy as np
import collections

w = 8
h = 8
d = 70

plt.figure(figsize=(w, h), dpi=d)
data = np.genfromtxt("./data.csv", delimiter=",", skip_header=1, names=["timeStamp", "Latency"])
#this part is to make the timestamps start at 0
data.sort(axis=0)
convert = lambda x: int(x)-int(data[0][0])
data = np.genfromtxt("./data.csv", delimiter=",", skip_header=1, names=["timeStamp", "Latency"], converters={0: convert})
data.sort(axis=0)
plt.plot(data['timeStamp'], data['Latency'],color='DarkSlateBlue')

#image 
plt.savefig("./data.png")

print("Overall requests " + str(len(data)))
#Average Throughput
print("Avg " + str(len(data)/(int(data[-1][0])/1000)))
print("Took " + str(int(data[-1][0])) + "ms")

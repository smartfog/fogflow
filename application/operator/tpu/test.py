#!/usr/bin/env python
import time
import os
import threading
import signal
import sys
import json
import requests 
from datetime import datetime
from PIL import Image
from io import BytesIO

# edge-tpu
from embedding import kNNEmbeddingEngine

## parameter configuration
model_path = "test_data/mobilenet_v2_1.0_224_quant_edgetpu.tflite"
width = 224
height = 224

kNN = 3
engine = kNNEmbeddingEngine(model_path, kNN)

cameraURL = "http://192.168.1.110:8040/image"

def train(label):
    global cameraURL
            
    response = requests.get(cameraURL)
    img = Image.open(BytesIO(response.content))
    
    emb = engine.DetectWithImage(img)
    engine.addEmbedding(emb, label)    


def detect():
    global cameraURL    
            
    response = requests.get(cameraURL)
    img = Image.open(BytesIO(response.content))
    
    emb = engine.DetectWithImage(img)
    result = engine.kNNEmbedding(emb)
    print("detected result %s" %(result))              

def run():
	# waiting for the inputs to decide which event to report
    print('please select a specific event to report....')  
    print('    1: normal       ')
    print('    2: defect   ')
    print('    3: detecting    ')        
    print('    0: exit         ')            
	
    while True:
        choice = input("> ")		
        
        if choice == '1':
            print("training for normal")
            for i in range(20):
                train('normal')
                print(i)
        elif choice == '2':
            print("training for normal")
            for i in range(20):
                train('defect')
                print(i)
        elif choice == '3':
            detect()
        elif choice == '0':
            break
        else:
            print('please choose 1, 2, 3')
                           
  
if __name__ == '__main__':
    run()


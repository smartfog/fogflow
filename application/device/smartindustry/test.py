print('Initializing...')

import RPi.GPIO as GPIO # Import Raspberry Pi GPIO library

## imports 
# camera
import io
from picamera import PiCamera
import time
# images
from PIL import Image
#import matplotlib.pyplot as plt
import numpy as np
from scipy import ndimage

# edge-tpu
from embedding import kNNEmbeddingEngine

import subprocess

GPIO.setwarnings(False) # Ignore warning for now
GPIO.setmode(GPIO.BOARD) # Use physical pin numbering
GPIO.setup(10, GPIO.IN, pull_up_down=GPIO.PUD_DOWN) # Set pin 10 to be an input pin and set initial value to be pulled low (off)
GPIO.setup(11, GPIO.IN, pull_up_down=GPIO.PUD_DOWN) # Set pin 11 to be an input pin and set initial value to be pulled low (off)
GPIO.setup(13, GPIO.OUT) # For LED.

## parameter configuration
model_path = "test_data/mobilenet_v2_1.0_224_quant_edgetpu.tflite"
width = 224
height = 224

kNN = 3
engine = kNNEmbeddingEngine(model_path, kNN)

def capture_image():
    stream = io.BytesIO()
    with PiCamera() as camera:
        camera.resolution = (640, 480)
        camera.capture(stream, format='jpeg')
    stream.seek(0)

    # Return image converted to PIL.
    return Image.open(stream)

def show_image(image):
    #plt.imshow(ndimage.rotate(np.asarray(image), 180))
    #plt.axis('off')
    #plt.show()
    pass
    
## learning
labels = ["good_block", "bad_block"]
count = 1

print('Ready')
while True: # Run forever
    if GPIO.input(10) == GPIO.HIGH or GPIO.input(11) == GPIO.HIGH:
        time.sleep(0.1)
        s10 = GPIO.input(10)
        s11 = GPIO.input(11)
        if s10 == GPIO.HIGH and s11 == GPIO.HIGH:
            # Reset model.
            engine = kNNEmbeddingEngine(model_path, kNN)
            GPIO.output(13,GPIO.LOW)
            for _ in range(3):
                GPIO.output(13,GPIO.HIGH)
                time.sleep(0.1)
                GPIO.output(13,GPIO.LOW)
                time.sleep(0.1)
            print("Both buttons pushed")
            pass
        else:
            label = ''
            if s10 == GPIO.HIGH:
                print("Button 1 was pushed!")
                label = labels[0]
                GPIO.output(13,GPIO.HIGH)
            elif s11 == GPIO.HIGH:
                print("Button 2 was pushed!")
                label = labels[1]
                
            print("Object: " + label)
            image = capture_image()
            show_image(image)

            # learning engine
            emb = engine.DetectWithImage(image)
            engine.addEmbedding(emb, label)
            count += 1
    else:
        image = capture_image()
        # classify the image
        emb_test = engine.DetectWithImage(image)
        res = engine.kNNEmbedding(emb_test)
        if res == 'bad_block':
            GPIO.output(13, GPIO.HIGH)
            subprocess.run(["python2", "/notebooks/motor.py", "100", "1"])
        else:
            GPIO.output(13, GPIO.LOW)
            subprocess.run(["python2", "/notebooks/motor.py", "-100", "1"])

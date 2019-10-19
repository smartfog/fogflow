sleep 2

echo "start pushbutton"
python button.py &

echo "start camera"
python camera.py &

echo "start motor"
sudo python motor.py &

echo "start tpu-based controller"
docker run -d --privileged -p 25:22 -p 8000:8000 -p 8888:8888  -v /dev/bus/usb:/dev/bus/usb --env-file env.list fogflow/tpu

echo "already set up everything"


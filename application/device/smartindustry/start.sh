sleep 2

echo "start pushbutton"
python button.py &

echo "start camera"
python camera.py &

echo "start motor"
sudo python motor.py &


echo "already set up everything"


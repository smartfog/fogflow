sleep 10

echo "start pushbutton"
stdbuf -oL  python button.py &  > button.log

echo "start camera"
stdbuf -oL  python camera.py &  > camera.log

echo "start motor"
stdbuf -oL  sudo python motor.py & > motor.log

echo "already set up everything"


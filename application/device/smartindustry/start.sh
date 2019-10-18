sleep 10

echo "start pushbutton"
python button.py &

echo "start camera"
python camera.py &

echo "start motor A"
python motorA.py &

echo "start motor B"
python motorB.py &

echo "start tpu"
python tpu.py &

echo "set it up already"


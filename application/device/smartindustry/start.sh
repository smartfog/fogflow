echo "start pushbutton"
sudo  python /home/pi/go/src/github.com/smartfog/fogflow/application/device/smartindustry/button.py & > /home/pi/button.log 2>&1

echo "start camera"
sudo  python /home/pi/go/src/github.com/smartfog/fogflow/application/device/smartindustry/camera.py & > /home/pi/camera.log 2>&1

echo "start motor"
sudo  python /home/pi/go/src/github.com/smartfog/fogflow/application/device/smartindustry/motor.py & > /home/pi/motor.log 2>&1

echo "already set up everything"


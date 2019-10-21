cd /home/pi/go/src/github.com/smartfog/fogflow/application/device/smartindustry 

echo "start pushbutton"
#sudo  python /home/pi/go/src/github.com/smartfog/fogflow/application/device/smartindustry/button.py & > /home/pi/button.log 2>&1
sudo  python button.py > /tmp/button.log &

echo "start camera"
#sudo  python /home/pi/go/src/github.com/smartfog/fogflow/application/device/smartindustry/camera.py &
sudo  python camera.py > /tmp/camera.log &

echo "start motor"
#sudo  python /home/pi/go/src/github.com/smartfog/fogflow/application/device/smartindustry/motor.py & > /home/pi/motor.log 2>&1
sudo  python motor.py > /tmp/motor.log &

echo "already set up everything"

cd powerpanel
rm *.log

node powerpanel.js profile1.json > 1.log & 
echo "start powerpanel #1"
sleep 2

node powerpanel.js profile2.json > 2.log & 
echo "start powerpanel #2"
sleep 2

node powerpanel.js profile3.json > 3.log & 
echo "start powerpanel #3"
sleep 2

cd ../camera1
rm tmp.log
python fakecamera.py > tmp.log   &
echo "start camera #1" 
sleep 2

cd ../camera2
rm tmp.log
python fakecamera.py > tmp.log   & 
echo "start camera #2"
sleep 2

cd ../
pwd


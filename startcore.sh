cd discovery
rm nohup.out
nohup ./discovery & 
echo 'started discovery'

cd ../broker
rm nohup.out
nohup ./broker & 
echo 'started broker'

cd ../master
rm nohup.out
nohup ./master &
echo 'started master'

cd ../worker
rm nohup.out
nohup ./worker &
echo 'started workers'

cd ../


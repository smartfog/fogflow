cd discovery
echo "======start to build fogflow/discovery======"
./build $1

cd ../broker
echo "======start to build fogflow/broker======"
./build $1

cd ../master
echo "======start to build fogflow/master======"
./build $1

cd ../worker
echo "======start to build fogflow/worker======"
./build $1

cd ../designer
echo "======start to build fogflow/designer======"
./build $1

cd ../
echo "!!!finished building all FogFlow core components and generating their docker images!!!" 




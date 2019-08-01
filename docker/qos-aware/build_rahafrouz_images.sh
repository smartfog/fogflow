#cd discovery
#./build
#
#cd ../broker
#./build
#
#cd ../master
#./build
#
#cd ../worker
#./build
#
#cd ../designer
#./build
#
#cd ../prometheus
#./build
#
#cd ../
#echo "finished building all FogFlow core components and generating their docker images"

build_module() { #first argument is module name
    base_path=$(dirname $(dirname $PWD))

    cd $base_path/$1
    ./build

    #build the test images
    docker build -t "rahafrouz/fogflow-$1" .

    #push to dockerhub
    docker login
    docker push "rahafrouz/fogflow-$1"
}
for module in "designer" "broker" "discovery" "worker" "master" "prometheus";
do
    build_module $module &
done

# build the images of all FogFlow core components
#./build.sh multistage

# tag all docker images with the version number
if [ $# -gt 0 ]; then
    VERSION=$1
    echo "releasing v${VERSION} to docker hub"
    
    # rename images from latest to the specific version
    docker image tag fogflow/discovery:latest fogflow/discovery:${VERSION}
    docker image tag fogflow/broker:latest fogflow/broker:${VERSION}
    docker image tag fogflow/master:latest fogflow/master:${VERSION}
    docker image tag fogflow/worker:latest fogflow/worker:${VERSION}
    docker image tag fogflow/designer:latest fogflow/designer:${VERSION}

    # publish images to the docker hub
    docker push fogflow/discovery:${VERSION}
    docker push fogflow/broker:${VERSION}
    docker push fogflow/master:${VERSION}
    docker push fogflow/worker:${VERSION}
    docker push fogflow/designer:${VERSION}

    # publish the arm version for both worker and broker
    docker image tag fogflow/worker:arm fogflow/worker_arm:${VERSION}
    docker image tag fogflow/broker:arm fogflow/broker_arm:${VERSION}

    docker push fogflow/worker_arm:${VERSION}
    docker push fogflow/broker_arm:${VERSION}
fi

FROM fogflow/mongo-iotagent-custom:latest

USER root

WORKDIR /opt/iotajson/

COPY ./function.js .
COPY ./main.js .
COPY ./iota-config.sh .
COPY ./Dockerfile .
COPY ./ngsi ./ngsi

RUN chmod +x ./iota-config.sh && \
    chmod 777 ./config.js && \
    chmod +x /usr/local/bin/docker-entrypoint.sh

CMD ["node","main.js"]

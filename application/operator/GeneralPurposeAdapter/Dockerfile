FROM dievexx/adapter:dev

WORKDIR "/app/source"

RUN apk update -f && \
    apk add --update nodejs nodejs-npm && \
    npm install shelljs && \
    npm install express && \
    npm install logops && \
    npm install axios

WORKDIR "/app/source"

COPY ./Dockerfile .
COPY ./main.js .
COPY ./ngsi ./ngsi/
COPY ./function.js .
COPY ./gpadapter-config.sh .
COPY ./build .

RUN chmod +x gpadapter-config.sh && \
    chmod 777 gpadapter-config.sh

ENTRYPOINT ["node","main.js"]

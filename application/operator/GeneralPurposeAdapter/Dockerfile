FROM python:3.7-alpine

WORKDIR "/app/source"

RUN apk add --no-cache gcc libc-dev unixodbc-dev && \
    wget https://gitlab.hopu.eu/public-projects/Py-IoTAgent-Adapter/-/archive/master/Py-IoTAgent-Adapter-master.zip && \
    unzip Py-IoTAgent-Adapter-master.zip && mv ./Py-IoTAgent-Adapter-master/* . && \
    rm -rf Py-IoTAgent-Adapter-master.zip Py-IoTAgent-Adapter-master


RUN pip install --no-cache-dir --trusted-host=pypi.python.org -r /app/source/requirements.txt && \
    apk update -f && \
    apk add --update nodejs nodejs-npm && \
    apk add zip && \
    npm install shelljs && \
    npm install express && \
    npm install logops && \
    npm install axios

EXPOSE 1026

ENV PYTHONPATH /app/source

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

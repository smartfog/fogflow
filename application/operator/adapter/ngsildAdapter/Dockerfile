FROM python:2.7-jessie

RUN mkdir /opt/ngsildAdapter && \
    apt-get update -y && \
    apt install python-pip -y && \
    pip install flask && \
    pip install requests && \
    apt-get install python-dev -y && \
    pip install ConfigParser && \ 
    apt-get install software-properties-common -y && \
    curl -sL https://deb.nodesource.com/setup_11.x | bash - && \
    apt-get install nodejs -y && \
    node -v && \
    npm install axios && \
    npm install express && \
    npm install logops && \
    npm install shelljs

WORKDIR /opt/ngsildAdapter

USER root

COPY ./ ./
RUN chmod +x transformer-config.sh && \
    chmod 777 transformer-config.sh

WORKDIR /opt/ngsildAdapter/module

ENTRYPOINT ["node","../main.js"]
EXPOSE 8880

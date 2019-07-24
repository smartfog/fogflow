FROM python:3.7-alpine

RUN apk add --no-cache gcc libc-dev unixodbc-dev

COPY requirements.txt /app/source/requirements.txt
RUN pip install --no-cache-dir --trusted-host=pypi.python.org -r /app/source/requirements.txt && \
    apk update -f && \
    apk add --update nodejs nodejs-npm && \
    npm install shelljs && \
    npm install express && \
    npm install logops && \
    npm install axios


COPY . /app/source
EXPOSE 1026
ENV PYTHONPATH /app/source
WORKDIR "/app/source"
RUN chmod +x gpadapter-config.sh && \
    chmod 777 gpadapter-config.sh


ENTRYPOINT ["node","main.js"]

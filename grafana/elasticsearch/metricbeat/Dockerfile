FROM docker.elastic.co/beats/metricbeat:7.6.0

# Metricbeat required host mount volumes
#VOLUME /hostfs/proc;
#VOLUME /hostfs/sys/fs/cgroup;
#VOLUME /hostfs;
#VOLUME /var/run/docker.sock;

# Elasticsearch host to push into
ENV ES_HOST    <ELASTICSEARCH_HOST_IP>

# Elasticsearch port to push into
ENV ES_PORT    <ELASTICSEARCH_HOST_PORT>

# Elasticsearch protocall to use
ENV ES_PROT    "http"

# Process only top X of CPU / ram processes only. 
ENV PROC_ENABLE "true"

# Top N processes by CPU or RAM to log
ENV PROC_CPU_TOP 25
ENV PROC_RAM_TOP 25

ENV PROC_METRIC_PERIOD 30m

# Socket connection monitoring
#ENV SOCKET_ENABLE "true"

# Enable / disable system monitoring
ENV SYSTEM_ENABLE "true"

# System metrics period
ENV SYSTEM_METRIC_PERIOD 10s

# Enable / disable docker conatainer monitoring
# Default : false
ENV DOCKER_ENABLE "true"

# Docker metrics, are polled on a 10 minute basis by default
ENV DOCKER_METRIC_PERIOD 10m

USER root

COPY metricbeat.docker.yml /usr/share/metricbeat/metricbeat.yml

# Perform a configration file test
#RUN metricbeat -c /usr/share/metricbeat/metricbeat.yml test config

# Assumes hostfs metricbeat binding 
# (as that is what this container is for)
ENTRYPOINT ["metricbeat"]
CMD [ "-e", "-system.hostfs=/hostfs"]

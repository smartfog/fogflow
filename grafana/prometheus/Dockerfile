

FROM prom/prometheus
USER root

COPY prometheus.yml /etc/prometheus/prometheus.yml
COPY tgroups/ /etc/prometheus/tgroups/
COPY PrometheusConfigUpdaterAPI /bin/PrometheusConfigUpdaterAPI


#A new entrypoint is create to run promtheus and our program
COPY CustomPromEntryPoint /bin/CustomPromEntryPoint

RUN    chown -R nobody:nogroup /etc/prometheus 


#From Prometheus Dockerfile
USER       nobody
#Added 4545 for talking to master
EXPOSE     9090 4545
VOLUME     [ "/prometheus" ]
WORKDIR    /prometheus
ENTRYPOINT ["/bin/CustomPromEntryPoint"]
#ENTRYPOINT [ "/bin/prometheus" ]
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--web.console.libraries=/usr/share/prometheus/console_libraries", \
             "--web.console.templates=/usr/share/prometheus/consoles" ]



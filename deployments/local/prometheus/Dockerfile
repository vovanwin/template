FROM prom/prometheus
ADD prometheus.yml /etc/prometheus/
EXPOSE 9090
ENTRYPOINT [ "/bin/prometheus" ]
CMD [ "--config.file=/etc/prometheus/prometheus.yml", \
      "--web.enable-remote-write-receiver" ]

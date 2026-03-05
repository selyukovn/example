FROM prom/prometheus:v3.5.0

COPY build/prometheus/prometheus.yml    /etc/prometheus/prometheus.yml
COPY build/prometheus/alerts.yml        /etc/prometheus/alerts.yml
COPY build/prometheus/records.yml       /etc/prometheus/records.yml

VOLUME /prometheus

CMD ["--config.file=/etc/prometheus/prometheus.yml"]

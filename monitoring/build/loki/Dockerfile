FROM grafana/loki:3.6.3

COPY build/loki/loki.yml /etc/loki/loki.yml

VOLUME /loki

# -config.expand-env -- можно использовать переменные вида ${VAR:-default_value}
CMD ["-config.file=/etc/loki/loki.yml", "-config.expand-env"]

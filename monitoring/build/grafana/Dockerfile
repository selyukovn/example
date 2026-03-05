FROM grafana/grafana:12.1

COPY build/grafana/config.ini       /etc/grafana/config.ini
COPY build/grafana/datasources.yml  /etc/grafana/provisioning/datasources/all.yaml
COPY build/grafana/dashboards.yml   /etc/grafana/provisioning/dashboards/all.yaml
COPY build/grafana/dashboards/      /etc/grafana/dashboards/

VOLUME /var/lib/grafana

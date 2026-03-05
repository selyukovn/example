FROM prom/alertmanager:v0.30.1

COPY build/alertmanager/alertmanager.yml /etc/alertmanager/alertmanager.yml

CMD ["--config.file=/etc/alertmanager/alertmanager.yml", "--log.level=DEBUG"]

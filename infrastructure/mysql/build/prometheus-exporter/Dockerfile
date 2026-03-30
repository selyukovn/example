FROM prom/mysqld-exporter:v0.17.2

ARG MYSQL_USER
ARG MYSQL_PASSWORD

WORKDIR /config

# Eсли нет секции [client], то контейнер запустится, но будет сгенерирована ошибка
# level=ERROR source=config.go:141
# msg="failed to validate config" section=client err="no user specified in section or parent"
# даже не смотря на то, что все используемые клиенты прописаны.
RUN echo "[client]"                      > /config/.my.cnf \
 && echo "user=${MYSQL_USER}"           >> /config/.my.cnf \
 && echo "password=${MYSQL_PASSWORD}"   >> /config/.my.cnf \
 && echo "[client.master]"              >> /config/.my.cnf \
 && echo "user=${MYSQL_USER}"           >> /config/.my.cnf \
 && echo "password=${MYSQL_PASSWORD}"   >> /config/.my.cnf \
 && echo "[client.replica]"             >> /config/.my.cnf \
 && echo "user=${MYSQL_USER}"           >> /config/.my.cnf \
 && echo "password=${MYSQL_PASSWORD}"   >> /config/.my.cnf

CMD ["--config.my-cnf=/config/.my.cnf"]

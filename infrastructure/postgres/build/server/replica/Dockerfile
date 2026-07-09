FROM postgres:17.6-alpine

ENV INIT_REPLICATOR_PASSWORD=""
ENV INIT_REPLICATION_HOST=""
ENV INIT_REPLICATION_SLOT=""

VOLUME /var/log/postgres

# ----------------------------------------------------------------------------------------------------------------------
# Конфиги
# ----------------------------------------------------------------------------------------------------------------------

COPY build/server/replica/pg_hba.conf   /etc/postgres/pg_hba.conf

COPY build/server/common.conf           /etc/postgres/common.conf
COPY build/server/replica/custom.conf   /etc/postgres/custom.conf
RUN cat /etc/postgres/common.conf /etc/postgres/custom.conf > /etc/postgres/postgresql.conf

# ----------------------------------------------------------------------------------------------------------------------
# Запуск
# ----------------------------------------------------------------------------------------------------------------------

COPY build/server/replica/entrypoint.sh     /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD ["postgres", "-c", "hba_file=/etc/postgres/pg_hba.conf", "-c", "config_file=/etc/postgres/postgresql.conf"]

# ----------------------------------------------------------------------------------------------------------------------
# Healthcheck
# ----------------------------------------------------------------------------------------------------------------------

COPY build/server/healthcheck.sh /healthcheck.sh
RUN chmod +x /healthcheck.sh

HEALTHCHECK \
    # Время на запуск контейнера и интервал проверок в это время.
    # Неудачи игнорируются до первой удачной проверки или до конца периода времени.
    --start-period=10s \
    --start-interval=5s \
    # С момента запуска контейнера,
    # если вернет не 0, не выполнится или будет выполняться дольше -- неудача.
    # Если много неудач -- провал.
    --interval=1m \
    --timeout=3s \
    --retries=3 \
    CMD ["/healthcheck.sh"]

# ----------------------------------------------------------------------------------------------------------------------

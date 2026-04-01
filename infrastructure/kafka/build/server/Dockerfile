# В `confluentinc/cp-kafka:8.2.0` есть бага с параметром `KAFKA_ADVERTISED_LISTENERS`.
# https://github.com/confluentinc/kafka-images/issues/467

FROM confluentinc/cp-kafka:8.1.2

# ----------------------------------------------------------------------------------------------------------------------
# Logs
# ----------------------------------------------------------------------------------------------------------------------

# https://docs.confluent.io/platform/current/installation/docker/development.html#log-to-external-volumes
COPY build/server/log4j2.yaml.template      /etc/confluent/docker/log4j2.yaml.template
VOLUME /var/log/kafka/custom

# ----------------------------------------------------------------------------------------------------------------------
# Healthcheck
# ----------------------------------------------------------------------------------------------------------------------

COPY --chown=appuser:appuser build/server/healthcheck.sh /etc/confluent/docker/healthcheck.sh
RUN chmod +x /etc/confluent/docker/healthcheck.sh

HEALTHCHECK \
    # Время на запуск контейнера и интервал проверок в это время.
    # Неудачи игнорируются до первой удачной проверки или до конца периода времени.
    --start-period=10s \
    --start-interval=5s \
    # С момента запуска контейнера,
    # если вернет не 0, не выполнится или будет выполняться дольше -- неудача.
    # Если много неудач -- провал.
    --interval=10s \
    --timeout=3s \
    --retries=3 \
    CMD ["/etc/confluent/docker/healthcheck.sh"]

# ----------------------------------------------------------------------------------------------------------------------

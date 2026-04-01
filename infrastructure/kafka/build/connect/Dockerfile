
# ######################################################################################################################
# PLUGINS
# ######################################################################################################################

# Установка коннектора через `confluent-hub` внутри самого `confluentinc/cp-kafka-connect` недоступна без VPN.
# https://docs.confluent.io/platform/current/installation/docker/development.html#create-a-docker-image-containing-c-hub-connectors
# https://www.confluent.io/hub/debezium/debezium-connector-mysql
# RUN confluent-hub install --no-prompt debezium/debezium-connector-mysql:3.2.6 -- зависает намертво.
#
# Кроме того, нет базовых wget, tar и нет инсталлятора (по крайней мере, ни apk, ни apt, хотя образ основан на ubuntu).
# Так что проще загрузить плагины отдельным шагом.
# https://debezium.io/documentation/reference/3.5/install.html
FROM alpine:3.22.1 AS plugins

WORKDIR /plugins

# debezium-connector-mysql
RUN wget https://repo1.maven.org/maven2/io/debezium/debezium-connector-mysql/3.5.0.Final/debezium-connector-mysql-3.5.0.Final-plugin.tar.gz && \
    tar -xzf debezium-connector-mysql-3.5.0.Final-plugin.tar.gz && rm debezium-connector-mysql-3.5.0.Final-plugin.tar.gz

# jmx prometheus exporter
# Это не плагин, но не хочется отдельный шаг делать для его скачивания.
# https://prometheus.github.io/jmx_exporter/1.5.0/java-agent/
RUN wget https://github.com/prometheus/jmx_exporter/releases/download/1.5.0/jmx_prometheus_javaagent-1.5.0.jar

# ######################################################################################################################
# KAFKA-CONNECT
# ######################################################################################################################

FROM confluentinc/cp-kafka-connect:8.2.0

# Часть переменных, требуемых для работы, устанавливается в `docker-compose`-файле,
# но некоторые значения можно зафиксировать здесь, поскольку они "статичны".

ENV CONNECT_GROUP_ID=1
ENV CONNECT_CONFIG_STORAGE_TOPIC="infrastructure-kafka-connect-config"
ENV CONNECT_OFFSET_STORAGE_TOPIC="infrastructure-kafka-connect-offset"
ENV CONNECT_STATUS_STORAGE_TOPIC="infrastructure-kafka-connect-status"

ENV CONNECT_KEY_CONVERTER="org.apache.kafka.connect.storage.StringConverter"
ENV CONNECT_VALUE_CONVERTER="org.apache.kafka.connect.json.JsonConverter"

# требует переменную, а 8083 -- значение по умолчанию.
ENV CONNECT_REST_PORT=8083

# ----------------------------------------------------------------------------------------------------------------------
# Plugins
# ----------------------------------------------------------------------------------------------------------------------

COPY --from=plugins /plugins/debezium-connector-mysql /usr/share/java/debezium-connector-mysql

# ----------------------------------------------------------------------------------------------------------------------
# Monitoring
# ----------------------------------------------------------------------------------------------------------------------

WORKDIR /monitoring

COPY --from=plugins /plugins/jmx_prometheus_javaagent-1.5.0.jar ./

# https://docs.confluent.io/platform/current/installation/docker/operations/monitoring.html
ENV KAFKA_JMX_OPTS="-javaagent:/monitoring/jmx_prometheus_javaagent-1.5.0.jar=7073:/monitoring/monitoring.yml"

COPY build/connect/monitoring.yml ./

EXPOSE 7073

# ----------------------------------------------------------------------------------------------------------------------
# Launcher
# ----------------------------------------------------------------------------------------------------------------------

WORKDIR /launch

COPY --chown=appuser:appuser build/connect/launcher.sh  ./
RUN chmod +x ./launcher.sh

COPY build/connect/*.py                                 ./

CMD ["/launch/launcher.sh"]

# ######################################################################################################################

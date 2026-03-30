FROM mysql:8.4

# master | replica
ARG SERVER_TYPE
# int > 0
ARG SERVER_ID
ARG EXTRA_REPLICA_SOURCE_HOST="localhost"
ARG EXTRA_REPLICA_SOURCE_PORT=3306

# Использование в docker-compose-файле переменной окружения исходного образа
# с именем, не сообщающем об однократном использовании этой переменной лишь при инициализации, вводит в заблуждение.
# Получение значения из аргумента с более говорящим именем яснее отражает суть.
ARG INITIAL_ROOT_PASSWORD
ENV MYSQL_ROOT_PASSWORD=${INITIAL_ROOT_PASSWORD}

EXPOSE 3306

# /var/lib/mysql -- уже объявлен как VOLUME в базовом образе, поэтому изменить ничего нельзя.

# ----------------------------------------------------------------------------------------------------------------------
# Конфиги
# ----------------------------------------------------------------------------------------------------------------------

# В образе, в конце /etc/my.cnf, есть строка "!includedir /etc/mysql/conf.d/"
COPY build/server/common.cnf                     /etc/mysql/conf.d/common.cnf
COPY build/server/${SERVER_TYPE}/custom.cnf      /etc/mysql/conf.d/custom.cnf

RUN echo "[mysqld]"                          >> /etc/mysql/conf.d/server_id.cnf \
 && echo "server-id=${SERVER_ID}"            >> /etc/mysql/conf.d/server_id.cnf

# ----------------------------------------------------------------------------------------------------------------------
# Скрипты инициализации
# ----------------------------------------------------------------------------------------------------------------------

# При первом запуске выполнятся в алфавитном порядке файлы из /docker-entrypoint-initdb.d
COPY build/server/${SERVER_TYPE}/init.sql    /docker-entrypoint-initdb.d/init.sql

RUN sed -i -e "s|__REPLICA_SOURCE_HOST__|${EXTRA_REPLICA_SOURCE_HOST}|g" /docker-entrypoint-initdb.d/init.sql \
 && sed -i -e "s|__REPLICA_SOURCE_PORT__|${EXTRA_REPLICA_SOURCE_PORT}|g" /docker-entrypoint-initdb.d/init.sql

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

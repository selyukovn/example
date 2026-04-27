# https://hub.docker.com/_/redis#redisversion-alpine
# "... The main caveat to note is that it does use musl libc instead of glibc and friends,
# so software will often run into issues depending on the depth of their libc requirements/assumptions. ..."
FROM redis:8.4.2-alpine

# *.conf файлы
ARG SERVER_TYPE

WORKDIR /launch

COPY build/server/launcher.sh       /launch
COPY build/server/init.sh           /launch
COPY build/server/healthcheck.sh    /launch

RUN chmod +x /launch/launcher.sh && \
    chmod +x /launch/init.sh && \
    chmod +x /launch/healthcheck.sh

COPY build/server/${SERVER_TYPE}.conf   /usr/local/etc/redis/redis.conf

ENTRYPOINT [ "/launch/launcher.sh", "/usr/local/etc/redis/redis.conf" ]

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
    CMD ["/launch/healthcheck.sh"]
